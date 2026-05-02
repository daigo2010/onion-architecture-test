package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"onion/domain/product"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS products (
    id         TEXT    PRIMARY KEY,
    name       TEXT    NOT NULL,
    price      INTEGER NOT NULL,
    stock      INTEGER NOT NULL,
    created_at TEXT    NOT NULL,
    updated_at TEXT    NOT NULL
);`

type SQLiteProductRepository struct {
	db *sql.DB
}

func NewSQLiteProductRepository(dsn string) (*SQLiteProductRepository, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &SQLiteProductRepository{db: db}, nil
}

func (r *SQLiteProductRepository) Close() error { return r.db.Close() }

func (r *SQLiteProductRepository) Save(ctx context.Context, p *product.Product) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO products (id, name, price, stock, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            name       = excluded.name,
            price      = excluded.price,
            stock      = excluded.stock,
            updated_at = excluded.updated_at
    `, p.ID, p.Name, p.Price, p.Stock,
		p.CreatedAt.UTC().Format(time.RFC3339Nano),
		p.UpdatedAt.UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (r *SQLiteProductRepository) FindByID(ctx context.Context, id string) (*product.Product, error) {
	row := r.db.QueryRowContext(ctx, `
        SELECT id, name, price, stock, created_at, updated_at
        FROM products WHERE id = ?
    `, id)
	p, err := scanProduct(row.Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, product.ErrNotFound
	}
	return p, err
}

func (r *SQLiteProductRepository) FindAll(ctx context.Context) ([]*product.Product, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, name, price, stock, created_at, updated_at
        FROM products ORDER BY created_at ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*product.Product, 0)
	for rows.Next() {
		p, err := scanProduct(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *SQLiteProductRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return product.ErrNotFound
	}
	return nil
}

func scanProduct(scan func(...any) error) (*product.Product, error) {
	var (
		p                       product.Product
		createdAt, updatedAt string
	)
	if err := scan(&p.ID, &p.Name, &p.Price, &p.Stock, &createdAt, &updatedAt); err != nil {
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
	p.CreatedAt = ca
	p.UpdatedAt = ua
	return &p, nil
}
