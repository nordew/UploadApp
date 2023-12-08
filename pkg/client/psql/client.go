package psql

import (
	"database/sql"
	"fmt"
)

type ConnectionInfo struct {
	Host     string
	Port     string
	User     string
	DBName   string
	SSLMode  string
	Password string
}

func NewPostgres(connInfo *ConnectionInfo) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		connInfo.Host, connInfo.Port, connInfo.User, connInfo.DBName, connInfo.SSLMode, connInfo.Password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
