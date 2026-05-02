package user

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
	Email string
}

func (u *UseCase) Create(ctx context.Context, in CreateInput) (*User, error) {
	usr, err := newUser(u.idGen.NewID(), in.Name, in.Email, u.clock.Now())
	if err != nil {
		return nil, err
	}
	if err := u.repo.Save(ctx, usr); err != nil {
		return nil, err
	}
	return usr, nil
}

type UpdateInput struct {
	ID    string
	Name  string
	Email string
}

func (u *UseCase) Update(ctx context.Context, in UpdateInput) (*User, error) {
	usr, err := u.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if err := usr.update(in.Name, in.Email, u.clock.Now()); err != nil {
		return nil, err
	}
	if err := u.repo.Save(ctx, usr); err != nil {
		return nil, err
	}
	return usr, nil
}

func (u *UseCase) Get(ctx context.Context, id string) (*User, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *UseCase) List(ctx context.Context) ([]*User, error) {
	return u.repo.FindAll(ctx)
}

func (u *UseCase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
