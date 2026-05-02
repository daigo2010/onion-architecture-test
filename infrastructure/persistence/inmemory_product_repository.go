package persistence

import (
	"context"
	"sort"
	"sync"

	"onion/domain/product"
)

type InMemoryProductRepository struct {
	mu    sync.RWMutex
	store map[string]product.Product
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{store: make(map[string]product.Product)}
}

func (r *InMemoryProductRepository) Save(_ context.Context, p *product.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[p.ID] = *p
	return nil
}

func (r *InMemoryProductRepository) FindByID(_ context.Context, id string) (*product.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[id]
	if !ok {
		return nil, product.ErrNotFound
	}
	cp := p
	return &cp, nil
}

func (r *InMemoryProductRepository) FindAll(_ context.Context) ([]*product.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*product.Product, 0, len(r.store))
	for _, p := range r.store {
		cp := p
		out = append(out, &cp)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

func (r *InMemoryProductRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return product.ErrNotFound
	}
	delete(r.store, id)
	return nil
}
