// +build !ignore_autogenerated

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
// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AclEntry) DeepCopyInto(out *AclEntry) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AclEntry.
func (in *AclEntry) DeepCopy() *AclEntry {
	if in == nil {
		return nil
	}
	out := new(AclEntry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupConfiguration) DeepCopyInto(out *BackupConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupConfiguration.
func (in *BackupConfiguration) DeepCopy() *BackupConfiguration {
	if in == nil {
		return nil
	}
	out := new(BackupConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseFlags) DeepCopyInto(out *DatabaseFlags) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseFlags.
func (in *DatabaseFlags) DeepCopy() *DatabaseFlags {
	if in == nil {
		return nil
	}
	out := new(DatabaseFlags)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseInstanceFailoverReplica) DeepCopyInto(out *DatabaseInstanceFailoverReplica) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseInstanceFailoverReplica.
func (in *DatabaseInstanceFailoverReplica) DeepCopy() *DatabaseInstanceFailoverReplica {
	if in == nil {
		return nil
	}
	out := new(DatabaseInstanceFailoverReplica)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseInstanceSpec) DeepCopyInto(out *DatabaseInstanceSpec) {
	*out = *in
	if in.ReplicaConfiguration != nil {
		in, out := &in.ReplicaConfiguration, &out.ReplicaConfiguration
		*out = new(ReplicaConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.Settings != nil {
		in, out := &in.Settings, &out.Settings
		*out = new(Settings)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseInstanceSpec.
func (in *DatabaseInstanceSpec) DeepCopy() *DatabaseInstanceSpec {
	if in == nil {
		return nil
	}
	out := new(DatabaseInstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseInstanceStatus) DeepCopyInto(out *DatabaseInstanceStatus) {
	*out = *in
	if in.FailoverReplica != nil {
		in, out := &in.FailoverReplica, &out.FailoverReplica
		*out = new(DatabaseInstanceFailoverReplica)
		**out = **in
	}
	if in.IpAddresses != nil {
		in, out := &in.IpAddresses, &out.IpAddresses
		*out = make([]*IpMapping, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(IpMapping)
				**out = **in
			}
		}
	}
	if in.ReplicaNames != nil {
		in, out := &in.ReplicaNames, &out.ReplicaNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SuspensionReason != nil {
		in, out := &in.SuspensionReason, &out.SuspensionReason
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseInstanceStatus.
func (in *DatabaseInstanceStatus) DeepCopy() *DatabaseInstanceStatus {
	if in == nil {
		return nil
	}
	out := new(DatabaseInstanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Instance) DeepCopyInto(out *Instance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Instance.
func (in *Instance) DeepCopy() *Instance {
	if in == nil {
		return nil
	}
	out := new(Instance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Instance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceList) DeepCopyInto(out *InstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Instance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceList.
func (in *InstanceList) DeepCopy() *InstanceList {
	if in == nil {
		return nil
	}
	out := new(InstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceSpec) DeepCopyInto(out *InstanceSpec) {
	*out = *in
	if in.DatabaseInstanceSpec != nil {
		in, out := &in.DatabaseInstanceSpec, &out.DatabaseInstanceSpec
		*out = new(DatabaseInstanceSpec)
		(*in).DeepCopyInto(*out)
	}
	out.ProviderRef = in.ProviderRef
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceSpec.
func (in *InstanceSpec) DeepCopy() *InstanceSpec {
	if in == nil {
		return nil
	}
	out := new(InstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceStatus) DeepCopyInto(out *InstanceStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	if in.DatabaseInstanceStatus != nil {
		in, out := &in.DatabaseInstanceStatus, &out.DatabaseInstanceStatus
		*out = new(DatabaseInstanceStatus)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceStatus.
func (in *InstanceStatus) DeepCopy() *InstanceStatus {
	if in == nil {
		return nil
	}
	out := new(InstanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IpConfiguration) DeepCopyInto(out *IpConfiguration) {
	*out = *in
	if in.AuthorizedNetworks != nil {
		in, out := &in.AuthorizedNetworks, &out.AuthorizedNetworks
		*out = make([]*AclEntry, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(AclEntry)
				**out = **in
			}
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IpConfiguration.
func (in *IpConfiguration) DeepCopy() *IpConfiguration {
	if in == nil {
		return nil
	}
	out := new(IpConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IpMapping) DeepCopyInto(out *IpMapping) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IpMapping.
func (in *IpMapping) DeepCopy() *IpMapping {
	if in == nil {
		return nil
	}
	out := new(IpMapping)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocationPreference) DeepCopyInto(out *LocationPreference) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocationPreference.
func (in *LocationPreference) DeepCopy() *LocationPreference {
	if in == nil {
		return nil
	}
	out := new(LocationPreference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MaintenanceWindow) DeepCopyInto(out *MaintenanceWindow) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MaintenanceWindow.
func (in *MaintenanceWindow) DeepCopy() *MaintenanceWindow {
	if in == nil {
		return nil
	}
	out := new(MaintenanceWindow)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MySqlReplicaConfiguration) DeepCopyInto(out *MySqlReplicaConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MySqlReplicaConfiguration.
func (in *MySqlReplicaConfiguration) DeepCopy() *MySqlReplicaConfiguration {
	if in == nil {
		return nil
	}
	out := new(MySqlReplicaConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicaConfiguration) DeepCopyInto(out *ReplicaConfiguration) {
	*out = *in
	if in.MysqlReplicaConfiguration != nil {
		in, out := &in.MysqlReplicaConfiguration, &out.MysqlReplicaConfiguration
		*out = new(MySqlReplicaConfiguration)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicaConfiguration.
func (in *ReplicaConfiguration) DeepCopy() *ReplicaConfiguration {
	if in == nil {
		return nil
	}
	out := new(ReplicaConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Settings) DeepCopyInto(out *Settings) {
	*out = *in
	if in.AuthorizedGaeApplications != nil {
		in, out := &in.AuthorizedGaeApplications, &out.AuthorizedGaeApplications
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.BackupConfiguration != nil {
		in, out := &in.BackupConfiguration, &out.BackupConfiguration
		*out = new(BackupConfiguration)
		**out = **in
	}
	if in.DatabaseFlags != nil {
		in, out := &in.DatabaseFlags, &out.DatabaseFlags
		*out = make([]*DatabaseFlags, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(DatabaseFlags)
				**out = **in
			}
		}
	}
	if in.IpConfiguration != nil {
		in, out := &in.IpConfiguration, &out.IpConfiguration
		*out = new(IpConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.LocationPreference != nil {
		in, out := &in.LocationPreference, &out.LocationPreference
		*out = new(LocationPreference)
		**out = **in
	}
	if in.MaintenanceWindow != nil {
		in, out := &in.MaintenanceWindow, &out.MaintenanceWindow
		*out = new(MaintenanceWindow)
		**out = **in
	}
	if in.StorageAutoResize != nil {
		in, out := &in.StorageAutoResize, &out.StorageAutoResize
		*out = new(bool)
		**out = **in
	}
	if in.UserLabels != nil {
		in, out := &in.UserLabels, &out.UserLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Settings.
func (in *Settings) DeepCopy() *Settings {
	if in == nil {
		return nil
	}
	out := new(Settings)
	in.DeepCopyInto(out)
	return out
}
