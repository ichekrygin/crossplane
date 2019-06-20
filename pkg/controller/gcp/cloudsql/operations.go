package cloudsql

import (
	"context"

	corev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/apis/gcp/cloudsql/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/clients/gcp/cloudsql"
	"github.com/crossplaneio/crossplane/pkg/logging"
	"github.com/crossplaneio/crossplane/pkg/util"
	"github.com/go-logr/logr"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type operations interface {
	cloudsql.InstanceClient

	// Finalizer operations
	addFinalizer(ctx context.Context) error
	removeFinalizer(ctx context.Context) error

	// Object and Instance operations
	isInstanceReady() bool
	isReclaimDelete() bool
	getInstanceName() string
	getInstance(*sqladmin.DatabaseInstance) *sqladmin.DatabaseInstance
	needsUpdate(*sqladmin.DatabaseInstance) bool

	// Sync operation - save instance to object
	syncInstance(context.Context, *sqladmin.DatabaseInstance) error
	syncInstanceSpec(context.Context, *sqladmin.DatabaseInstance) error
	syncInstanceStatus(context.Context, *sqladmin.DatabaseInstance) error

	reconcileStatus(context.Context, error) error
}

type instanceOperations struct {
	*cloudsql.InstanceService
	client client.Client
	object *v1alpha1.Instance
	log    logr.Logger
}

var _ operations = &instanceOperations{}

func (o *instanceOperations) addFinalizer(ctx context.Context) error {
	util.AddFinalizer(&o.object.ObjectMeta, finalizer)
	return o.updateObject(ctx)
}

func (o *instanceOperations) getInstanceName() string {
	return o.object.DatabaseInstanceName()
}

func (o *instanceOperations) getInstance(inst *sqladmin.DatabaseInstance) *sqladmin.DatabaseInstance {
	return o.object.DatabaseInstance(inst)
}

func (o *instanceOperations) removeFinalizer(ctx context.Context) error {
	util.RemoveFinalizer(&o.object.ObjectMeta, finalizer)
	return o.updateObject(ctx)
}

func (o *instanceOperations) isInstanceReady() bool {
	return o.object.IsAvailable()
}

func (o *instanceOperations) isReclaimDelete() bool {
	return o.object.Spec.ReclaimPolicy == corev1alpha1.ReclaimDelete
}

func (o *instanceOperations) isStatusReady() bool {
	return o.object.Status.State == v1alpha1.StateRunnable
}

func (o *instanceOperations) needsUpdate(inst *sqladmin.DatabaseInstance) bool {
	diff := o.object.DatabaseInstanceSpecDiff(inst)
	o.log.WithName("diff").V(logging.Debug).Info(diff)
	return diff != ""
}

func (o *instanceOperations) reconcileStatus(ctx context.Context, err error) error {
	if err == nil {
		o.object.Status.SetConditions(corev1alpha1.ReconcileSuccess())
	} else {
		o.object.Status.SetConditions(corev1alpha1.ReconcileError(err))
	}
	return o.updateStatus(ctx)
}

func (o *instanceOperations) syncInstance(ctx context.Context, inst *sqladmin.DatabaseInstance) error {
	if err := o.syncInstanceSpec(ctx, inst); err != nil {
		return err
	}
	return o.syncInstanceStatus(ctx, inst)
}

func (o *instanceOperations) syncInstanceSpec(ctx context.Context, inst *sqladmin.DatabaseInstance) error {
	o.object.DatabaseInstanceSpec(inst)
	return o.updateObject(ctx)
}

func (o *instanceOperations) syncInstanceStatus(ctx context.Context, inst *sqladmin.DatabaseInstance) error {
	o.object.DatabaseInstanceStatus(inst)
	return o.updateStatus(ctx)
}

//
// Controller-runtime Client operations
//
func (o *instanceOperations) updateObject(ctx context.Context) error {
	if err := o.client.Update(ctx, o.object); err != nil {
		o.log.Error(err, "failed to update object", "version", o.object.ResourceVersion)
		return err
	}
	return nil
}

func (o *instanceOperations) updateStatus(ctx context.Context) error {
	if err := o.client.Status().Update(ctx, o.object); err != nil {
		o.log.Error(err, "failed to update status", "version", o.object.ResourceVersion)
		return err
	}
	return nil
}
