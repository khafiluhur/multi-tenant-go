package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

type User struct {
	ID           int        `db:"id"`
	TenantID     int        `db:"tenant_id"`
	Username     string     `db:"username"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	LastLogin    *time.Time `db:"last_login"`
}

func CreateUser(db *pgx.Conn, user *User) error {
	err := db.QueryRow(context.Background(),
		"INSERT INTO users (username, tenant_id, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id",
		user.Username, user.TenantID, user.Email, user.PasswordHash).Scan(&user.ID)
	return err
}

func SoftDeleteUser(db *pgx.Conn, userID int) error {
	_, err := db.Exec(context.Background(), "UPDATE users SET deleted_at = NOW() WHERE id = $1", userID)
	return err
}

func UpdateUser(db *pgx.Conn, user *User) error {
	_, err := db.Exec(context.Background(),
		"UPDATE users SET username = $1, email = $2, password_hash = $3, updated_at = NOW() WHERE id = $4",
		user.Username, user.Email, user.PasswordHash, user.ID)
	return err
}

func SetLastLogin(db *pgx.Conn, userID int) error {
	_, err := db.Exec(context.Background(), "UPDATE users SET last_login = NOW() WHERE id = $1", userID)
	return err
}
