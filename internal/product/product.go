package product

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEmptyName     = errors.New("name must not be empty")
	ErrNegativePrice = errors.New("price must be non-negative")
	ErrNegativeStock = errors.New("stock must be non-negative")
	ErrNotFound      = errors.New("product not found")
)

type Product struct {
	ID        string
	Name      string
	Price     int
	Stock     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func newProduct(id, name string, price, stock int, now time.Time) (*Product, error) {
	p := &Product{ID: id, CreatedAt: now}
	if err := p.apply(name, price, stock, now); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Product) update(name string, price, stock int, now time.Time) error {
	return p.apply(name, price, stock, now)
}

func (p *Product) apply(name string, price, stock int, now time.Time) error {
	if name == "" {
		return ErrEmptyName
	}
	if price < 0 {
		return ErrNegativePrice
	}
	if stock < 0 {
		return ErrNegativeStock
	}
	p.Name = name
	p.Price = price
	p.Stock = stock
	p.UpdatedAt = now
	return nil
}

type Repository interface {
	Save(ctx context.Context, p *Product) error
	FindByID(ctx context.Context, id string) (*Product, error)
	FindAll(ctx context.Context) ([]*Product, error)
	Delete(ctx context.Context, id string) error
}
