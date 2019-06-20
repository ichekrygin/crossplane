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
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/util"
)

// CloudSQL instance states
const (
	// StateRunnable represents a CloudSQL instance in a running, available, and ready state
	StateRunnable = "RUNNABLE"
)

type InstanceSpec struct {
	*DatabaseInstanceSpec `json:"instance"`

	// NameFormat to format instance name passing it a object UID
	// If not provided, defaults to "%s", i.e. UID value
	NameFormat string `json:"nameFormat,omitempty"`

	// Provider reference
	ProviderRef core.LocalObjectReference `json:"providerRef"`

	// ReclaimPolicy identifies how to handle the cloud resource after the deletion of this type
	ReclaimPolicy corev1alpha1.ReclaimPolicy `json:"reclaimPolicy,omitempty"`
}

// InstanceStatus defines the observed state of Instance
type InstanceStatus struct {
	corev1alpha1.ConditionedStatus

	*DatabaseInstanceStatus `json:",inline"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Instance is the Schema for the instances API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="STATE",type="string",JSONPath=".status.state"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec,omitempty"`
	Status InstanceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstanceList contains a list of Instance
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}

// ObjectReference to this CloudSQL instance instance
func (in *Instance) ObjectReference() *core.ObjectReference {
	return util.ObjectReference(in.ObjectMeta, util.IfEmptyString(in.APIVersion, APIVersion), util.IfEmptyString(in.Kind, InstanceKind))
}

// OwnerReference to use this instance as an owner
func (in *Instance) OwnerReference() metav1.OwnerReference {
	return *util.ObjectToOwnerReference(in.ObjectReference())
}

// IsAvailable for usage/binding
func (in *Instance) IsAvailable() bool {
	return in.Status.State == StateRunnable
}

// DatabaseInstanceName based on the NameFormat spec value,
// If name format is not provided, instance name defaults to UID
// If name format provided with '%s' value, instance name will result in formatted string + UID,
//   NOTE: only single %s substitution is supported
// If name format does not contain '%s' substitution, i.e. a constant string, the
// constant string value is returned back
//
// Examples:
//   For all examples assume "UID" = "test-uid"
//   1. NameFormat = "", DatabaseInstanceName = "test-uid"
//   2. NameFormat = "%s", DatabaseInstanceName = "test-uid"
//   3. NameFormat = "foo", DatabaseInstanceName = "foo"
//   4. NameFormat = "foo-%s", DatabaseInstanceName = "foo-test-uid"
//   5. NameFormat = "foo-%s-bar-%s", DatabaseInstanceName = "foo-test-uid-bar-%!s(MISSING)"
func (in *Instance) DatabaseInstanceName() string {
	return util.ConditionalStringFormat(in.Spec.NameFormat, string(in.GetUID()))
}

// DatabaseInstance creates a new Cloudsql DatabaseInstance using properties defines in object spec
// Note: current instance can be nil - if there is no current instances
// The reason why we are passing current instance is to "copy" existing properties
func (in *Instance) DatabaseInstance(current *sqladmin.DatabaseInstance) *sqladmin.DatabaseInstance {
	instance := saveDatabaseInstanceSpec(in.Spec.DatabaseInstanceSpec)
	instance.Name = in.DatabaseInstanceName()
	if current != nil {
		instance.Name = current.Name
		instance.Settings.SettingsVersion = current.Settings.SettingsVersion
	}
	return instance
}

func (in *Instance) DatabaseInstanceSpec(instance *sqladmin.DatabaseInstance) {
	in.Spec.DatabaseInstanceSpec = newDatabaseInstanceSpec(instance)
}

func (in *Instance) DatabaseInstanceStatus(instance *sqladmin.DatabaseInstance) {
	in.Status.DatabaseInstanceStatus = newDatabaseInstanceStatus(instance)
}

func (in *Instance) DatabaseInstanceSpecDiff(instance *sqladmin.DatabaseInstance) string {
	desired := newDatabaseInstanceSpec(instance)
	return cmp.Diff(desired, in.Spec.DatabaseInstanceSpec, cmpopts.EquateEmpty())
}
