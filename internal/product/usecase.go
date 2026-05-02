package product

import (
	"context"

	"onion/internal/app"
)

type UseCase struct {
	repo  Repository
	idGen app.IDGenerator
	clock app.Clock
}

func NewUseCase(repo Repository, idGen app.IDGenerator, clock app.Clock) *UseCase {
	return &UseCase{repo: repo, idGen: idGen, clock: clock}
}

type CreateInput struct {
	Name  string
	Price int
	Stock int
}

func (u *UseCase) Create(ctx context.Context, in CreateInput) (*Product, error) {
	p, err := newProduct(u.idGen.NewID(), in.Name, in.Price, in.Stock, u.clock.Now())
	if err != nil {
		return nil, err
	}
	if err := u.repo.Save(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

type UpdateInput struct {
	ID    string
	Name  string
	Price int
	Stock int
}

func (u *UseCase) Update(ctx context.Context, in UpdateInput) (*Product, error) {
	p, err := u.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if err := p.update(in.Name, in.Price, in.Stock, u.clock.Now()); err != nil {
		return nil, err
	}
	if err := u.repo.Save(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (u *UseCase) Get(ctx context.Context, id string) (*Product, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *UseCase) List(ctx context.Context) ([]*Product, error) {
	return u.repo.FindAll(ctx)
}

func (u *UseCase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
