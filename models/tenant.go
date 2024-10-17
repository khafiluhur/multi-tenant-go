package models

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Tenant struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func CreateTenant(db *pgx.Conn, tenant *Tenant) error {
	err := db.QueryRow(context.Background(), "INSERT INTO tenants (name) VALUES ($1) RETURNING id", tenant.Name).Scan(&tenant.ID)
	return err
}

func SoftDeleteTenant(db *pgx.Conn, tenantID int) error {
	_, err := db.Exec(context.Background(), "UPDATE tenants SET deleted_at = NOW() WHERE id = $1", tenantID)
	return err
}
