package usecase

import (
	"context"

	"onion/domain/user"
)

type UserUseCase struct {
	repo  user.Repository
	idGen IDGenerator
	clock Clock
}

func NewUserUseCase(repo user.Repository, idGen IDGenerator, clock Clock) *UserUseCase {
	return &UserUseCase{repo: repo, idGen: idGen, clock: clock}
}

type CreateUserInput struct {
	Name  string
	Email string
}

func (u *UserUseCase) Create(ctx context.Context, in CreateUserInput) (*user.User, error) {
	usr, err := user.New(u.idGen.NewID(), in.Name, in.Email, u.clock.Now())
	if err != nil {
		return nil, err
	}
	if err := u.repo.Save(ctx, usr); err != nil {
		return nil, err
	}
	return usr, nil
}

type UpdateUserInput struct {
	ID    string
	Name  string
	Email string
}

func (u *UserUseCase) Update(ctx context.Context, in UpdateUserInput) (*user.User, error) {
	usr, err := u.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if err := usr.Update(in.Name, in.Email, u.clock.Now()); err != nil {
		return nil, err
	}
	if err := u.repo.Save(ctx, usr); err != nil {
		return nil, err
	}
	return usr, nil
}

func (u *UserUseCase) Get(ctx context.Context, id string) (*user.User, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *UserUseCase) List(ctx context.Context) ([]*user.User, error) {
	return u.repo.FindAll(ctx)
}

func (u *UserUseCase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
