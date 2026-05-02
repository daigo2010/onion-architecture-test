package usecase

import (
	"context"
	"time"

	"onion/domain/product"
)

type IDGenerator interface {
	NewID() string
}

type Clock interface {
	Now() time.Time
}

type ProductUseCase struct {
	repo  product.Repository
	idGen IDGenerator
	clock Clock
}

func NewProductUseCase(repo product.Repository, idGen IDGenerator, clock Clock) *ProductUseCase {
	return &ProductUseCase{repo: repo, idGen: idGen, clock: clock}
}

type CreateProductInput struct {
	Name  string
	Price int
	Stock int
}

func (u *ProductUseCase) Create(ctx context.Context, in CreateProductInput) (*product.Product, error) {
	p, err := product.New(u.idGen.NewID(), in.Name, in.Price, in.Stock, u.clock.Now())
	if err != nil {
		return nil, err
	}
	if err := u.repo.Save(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

type UpdateProductInput struct {
	ID    string
	Name  string
	Price int
	Stock int
}

func (u *ProductUseCase) Update(ctx context.Context, in UpdateProductInput) (*product.Product, error) {
	p, err := u.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if err := p.Update(in.Name, in.Price, in.Stock, u.clock.Now()); err != nil {
		return nil, err
	}
	if err := u.repo.Save(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (u *ProductUseCase) Get(ctx context.Context, id string) (*product.Product, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *ProductUseCase) List(ctx context.Context) ([]*product.Product, error) {
	return u.repo.FindAll(ctx)
}

func (u *ProductUseCase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
