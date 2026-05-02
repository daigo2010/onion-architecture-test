package product

import (
	"context"
	"sort"
	"sync"
)

type inMemoryRepository struct {
	mu    sync.RWMutex
	store map[string]Product
}

func newInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{store: make(map[string]Product)}
}

func (r *inMemoryRepository) Save(_ context.Context, p *Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[p.ID] = *p
	return nil
}

func (r *inMemoryRepository) FindByID(_ context.Context, id string) (*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := p
	return &cp, nil
}

func (r *inMemoryRepository) FindAll(_ context.Context) ([]*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Product, 0, len(r.store))
	for _, p := range r.store {
		cp := p
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
