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

package cloudsql

import (
	"context"
	"time"

	corev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/apis/gcp/cloudsql/v1alpha1"
	gcpv1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/gcp/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/clients/gcp/cloudsql"
	"github.com/crossplaneio/crossplane/pkg/logging"
	"github.com/crossplaneio/crossplane/pkg/util/googleapi"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
	corev1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	ctrlhandler "sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	controllerName = "cloudsql"
	finalizer      = "finalizer" + controllerName

	reconcileTimeout = 1 * time.Minute

	resourceSyncDelay   = 60 * time.Second
	resourceStatusDelay = 10 * time.Second
)

var (
	resultNoRequeue             = reconcile.Result{}
	resultRequeueNow            = reconcile.Result{Requeue: true}
	resultResourceSyncRequeue   = reconcile.Result{RequeueAfter: resourceSyncDelay}
	resultResourceStatusRequeue = reconcile.Result{RequeueAfter: resourceStatusDelay}
)

type handler interface {
	create(context.Context) (reconcile.Result, error)
	update(context.Context, *sqladmin.DatabaseInstance) (reconcile.Result, error)
	delete(context.Context) (reconcile.Result, error)
}

// instanceCreateUpdater implementation of createupdater interface
type resourceHandler struct {
	operations
}

// create new cloudsql instance
func (rh *resourceHandler) create(ctx context.Context) (reconcile.Result, error) {
	if err := rh.addFinalizer(ctx); err != nil {
		return resultRequeueNow, errors.Wrap(err, "failed to update object")
	}

	if err := rh.Create(ctx, rh.getInstance(nil)); err != nil {
		return resultRequeueNow, rh.reconcileStatus(ctx, err)
	}

	// Retrieve newly create instance
	inst, err := rh.Get(ctx, rh.getInstanceName())
	if err != nil {
		return resultRequeueNow, rh.reconcileStatus(ctx, err)
	}

	// Sync back instance properties to the object
	return resultNoRequeue, rh.reconcileStatus(ctx, rh.syncInstance(ctx, inst))
}

// delete cloudsql instance
func (rh *resourceHandler) delete(ctx context.Context) (reconcile.Result, error) {
	if rh.isReclaimDelete() {
		if err := rh.Delete(ctx, rh.getInstanceName()); err != nil && !googleapi.IsErrorNotFound(err) {
			return resultRequeueNow, rh.reconcileStatus(ctx, err)
		}
	}
	return resultNoRequeue, rh.removeFinalizer(ctx)
}

// update cloudsql instance if needed
func (rh *resourceHandler) update(ctx context.Context, instance *sqladmin.DatabaseInstance) (reconcile.Result, error) {
	// Sync status
	if err := rh.syncInstanceStatus(ctx, instance); err != nil {
		return resultRequeueNow, rh.reconcileStatus(ctx, err)
	}

	// Check if instance is ready, if not - wait for it
	if !rh.isInstanceReady() {
		// TODO(illya) consider a timeout to avoid a perpetual wait for resource that may never will be ready
		//   upon reaching the timeout convert from `resultResourceSyncRequeue` to `resultRequeueNow` w/ back-off
		return resultResourceStatusRequeue, rh.reconcileStatus(ctx, nil)
	}

	// Check if the instance needs to be updated
	if rh.needsUpdate(instance) {
		// Patch instance
		if err := rh.Patch(ctx, rh.getInstance(instance)); err != nil {
			// If patch failed due to operations conflict - delay requeue
			result := resultRequeueNow
			if googleapi.IsErrorConflict(err) {
				result = resultResourceStatusRequeue
			}
			return result, rh.reconcileStatus(ctx, err)
		}
		return resultRequeueNow, rh.reconcileStatus(ctx, nil)
	}

	// No changes detected - delay requeue
	return resultResourceSyncRequeue, rh.reconcileStatus(ctx, nil)
}

// newInstanceHandler new instance of resourceCreateUpdater
func newInstanceHandler(ops operations) *resourceHandler {
	return &resourceHandler{
		operations: ops,
	}
}

// objectReconciler defines reconciliation for cloudsql instance object
type objectReconciler interface {
	reconcile(ctx context.Context, object *v1alpha1.Instance) (reconcile.Result, error)
}

// instanceReconciler implements reconciliation of cloudsql instance object
type instanceReconciler struct {
	handler
	operations
	log logr.Logger
}

// reconcile cloudsql instance object
func (r instanceReconciler) reconcile(ctx context.Context, obj *v1alpha1.Instance) (reconcile.Result, error) {
	name := r.getInstanceName()
	log := r.log.WithName("reconcile").WithValues("instance", name)
	debug := log.V(logging.Debug).Info

	if obj.DeletionTimestamp != nil {
		debug("Delete")
		return r.delete(ctx)
	}

	debug("Retrieve")
	resource, err := r.Get(ctx, name)
	if err != nil && !googleapi.IsErrorNotFound(err) {
		return resultRequeueNow, r.reconcileStatus(ctx, err)
	}

	if resource == nil {
		debug("Create")
		return r.create(ctx)
	}

	debug("Updating")
	return r.update(ctx, resource)
}

// reconcilerMaker defines operations to create new cloudsql instance reconciler
type reconcilerMaker interface {
	newReconciler(context.Context, *v1alpha1.Instance, logr.Logger) (objectReconciler, error)
}

// instanceReconcilerMaker implements reconcilerMaker interface and defines new reconciler constructor
type instanceReconcileMaker struct {
	client.Client
}

// newReconciler returns a reconciler for a specific object instance
func (m instanceReconcileMaker) newReconciler(ctx context.Context, obj *v1alpha1.Instance, log logr.Logger) (objectReconciler, error) {
	// Fetch provider object
	p := &gcpv1alpha1.Provider{}
	n := types.NamespacedName{Namespace: obj.GetNamespace(), Name: obj.Spec.ProviderRef.Name}
	if err := m.Get(ctx, n, p); err != nil {
		return nil, err
	}

	// Fetch provider credentials secret
	s := &corev1.Secret{}
	n = types.NamespacedName{Namespace: p.GetNamespace(), Name: p.Spec.Secret.Name}
	if err := m.Get(ctx, n, s); err != nil {
		return nil, errors.Wrapf(err, "cannot get provider's secret %s", n)
	}

	// Create google credentials
	creds, err := google.CredentialsFromJSON(context.Background(), s.Data[p.Spec.Secret.Key], sqladmin.SqlserviceAdminScope)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot retrieve creds from json")
	}

	// Create sqladmin service client
	sc, err := sqladmin.New(oauth2.NewClient(ctx, creds.TokenSource))
	if err != nil {
		return nil, errors.Wrapf(err, "error creating sqladmin service")
	}

	// Create operations
	ops := &instanceOperations{
		client:          m.Client,
		object:          obj,
		InstanceService: cloudsql.NewInstanceService(creds.ProjectID, sc),
		log:             log.WithName("operations"),
	}

	// Return initialized instance reconciler
	return &instanceReconciler{
		handler:    newInstanceHandler(ops),
		operations: ops,
		log:        log,
	}, nil
}

// Reconciler reconciles a Instance object
type Reconciler struct {
	client.Client
	maker reconcilerMaker
	log   logr.Logger
}

// Add creates a new Instance Controller and adds it to the Manager with default RBAC.
// The Manager will set fields on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	r := &Reconciler{
		Client: mgr.GetClient(),
		maker: &instanceReconcileMaker{
			Client: mgr.GetClient(),
		},
		log: logging.Logger.WithName(controllerName).WithName("controller"),
	}

	// Create a newSyncDeleter controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to bucket
	return c.Watch(&source.Kind{Type: &v1alpha1.Instance{}}, &ctrlhandler.EnqueueRequestForObject{})
}

// Reconcile reads that state of the cluster for a Provider bucket and makes changes based on the state read
// and what is in the Provider.Spec
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("key", request.NamespacedName)

	log.V(logging.Debug).Info("Reconcile started", "kind", v1alpha1.InstanceKindAPIVersion)

	ctx, cancel := context.WithTimeout(context.Background(), reconcileTimeout)
	defer cancel()

	obj := &v1alpha1.Instance{}
	if err := r.Get(ctx, request.NamespacedName, obj); err != nil {
		if apierr.IsNotFound(err) {
			return resultNoRequeue, nil
		}
		return resultNoRequeue, err
	}
	log.V(logging.Debug).Info("Object", "version", obj.ResourceVersion)

	//time.Sleep(2 * time.Second)
	//
	//obj = &v1alpha1.Instance{}
	//if err := r.Get(ctx, request.NamespacedName, obj); err != nil {
	//	if apierr.IsNotFound(err) {
	//		return resultNoRequeue, nil
	//	}
	//	return resultNoRequeue, err
	//}
	//glog.WithValues("key", request.NamespacedName).V(logging.Debug).Info("Reconciling-2", "version", obj.ResourceVersion)

	rr, err := r.maker.newReconciler(ctx, obj, log)
	if err != nil {
		obj.Status.SetConditions(corev1alpha1.ReconcileError(err))
		return resultRequeueNow, r.Status().Update(ctx, obj)
	}

	rs, err := rr.reconcile(ctx, obj)
	log.V(logging.Debug).Info("Reconcile finished",
		"requeue", rs.Requeue, "delay", rs.RequeueAfter.String(),
		"error", err)
	return rs, err
}
