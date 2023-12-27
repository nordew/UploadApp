package service

import (
	"context"
	psqldb "github.com/nordew/UploadApp/internal/adapters/db/postgres"
	"github.com/nordew/UploadApp/internal/domain/entity"
)

type Dashboards interface {
	CreateLog(ctx context.Context, log *entity.AuditLog) error

	GetLogs(ctx context.Context) ([]entity.AuditLog, error)

	DeleteLog(ctx context.Context, id int64) error
}

type dashboardService struct {
	dashboardStorage psqldb.DashboardStorage
}

func NewDashboardService(dashboardStorage psqldb.DashboardStorage) *dashboardService {
	return &dashboardService{
		dashboardStorage: dashboardStorage,
	}
}

func (s *dashboardService) CreateLog(ctx context.Context, log *entity.AuditLog) error {
	return s.dashboardStorage.CreateLog(ctx, log)
}

func (s *dashboardService) GetLogs(ctx context.Context) ([]entity.AuditLog, error) {
	return s.dashboardStorage.GetLogs(ctx)
}

func (s *dashboardService) DeleteLog(ctx context.Context, id int64) error {
	return s.dashboardStorage.DeleteLog(ctx, id)
}
