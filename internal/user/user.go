package user

import (
	"context"
	"errors"
	"net/mail"
	"time"
)

var (
	ErrEmptyName    = errors.New("name must not be empty")
	ErrInvalidEmail = errors.New("email is invalid")
	ErrNotFound     = errors.New("user not found")
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func newUser(id, name, email string, now time.Time) (*User, error) {
	u := &User{ID: id, CreatedAt: now}
	if err := u.apply(name, email, now); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) update(name, email string, now time.Time) error {
	return u.apply(name, email, now)
}

func (u *User) apply(name, email string, now time.Time) error {
	if name == "" {
		return ErrEmptyName
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}
	u.Name = name
	u.Email = email
	u.UpdatedAt = now
	return nil
}

type Repository interface {
	Save(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindAll(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, id string) error
}
