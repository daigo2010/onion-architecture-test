package user

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

const sqliteSchema = `
CREATE TABLE IF NOT EXISTS users (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);`

type sqliteRepository struct {
	db *sql.DB
}

func newSQLiteRepository(db *sql.DB) (*sqliteRepository, error) {
	if _, err := db.Exec(sqliteSchema); err != nil {
		return nil, err
	}
	return &sqliteRepository{db: db}, nil
}

func (r *sqliteRepository) Save(ctx context.Context, u *User) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO users (id, name, email, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            name       = excluded.name,
            email      = excluded.email,
            updated_at = excluded.updated_at
    `, u.ID, u.Name, u.Email,
		u.CreatedAt.UTC().Format(time.RFC3339Nano),
		u.UpdatedAt.UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (r *sqliteRepository) FindByID(ctx context.Context, id string) (*User, error) {
	row := r.db.QueryRowContext(ctx, `
        SELECT id, name, email, created_at, updated_at
        FROM users WHERE id = ?
    `, id)
	u, err := scanRow(row.Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

func (r *sqliteRepository) FindAll(ctx context.Context) ([]*User, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, name, email, created_at, updated_at
        FROM users ORDER BY created_at ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*User, 0)
	for rows.Next() {
		u, err := scanRow(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (r *sqliteRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func scanRow(scan func(...any) error) (*User, error) {
	var (
		u                    User
		createdAt, updatedAt string
	)
	if err := scan(&u.ID, &u.Name, &u.Email, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	ca, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, err
	}
	ua, err := time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return nil, err
	}
	u.CreatedAt = ca
	u.UpdatedAt = ua
	return &u, nil
}
