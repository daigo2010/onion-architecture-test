package product

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("product not found")

type Repository interface {
	Save(ctx context.Context, p *Product) error
	FindByID(ctx context.Context, id string) (*Product, error)
	FindAll(ctx context.Context) ([]*Product, error)
	Delete(ctx context.Context, id string) error
}
