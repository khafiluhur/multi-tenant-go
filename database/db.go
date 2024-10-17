package database

import (
	"context"
	"jatis_mobile_api/logs"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

var (
	conn   *pgx.Conn
	logger = logs.SetupLogger()
)

func ConnectDB(postgresURL string) error {
	var err error
	conn, err = pgx.Connect(context.Background(), postgresURL)
	if err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to connect to PostgreSQL", struct {
			PostgresURL string
			Error       error
		}{PostgresURL: postgresURL, Error: err})
		return err
	}
	logs.LogWithFields(logger, logrus.InfoLevel, "Connected to PostgreSQL", struct{ PostgresURL string }{PostgresURL: postgresURL})
	return nil
}

func GetDB() *pgx.Conn {
	return conn
}

func Close() {
	if conn != nil {
		if err := conn.Close(context.Background()); err != nil {
			logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to close PostgreSQL connection", struct{ Error error }{Error: err})
		} else {
			logs.LogWithFields(logger, logrus.InfoLevel, "PostgreSQL connection closed successfully", struct{}{})
		}
	}
}
