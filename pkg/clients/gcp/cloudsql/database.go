package cloudsql

import (
	"context"

	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

type DatabaseClient interface {
	Get(ctx context.Context, name string) (*sqladmin.Database, error)
	Create(ctx context.Context, instance *sqladmin.Database) error
	Update(ctx context.Context, instance *sqladmin.Database) error
	Delete(ctx context.Context, name string) error
}

type DatabaseService struct {
	projectID string
	instance  string
	service   *sqladmin.DatabasesService
}

var _ DatabaseClient = &DatabaseService{}

func NewDatabaseService(projectID, instance string, service *sqladmin.Service) *DatabaseService {
	return &DatabaseService{
		projectID: projectID,
		instance:  instance,
		service:   sqladmin.NewDatabasesService(service),
	}
}

func (s *DatabaseService) Get(ctx context.Context, name string) (*sqladmin.Database, error) {
	return s.service.Get(s.projectID, s.instance, name).Context(ctx).Do()
}

func (s *DatabaseService) Create(ctx context.Context, db *sqladmin.Database) error {
	_, err := s.service.Insert(s.projectID, s.instance, db).Context(ctx).Do()
	return err
}

func (s *DatabaseService) Update(ctx context.Context, db *sqladmin.Database) error {
	_, err := s.service.Update(s.projectID, s.instance, db.Name, db).Context(ctx).Do()
	return err
}

func (s *DatabaseService) Delete(ctx context.Context, name string) error {
	if _, err := s.Get(ctx, name); err != nil {
		return err
	}
	_, err := s.service.Delete(s.projectID, s.instance, name).Context(ctx).Do()
	return err
}
