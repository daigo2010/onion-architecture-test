package user

import (
	"context"
	"sort"
	"sync"
)

type inMemoryRepository struct {
	mu    sync.RWMutex
	store map[string]User
}

func newInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{store: make(map[string]User)}
}

func (r *inMemoryRepository) Save(_ context.Context, u *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[u.ID] = *u
	return nil
}

func (r *inMemoryRepository) FindByID(_ context.Context, id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := u
	return &cp, nil
}

func (r *inMemoryRepository) FindAll(_ context.Context) ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*User, 0, len(r.store))
	for _, u := range r.store {
		cp := u
		out = append(out, &cp)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

func (r *inMemoryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}
