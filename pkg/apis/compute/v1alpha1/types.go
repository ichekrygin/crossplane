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

package v1alpha1

import (
	corev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/rbac/v1beta1"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// KubernetesClusterSpec specifies the configuration of a Kubernetes cluster.
type KubernetesClusterSpec struct {
	ClassRef    *corev1.ObjectReference `json:"classReference,omitempty"`
	ResourceRef *corev1.ObjectReference `json:"resourceName,omitempty"`
	Selector    metav1.LabelSelector    `json:"selector,omitempty"`

	// cluster properties
	ClusterVersion string `json:"clusterVersion,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesCluster is the Schema for the instances API
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.bindingPhase"
// +kubebuilder:printcolumn:name="CLUSTER-CLASS",type="string",JSONPath=".spec.classReference.name"
// +kubebuilder:printcolumn:name="CLUSTER-REF",type="string",JSONPath=".spec.resourceName.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type KubernetesCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubernetesClusterSpec            `json:"spec,omitempty"`
	Status corev1alpha1.ResourceClaimStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesClusterList contains a list of KubernetesClusters.
type KubernetesClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubernetesCluster `json:"items"`
}

// ObjectReference to using this object as a reference
func (kc *KubernetesCluster) ObjectReference() *corev1.ObjectReference {
	if kc.Kind == "" {
		kc.Kind = KubernetesInstanceKind
	}
	if kc.APIVersion == "" {
		kc.APIVersion = APIVersion
	}
	return &corev1.ObjectReference{
		APIVersion: kc.APIVersion,
		Kind:       kc.Kind,
		Name:       kc.Name,
		Namespace:  kc.Namespace,
		UID:        kc.UID,
	}
}

// OwnerReference to use this object as an owner
func (kc *KubernetesCluster) OwnerReference() metav1.OwnerReference {
	return *util.ObjectToOwnerReference(kc.ObjectReference())
}

// ClaimStatus returns the claim status of this Kubernetes cluster.
func (kc *KubernetesCluster) ClaimStatus() *corev1alpha1.ResourceClaimStatus {
	return &kc.Status
}

// ClassRef returns the resource class used by this Kubernetes cluster.
func (kc *KubernetesCluster) ClassRef() *corev1.ObjectReference {
	return kc.Spec.ClassRef
}

// ResourceRef returns the resource claimed by this Kubernetes cluster.
func (kc *KubernetesCluster) ResourceRef() *corev1.ObjectReference {
	return kc.Spec.ResourceRef
}

// SetResourceRef sets the resource claimed by this Kubernetes cluster.
func (kc *KubernetesCluster) SetResourceRef(ref *corev1.ObjectReference) {
	kc.Spec.ResourceRef = ref
}

// ResourceReference is generic resource represented by the resource name and the secret name that will be generated
// for the consumption inside the Workload.
// TODO: Note, currently resource reference is a general type, however, this will be change in the future and replaced with concrete resource types
type ResourceReference struct {
	// reference to a resource object in the same namespace
	corev1.ObjectReference `json:",inline"`
	// name of the generated resource secret
	SecretName string `json:"secretName"`
}

// WorkloadState represents the state of a workload.
type WorkloadState string

// Workload states.
const (
	WorkloadStateCreating WorkloadState = "CREATING"
	WorkloadStateRunning  WorkloadState = "RUNNING"
)

// WorkloadSpec specifies the configuration of a workload.
type WorkloadSpec struct {
	ClusterSelector map[string]string `json:"clusterSelector,omitempty"`

	TargetNamespace           string                            `json:"targetNamespace"`
	TargetCRDs                []apiext.CustomResourceDefinition `json:"targetCrds,omitempty"`
	TargetCRs                 []unstructured.Unstructured       `json:"targetCrs,omitempty"`
	TargetServiceAccounts     []corev1.ServiceAccount           `json:"targetServiceAccounts,omitempty"`
	TargetClusterRoles        []v1beta1.ClusterRole             `json:"targetClusterRoles,omitempty"`
	TargetClusterRoleBindings []v1beta1.ClusterRoleBinding      `json:"targetClusterRoleBindings,omitempty"`
	TargetRoles               []v1beta1.Role                    `json:"targetRoles,omitempty"`
	TargetRoleBindings        []v1beta1.RoleBinding             `json:"targetRoleBindings,omitempty"`
	TargetDeployments         []appsv1.Deployment               `json:"targetDeployments,omitempty"`
	TargetServices            []corev1.Service                  `json:"targetServices,omitempty"`

	// Resources
	Resources []ResourceReference `json:"resources,omitempty"`
}

type WorkloadStatusObjectReferences struct {
	// CRDs
	CRDs                []corev1.ObjectReference `json:"crds,omitempty"`
	CRs                 []corev1.ObjectReference `json:"crs,omitempty"`
	SAs                 []corev1.ObjectReference `json:"serviceAccounts,omitempty"`
	ClusterRoles        []corev1.ObjectReference `json:"clusterRoles,omitempty"`
	ClusterRoleBindings []corev1.ObjectReference `json:"clusterRoleBindings,omitempty"`
	Deployments         []corev1.ObjectReference `json:"deployments,omitempty"`
	Services            []corev1.ObjectReference `json:"services,omitempty"`
}

type WorkloadStatusObjectStatuts struct {
	CRDs        []apiext.CustomResourceDefinitionStatus `json:"crds,omitempty"`
	Deployments []appsv1.DeploymentStatus               `json:"deployments,omitempty"`
	Services    []corev1.ServiceStatus                  `json:"services,omitempty"`
}

// WorkloadStatus represents the status of a workload.
type WorkloadStatus struct {
	corev1alpha1.ConditionedStatus

	Cluster *corev1.ObjectReference `json:"clusterRef,omitempty"`
	State   WorkloadState           `json:"state,omitempty"`

	TargetReferences WorkloadStatusObjectReferences `json:"targetReferences,omitempty"`
	TargetStatuses   WorkloadStatusObjectStatuts    `json:"targetStatuses,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Workload is the Schema for the instances API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.state"
// +kubebuilder:printcolumn:name="CLUSTER",type="string",JSONPath=".status.clusterRef.name"
// +kubebuilder:printcolumn:name="NAMESPACE",type="string",JSONPath=".spec.targetNamespace"
// +kubebuilder:printcolumn:name="DEPLOYMENT",type="string",JSONPath=".spec.targetDeployment.metadata.name"
// +kubebuilder:printcolumn:name="SERVICE-EXTERNAL-IP",type="string",JSONPath=".status.service.loadBalancer.ingress[0].ip"
type Workload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkloadSpec   `json:"spec,omitempty"`
	Status WorkloadStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WorkloadList contains a list of Workloads.
type WorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workload `json:"items"`
}

// ObjectReference to using this object as a reference
func (kc *Workload) ObjectReference() *corev1.ObjectReference {
	if kc.Kind == "" {
		kc.Kind = KubernetesInstanceKind
	}
	if kc.APIVersion == "" {
		kc.APIVersion = APIVersion
	}
	return &corev1.ObjectReference{
		APIVersion: kc.APIVersion,
		Kind:       kc.Kind,
		Name:       kc.Name,
		Namespace:  kc.Namespace,
		UID:        kc.UID,
	}
}

// OwnerReference to use this object as an owner
func (kc *Workload) OwnerReference() metav1.OwnerReference {
	return *util.ObjectToOwnerReference(kc.ObjectReference())
}
