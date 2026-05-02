package user

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
	mux.HandleFunc("POST /users", m.handler.Create)
	mux.HandleFunc("GET /users", m.handler.List)
	mux.HandleFunc("GET /users/{id}", m.handler.Get)
	mux.HandleFunc("PUT /users/{id}", m.handler.Update)
	mux.HandleFunc("DELETE /users/{id}", m.handler.Delete)
}
