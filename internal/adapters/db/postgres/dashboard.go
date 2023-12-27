package psqldb

import (
	"context"
	"database/sql"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/sirupsen/logrus"
)

// DashboardStorage is an interface that defines methods for interacting with the
// database to manage audit logs in a dashboard.
type DashboardStorage interface {
	// CreateLog creates a new audit log entry in the database.
	// It takes an AuditLog entity as input and returns an error if the operation fails.
	CreateLog(ctx context.Context, log *entity.AuditLog) error

	// GetLogs retrieves all audit logs from the database.
	// It returns a slice of AuditLog entities and an error if the retrieval fails.
	GetLogs(ctx context.Context) ([]entity.AuditLog, error)

	// DeleteLog deletes an audit log entry from the database based on the provided ID.
	// It returns an error if the deletion operation fails.
	DeleteLog(ctx context.Context, id int64) error
}

type dashboardStorage struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewDashboardStorage(db *sql.DB, logger *logrus.Logger) *dashboardStorage {
	return &dashboardStorage{
		db:     db,
		logger: logger,
	}
}

func (d *dashboardStorage) CreateLog(ctx context.Context, log *entity.AuditLog) error {
	logger := d.logger.WithField("function", "CreateLog")

	_, err := d.db.ExecContext(ctx,
		"INSERT INTO audit_logs (user_id, action_type, timestamp) VALUES ($1, $2, $3)",
		log.UserID, log.ActionType, log.Timestamp)
	if err != nil {
		logger.WithError(err).Error("failed to create log")
		return err
	}

	logger.Info("CreateLog: log created successfully")
	return nil
}

func (d *dashboardStorage) GetLogs(ctx context.Context) ([]entity.AuditLog, error) {
	logger := d.logger.WithField("function", "GetLogs")

	var logs []entity.AuditLog

	rows, err := d.db.QueryContext(ctx, "SELECT * FROM audit_logs")
	if err != nil {
		logger.WithError(err).Error("failed to retrieve logs")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var log entity.AuditLog
		if err := rows.Scan(&log.LogID, &log.UserID, &log.ActionType, &log.Timestamp); err != nil {
			logger.WithError(err).Error("failed to scan log")
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		logger.WithError(err).Error("error while iterating over logs")
		return nil, err
	}

	logger.Info("GetLogs: logs retrieved successfully")
	return logs, nil
}

func (d *dashboardStorage) DeleteLog(ctx context.Context, id int64) error {
	logger := d.logger.WithField("function", "DeleteLog")

	_, err := d.db.ExecContext(ctx, "DELETE FROM audit_logs WHERE id = $1", id)
	if err != nil {
		logger.WithError(err).Error("failed to delete log")
		return err
	}

	logger.Info("DeleteLog: log deleted successfully")
	return nil
}
