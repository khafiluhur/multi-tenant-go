package migrations

import (
	"context"

	"jatis_mobile_api/logs"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

func CreateTenantsTable(db *pgx.Conn, logger *logrus.Logger) error {
	query := `
    CREATE TABLE IF NOT EXISTS tenants (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT now(),
        deleted_at TIMESTAMP NULL
    );
    `
	_, err := db.Exec(context.Background(), query)
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Unable to create tenants table", struct{ Error error }{Error: err})
		return err
	}

	logs.LogWithFields(logger, logrus.InfoLevel, "Tenants table created successfully", struct{}{})
	return nil
}
