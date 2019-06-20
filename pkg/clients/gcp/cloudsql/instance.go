package cloudsql

import (
	"context"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

type InstanceClient interface {
	Create(ctx context.Context, instance *sqladmin.DatabaseInstance) error
	Delete(ctx context.Context, name string) error
	Get(ctx context.Context, name string) (*sqladmin.DatabaseInstance, error)
	Patch(ctx context.Context, instance *sqladmin.DatabaseInstance) error
	Update(ctx context.Context, instance *sqladmin.DatabaseInstance) error
}

type InstanceService struct {
	projectID string
	service   *sqladmin.InstancesService
}

var _ InstanceClient = &InstanceService{}

func NewInstanceService(projectID string, service *sqladmin.Service) *InstanceService {
	return &InstanceService{
		projectID: projectID,
		service:   sqladmin.NewInstancesService(service),
	}
}

func (s *InstanceService) Get(ctx context.Context, name string) (*sqladmin.DatabaseInstance, error) {
	return s.service.Get(s.projectID, name).Do()
}

func (s *InstanceService) Create(ctx context.Context, instance *sqladmin.DatabaseInstance) error {
	_, err := s.service.Insert(s.projectID, instance).Context(ctx).Do()
	return err
}

func (s *InstanceService) Patch(ctx context.Context, instance *sqladmin.DatabaseInstance) error {
	_, err := s.service.Patch(s.projectID, instance.Name, instance).Context(ctx).Do()
	return err
}

func (s *InstanceService) Update(ctx context.Context, instance *sqladmin.DatabaseInstance) error {
	_, err := s.service.Update(s.projectID, instance.Name, instance).Context(ctx).Do()
	return err
}

func (s *InstanceService) Delete(ctx context.Context, name string) error {
	if _, err := s.Get(ctx, name); err != nil {
		return err
	}
	_, err := s.service.Delete(s.projectID, name).Context(ctx).Do()
	return err
}
