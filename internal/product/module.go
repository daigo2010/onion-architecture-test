package product

import (
	"net/http"

	"onion/internal/app"
)

type Module struct {
	handler *Handler
}

func New(deps app.Deps) (app.Module, error) {
	repo := newInMemoryRepository()
	uc := NewUseCase(repo, deps.IDGen, deps.Clock)
	return &Module{handler: NewHandler(uc)}, nil
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /products", m.handler.Create)
	mux.HandleFunc("GET /products", m.handler.List)
	mux.HandleFunc("GET /products/{id}", m.handler.Get)
	mux.HandleFunc("PUT /products/{id}", m.handler.Update)
	mux.HandleFunc("DELETE /products/{id}", m.handler.Delete)
}
