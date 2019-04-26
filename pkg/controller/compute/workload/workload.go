/*
Copyright 2018 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package workload

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/rbac/v1beta1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	computev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/compute/v1alpha1"
	corev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/logging"
	"github.com/crossplaneio/crossplane/pkg/util"
)

const (
	controllerName = "workload.compute.crossplane.io"
	finalizer      = "finalizer." + controllerName

	errorClusterClient = "Failed to create cluster client"
	errorCreating      = "Failed to create"
	errorSynchronizing = "Failed to sync"
	errorDeleting      = "Failed to delete"

	workloadReferenceLabelKey = "workloadRef"

	// supported kinds
	kindCrd                = "customresourcedefinitions"
	kindSA                 = "serviceaccounts"
	kindClusterRole        = "clusterroles"
	kindClusterRoleBinding = "clusterrolebindings"
	kindDeployment         = "deployments"
	kindService            = "services"
)

var (
	log           = logging.Logger.WithName("controller." + controllerName)
	ctx           = context.Background()
	resultDone    = reconcile.Result{}
	resultRequeue = reconcile.Result{Requeue: true}
)

// Add creates a new Instance Controller and adds it to the Manager with default RBAC.
// The Manager will set fields on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// Reconciler reconciles a Instance object
type Reconciler struct {
	client.Client
	scheme     *runtime.Scheme
	kubeclient kubernetes.Interface
	recorder   record.EventRecorder

	connect func(*computev1alpha1.Workload) (client.Client, error)
	create  func(*computev1alpha1.Workload, client.Client) (reconcile.Result, error)
	sync    func(*computev1alpha1.Workload, client.Client) (reconcile.Result, error)
	delete  func(*computev1alpha1.Workload, client.Client) (reconcile.Result, error)

	propagate                    func(client.Client, runtime.Object, string, string) error
	propagateCRs                 func(client.Client, []unstructured.Unstructured, string, string) ([]corev1.ObjectReference, error)
	propagateCRDs                func(client.Client, []apiextensions.CustomResourceDefinition, string, string) ([]corev1.ObjectReference, error)
	propagateServiceAccounts     func(client.Client, []corev1.ServiceAccount, string, string) ([]corev1.ObjectReference, error)
	propagateClusterRoles        func(client.Client, []v1beta1.ClusterRole, string, string) ([]corev1.ObjectReference, error)
	propagateClusterRoleBindings func(client.Client, []v1beta1.ClusterRoleBinding, string, string) ([]corev1.ObjectReference, error)
	propagateDeployments         func(client.Client, []appsv1.Deployment, string, string) ([]corev1.ObjectReference, error)
	propagateServices            func(client.Client, []corev1.Service, string, string) ([]corev1.ObjectReference, error)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	r := &Reconciler{
		Client:                       mgr.GetClient(),
		scheme:                       mgr.GetScheme(),
		kubeclient:                   kubernetes.NewForConfigOrDie(mgr.GetConfig()),
		recorder:                     mgr.GetRecorder(controllerName),
		propagate:                    propagate,
		propagateCRs:                 propagateCRs,
		propagateCRDs:                propagateCRDs,
		propagateServiceAccounts:     propagateServiceAccounts,
		propagateClusterRoles:        propagateClusterRoles,
		propagateClusterRoleBindings: propagateClusterRoleBindings,
		propagateDeployments:         propagateDeployments,
		propagateServices:            propagateServices,
	}
	r.connect = r._connect
	r.create = r._create
	r.sync = r._sync
	r.delete = r._delete

	return r
}

// CreatePredicate accepts Workload instances with set `Status.Cluster` reference value
func CreatePredicate(event event.CreateEvent) bool {
	wl, ok := event.Object.(*computev1alpha1.Workload)
	if !ok {
		return false
	}
	return wl.Status.Cluster != nil
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Instance
	err = c.Watch(&source.Kind{Type: &computev1alpha1.Workload{}}, &handler.EnqueueRequestForObject{}, &predicate.Funcs{CreateFunc: CreatePredicate})
	if err != nil {
		return err
	}

	return nil
}

// fail - helper function to set fail condition with reason and message
func (r *Reconciler) fail(instance *computev1alpha1.Workload, reason, msg string) (reconcile.Result, error) {
	instance.Status.SetCondition(corev1alpha1.NewCondition(corev1alpha1.Failed, reason, msg))
	return resultRequeue, r.Status().Update(ctx, instance)
}

// _connect establish connection to the target cluster
func (r *Reconciler) _connect(instance *computev1alpha1.Workload) (client.Client, error) {
	ref := instance.Status.Cluster
	if ref == nil {
		return nil, fmt.Errorf("workload is not scheduled")
	}

	k := &computev1alpha1.KubernetesCluster{}

	err := r.Get(ctx, client.ObjectKey{Namespace: ref.Namespace, Name: ref.Name}, k)
	if err != nil {
		return nil, err
	}

	s, err := r.kubeclient.CoreV1().Secrets(k.Namespace).Get(k.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// read the individual connection config fields
	user := s.Data[corev1alpha1.ResourceCredentialsSecretUserKey]
	pass := s.Data[corev1alpha1.ResourceCredentialsSecretPasswordKey]
	ca := s.Data[corev1alpha1.ResourceCredentialsSecretCAKey]
	cert := s.Data[corev1alpha1.ResourceCredentialsSecretClientCertKey]
	key := s.Data[corev1alpha1.ResourceCredentialsSecretClientKeyKey]
	token := s.Data[corev1alpha1.ResourceCredentialsTokenKey]
	host, ok := s.Data[corev1alpha1.ResourceCredentialsSecretEndpointKey]
	if !ok {
		return nil, fmt.Errorf("kubernetes cluster endpoint/host is not found")
	}
	u, err := url.Parse(string(host))
	if err != nil {
		return nil, fmt.Errorf("cannot parse Kubernetes endpoint as URL: %+v", err)
	}

	config := &rest.Config{
		Host:     u.String(),
		Username: string(user),
		Password: string(pass),
		TLSClientConfig: rest.TLSClientConfig{
			// This field's godoc claims clients will use 'the hostname used to
			// contact the server' when it is left unset. In practice clients
			// appear to use the URL, including scheme and port.
			ServerName: u.Hostname(),
			CAData:     ca,
			CertData:   cert,
			KeyData:    key,
		},
		BearerToken: string(token),
	}
	apiextensions.SchemeBuilder.AddToScheme(scheme.Scheme)
	return client.New(config, client.Options{Scheme: scheme.Scheme})
}

func addWorkloadReferenceLabel(m metav1.Object, uid string) {
	labels := m.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[workloadReferenceLabelKey] = uid
	m.SetLabels(labels)
}

func getWorkloadReferenceLabel(m metav1.Object) string {
	labels := m.GetLabels()
	if len(labels) > 0 {
		return labels[workloadReferenceLabelKey]
	}
	return ""
}

func propagate(c client.Client, o runtime.Object, ns, uid string) error {
	mo, ok := o.(metav1.Object)
	if !ok {
		return fmt.Errorf("not valid runtime.Object: %T", o)
	}
	mo.SetNamespace(util.IfEmptyString(mo.GetNamespace(), ns))

	addWorkloadReferenceLabel(mo, uid)

	if err := c.Create(ctx, o); err != nil {
		if kerrors.IsAlreadyExists(err) {
			if getWorkloadReferenceLabel(mo) == uid {
				//return c.Update(ctx, o)
				return nil
			}
		}
		return err
	}
	return nil
}

func propagateCRs(c client.Client, crs []unstructured.Unstructured, ns, uid string) ([]corev1.ObjectReference, error) {
	var refs []corev1.ObjectReference

	for _, i := range crs {
		if err := propagate(c, &i, ns, uid); err != nil {
			return nil, err
		}
		ref := corev1.ObjectReference{
			Kind:            i.GetKind(),
			APIVersion:      i.GetAPIVersion(),
			Namespace:       i.GetNamespace(),
			Name:            i.GetName(),
			UID:             i.GetUID(),
			ResourceVersion: i.GetResourceVersion(),
		}
		refs = append([]corev1.ObjectReference{ref}, refs...)
	}
	return refs, nil
}

func propagateCRDs(c client.Client, list []apiextensions.CustomResourceDefinition, ns, uid string) ([]corev1.ObjectReference, error) {
	if len(list) == 0 {
		return nil, nil
	}

	var refs []corev1.ObjectReference

	for _, i := range list {
		if err := propagate(c, &i, ns, uid); err != nil {
			return nil, err
		}
		refs = append([]corev1.ObjectReference{
			newObjectReference(i.TypeMeta, i.ObjectMeta, kindCrd, apiextensions.SchemeGroupVersion.String())}, refs...)
	}
	return refs, nil
}

func propagateServiceAccounts(c client.Client, list []corev1.ServiceAccount, ns, uid string) ([]corev1.ObjectReference, error) {
	if len(list) == 0 {
		return nil, nil
	}

	var refs []corev1.ObjectReference

	for _, i := range list {
		if err := propagate(c, &i, ns, uid); err != nil {
			return nil, err
		}
		refs = append([]corev1.ObjectReference{
			newObjectReference(i.TypeMeta, i.ObjectMeta, kindSA, corev1.SchemeGroupVersion.String())}, refs...)
	}
	return refs, nil
}

func propagateClusterRoles(c client.Client, list []v1beta1.ClusterRole, ns, uid string) ([]corev1.ObjectReference, error) {
	if len(list) == 0 {
		return nil, nil
	}

	var refs []corev1.ObjectReference

	for _, i := range list {
		if err := propagate(c, &i, ns, uid); err != nil {
			return nil, err
		}
		refs = append([]corev1.ObjectReference{
			newObjectReference(i.TypeMeta, i.ObjectMeta, kindClusterRole, v1beta1.SchemeGroupVersion.String())}, refs...)
	}
	return refs, nil
}

func propagateClusterRoleBindings(c client.Client, list []v1beta1.ClusterRoleBinding, ns, uid string) ([]corev1.ObjectReference, error) {
	if len(list) == 0 {
		return nil, nil
	}

	var refs []corev1.ObjectReference

	for _, i := range list {
		if err := propagate(c, &i, ns, uid); err != nil {
			return nil, err
		}
		refs = append([]corev1.ObjectReference{
			newObjectReference(i.TypeMeta, i.ObjectMeta, kindClusterRoleBinding, v1beta1.SchemeGroupVersion.String())}, refs...)
	}
	return refs, nil
}

func propagateDeployments(c client.Client, list []appsv1.Deployment, ns, uid string) ([]corev1.ObjectReference, error) {
	if len(list) == 0 {
		return nil, nil
	}

	var refs []corev1.ObjectReference

	for _, i := range list {
		if err := propagate(c, &i, ns, uid); err != nil {
			return nil, err
		}
		refs = append([]corev1.ObjectReference{
			newObjectReference(i.TypeMeta, i.ObjectMeta, kindDeployment, appsv1.SchemeGroupVersion.String())}, refs...)
	}
	return refs, nil
}

func propagateServices(c client.Client, list []corev1.Service, ns, uid string) ([]corev1.ObjectReference, error) {
	if len(list) == 0 {
		return nil, nil
	}

	var refs []corev1.ObjectReference

	for _, i := range list {
		if err := propagate(c, &i, ns, uid); err != nil {
			return nil, err
		}
		refs = append([]corev1.ObjectReference{
			newObjectReference(i.TypeMeta, i.ObjectMeta, kindService, corev1.SchemeGroupVersion.String())}, refs...)
	}
	return refs, nil
}

// _create workload
func (r *Reconciler) _create(instance *computev1alpha1.Workload, remote client.Client) (reconcile.Result, error) {
	log.V(1).Info("creating workload")

	if !util.HasFinalizer(&instance.ObjectMeta, finalizer) {
		util.AddFinalizer(&instance.ObjectMeta, finalizer)
		return resultDone, errors.Wrapf(r.Update(ctx, instance), "failed to update the object after adding finalizer")
	}

	instance.Status.SetCreating()

	// create target namespace
	targetNamespace := instance.Spec.TargetNamespace

	if err := remote.Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: targetNamespace}}); err != nil && !kerrors.IsAlreadyExists(err) {
		return r.fail(instance, errorCreating, err.Error())
	}

	uid := string(instance.UID)

	// propagate resources secrets
	for _, resource := range instance.Spec.Resources {
		// retrieve secret
		secretName := util.IfEmptyString(resource.SecretName, resource.Name)
		sec, err := r.kubeclient.CoreV1().Secrets(instance.Namespace).Get(secretName, metav1.GetOptions{})
		if err != nil {
			return r.fail(instance, errorCreating, err.Error())
		}

		// create secret
		sec.ObjectMeta = metav1.ObjectMeta{
			Name:      sec.Name,
			Namespace: targetNamespace,
		}
		addWorkloadReferenceLabel(&sec.ObjectMeta, uid)
		if err := util.Apply(ctx, remote, sec); err != nil {
			return r.fail(instance, errorCreating, err.Error())
		}
	}

	var err error
	trefs := &instance.Status.TargetReferences

	// propagate crds
	if trefs.CRDs, err = r.propagateCRDs(remote, instance.Spec.TargetCRDs, targetNamespace, uid); err != nil {
		return r.fail(instance, errorCreating, err.Error())
	}

	// propagate crs
	if trefs.CRs, err = r.propagateCRs(remote, instance.Spec.TargetCRs, targetNamespace, uid); err != nil {
		return r.fail(instance, errorCreating, err.Error())
	}

	// propagate service accounts
	if trefs.SAs, err = r.propagateServiceAccounts(remote, instance.Spec.TargetServiceAccounts, targetNamespace, uid); err != nil {
		return r.fail(instance, errorCreating, err.Error())
	}

	// propagate cluster roles
	if trefs.ClusterRoles, err = r.propagateClusterRoles(remote, instance.Spec.TargetClusterRoles, targetNamespace, uid); err != nil {
		return r.fail(instance, errorCreating, err.Error())
	}

	// propagate cluster role bindings
	if trefs.ClusterRoleBindings, err = r.propagateClusterRoleBindings(remote, instance.Spec.TargetClusterRoleBindings, targetNamespace, uid); err != nil {
		return r.fail(instance, errorCreating, err.Error())
	}

	// propagate deployment
	if trefs.Deployments, err = r.propagateDeployments(remote, instance.Spec.TargetDeployments, targetNamespace, uid); err != nil {
		return r.fail(instance, errorCreating, err.Error())
	}

	// propagate service
	if trefs.Services, err = r.propagateServices(remote, instance.Spec.TargetServices, targetNamespace, uid); err != nil {
		return r.fail(instance, errorCreating, err.Error())
	}

	instance.Status.State = computev1alpha1.WorkloadStateCreating
	err = r.Status().Update(ctx, instance)
	// update instance
	return resultDone, err
}

func syncCRDs(c client.Client, instance *computev1alpha1.Workload) ([]apiextensions.CustomResourceDefinitionStatus, error) {
	var statuses []apiextensions.CustomResourceDefinitionStatus
	for _, v := range instance.Status.TargetReferences.CRDs {
		crd := &apiextensions.CustomResourceDefinition{}
		nn := util.NamespaceNameFromObjectRef(&v)
		if err := c.Get(ctx, nn, crd); err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve crd: %s", nn)
		}
		statuses = append(statuses, crd.Status)
	}
	return statuses, nil
}

func syncDeployments(c client.Client, instance *computev1alpha1.Workload) ([]appsv1.DeploymentStatus, error) {
	var statuses []appsv1.DeploymentStatus
	for _, v := range instance.Status.TargetReferences.Deployments {
		i := &appsv1.Deployment{}
		nn := util.NamespaceNameFromObjectRef(&v)
		if err := c.Get(ctx, nn, i); err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve crd: %s", nn)
		}
		statuses = append(statuses, i.Status)
	}
	return statuses, nil
}

func syncServices(c client.Client, instance *computev1alpha1.Workload) ([]corev1.ServiceStatus, error) {
	var statuses []corev1.ServiceStatus
	for _, v := range instance.Status.TargetReferences.Deployments {
		i := &corev1.Service{}
		nn := util.NamespaceNameFromObjectRef(&v)
		if err := c.Get(ctx, nn, i); err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve crd: %s", nn)
		}
		statuses = append(statuses, i.Status)
	}
	return statuses, nil
}

// _sync Workload status
func (r *Reconciler) _sync(instance *computev1alpha1.Workload, c client.Client) (reconcile.Result, error) {
	var first, err error
	if instance.Status.TargetStatuses.CRDs, err = syncCRDs(c, instance); err != nil {
		first = errors.Wrapf(err, "failed to sync crds")
	}

	if instance.Status.TargetStatuses.Deployments, err = syncDeployments(c, instance); err != nil {
		if first != nil {
			first = errors.Wrapf(err, "failed to sync deployments")
		}
	}

	if instance.Status.TargetStatuses.Services, err = syncServices(c, instance); err != nil {
		if first != nil {
			first = errors.Wrapf(err, "failed to sync services")
		}
	}

	if first != nil {
		return resultDone, first
	}

	ready := true

	for _, ds := range instance.Status.TargetStatuses.Deployments {
		ready = ready && util.LatestDeploymentCondition(ds.Conditions).Type == appsv1.DeploymentAvailable
	}

	if ready {
		instance.Status.State = computev1alpha1.WorkloadStateRunning
		instance.Status.SetCondition(corev1alpha1.NewCondition(corev1alpha1.Ready, "", ""))
	}

	return resultDone, r.Status().Update(ctx, instance)
}

func deleteCRs(c client.Client, refs []corev1.ObjectReference) error {
	for _, ref := range refs {
		obj := &unstructured.Unstructured{}
		obj.SetName(ref.Name)
		obj.SetNamespace(ref.Namespace)
		if err := c.Delete(ctx, obj); err != nil {
			return errors.Wrapf(err, "failed to delete custom resource: %v", obj)
		}
	}
	return nil
}

func deleteCRDs(c client.Client, refs []corev1.ObjectReference) error {
	for _, ref := range refs {
		meta := newMetaFromObjectReference(ref)
		if err := c.Delete(ctx, &apiextensions.CustomResourceDefinition{ObjectMeta: meta}); err != nil && !kerrors.IsNotFound(err) {
			return errors.Wrapf(err, "failed to delete crd: %v", meta)
		}
	}
	return nil
}

func deleteServiceAccount(c client.Client, refs []corev1.ObjectReference) error {
	for _, ref := range refs {
		meta := newMetaFromObjectReference(ref)
		if err := c.Delete(ctx, &corev1.ServiceAccount{ObjectMeta: meta}); err != nil && !kerrors.IsNotFound(err) {
			return errors.Wrapf(err, "failed to delete service account: %v", meta)
		}
	}
	return nil
}

func deleteClusterRole(c client.Client, refs []corev1.ObjectReference) error {
	for _, ref := range refs {
		meta := newMetaFromObjectReference(ref)
		if err := c.Delete(ctx, &v1beta1.ClusterRole{ObjectMeta: meta}); err != nil && !kerrors.IsNotFound(err) {
			return errors.Wrapf(err, "failed to delete cluster role: %v", meta)
		}
	}
	return nil
}

func deleteClusterRoleBindings(c client.Client, refs []corev1.ObjectReference) error {
	for _, ref := range refs {
		meta := newMetaFromObjectReference(ref)
		if err := c.Delete(ctx, &v1beta1.ClusterRoleBinding{ObjectMeta: meta}); err != nil && !kerrors.IsNotFound(err) {
			return errors.Wrapf(err, "failed to delete cluster role binding: %v", meta)
		}
	}
	return nil
}

func deleteDeployments(c client.Client, refs []corev1.ObjectReference) error {
	for _, ref := range refs {
		meta := newMetaFromObjectReference(ref)
		if err := c.Delete(ctx, &appsv1.Deployment{ObjectMeta: meta}); err != nil && !kerrors.IsNotFound(err) {
			return errors.Wrapf(err, "failed to delete deployment: %v", meta)
		}
	}
	return nil
}

func deleteServices(c client.Client, refs []corev1.ObjectReference) error {
	for _, ref := range refs {
		meta := newMetaFromObjectReference(ref)
		if err := c.Delete(ctx, &corev1.Service{ObjectMeta: meta}); err != nil && !kerrors.IsNotFound(err) {
			return errors.Wrapf(err, "failed to delete service: %v", meta)
		}
	}
	return nil
}

// _delete workload
func (r *Reconciler) _delete(instance *computev1alpha1.Workload, c client.Client) (reconcile.Result, error) {
	if err := deleteServices(c, instance.Status.TargetReferences.Services); err != nil {
		return resultDone, errors.Wrapf(err, "failed to delete services")
	}
	if err := deleteDeployments(c, instance.Status.TargetReferences.Deployments); err != nil {
		return resultDone, errors.Wrapf(err, "failed to delete deployments")
	}
	if err := deleteServiceAccount(c, instance.Status.TargetReferences.SAs); err != nil {
		return resultDone, errors.Wrapf(err, "failed to delete service accounts")
	}
	if err := deleteClusterRoleBindings(c, instance.Status.TargetReferences.ClusterRoleBindings); err != nil {
		return resultDone, errors.Wrapf(err, "failed to delete cluster role bindings")
	}
	if err := deleteClusterRole(c, instance.Status.TargetReferences.ClusterRoles); err != nil {
		return resultDone, errors.Wrapf(err, "failed to delete cluster roles")
	}
	if err := deleteCRs(c, instance.Status.TargetReferences.CRs); err != nil {
		return resultDone, errors.Wrapf(err, "failed to delete custom resources")
	}
	if err := deleteCRDs(c, instance.Status.TargetReferences.CRDs); err != nil {
		return resultDone, errors.Wrapf(err, "failed to delete crds")
	}

	instance.Status.SetCondition(corev1alpha1.NewCondition(corev1alpha1.Deleting, "", ""))
	if err := r.Status().Update(ctx, instance); err != nil {
		return resultDone, errors.Wrapf(err, "failed to updated status")
	}

	util.RemoveFinalizer(&instance.ObjectMeta, finalizer)
	return resultDone, r.Update(ctx, instance)
}

// Reconcile reads that state of the cluster for a Instance object and makes changes based on the state read
// and what is in the Instance.Spec
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.V(logging.Debug).Info("reconciling", "kind", computev1alpha1.WorkloadKindAPIVersion, "request", request)
	// fetch the CRD instance
	instance := &computev1alpha1.Workload{}
	instance.GetNamespace()

	err := r.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if kerrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return resultDone, nil
		}
		return resultDone, err
	}

	if instance.Status.Cluster == nil {
		return resultDone, nil
	}

	// target cluster client
	targetClient, err := r.connect(instance)
	if err != nil {
		return r.fail(instance, errorClusterClient, err.Error())
	}

	// Check for deletion
	if instance.DeletionTimestamp != nil {
		return r.delete(instance, targetClient)
	}

	if instance.Status.State == "" {
		return r.create(instance, targetClient)
	}

	return r.sync(instance, targetClient)
}

func newObjectReference(t metav1.TypeMeta, o metav1.ObjectMeta, kind, apiversion string) corev1.ObjectReference {
	return corev1.ObjectReference{
		Kind:            util.IfEmptyString(t.Kind, kind),
		APIVersion:      util.IfEmptyString(t.APIVersion, apiversion),
		Namespace:       o.Namespace,
		Name:            o.Name,
		UID:             o.UID,
		ResourceVersion: o.ResourceVersion,
	}
}

func newMetaFromObjectReference(ref corev1.ObjectReference) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: ref.Namespace,
		Name:      ref.Name,
	}
}
