package product

import (
	"errors"
	"net/http"
	"time"

	"onion/internal/shared/httpx"
)

type Handler struct {
	uc *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{uc: uc}
}

type response struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toResponse(p *Product) response {
	return response{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type request struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}
	p, err := h.uc.Create(r.Context(), CreateInput{
		Name: req.Name, Price: req.Price, Stock: req.Stock,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, toResponse(p))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req request
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}
	p, err := h.uc.Update(r.Context(), UpdateInput{
		ID: id, Name: req.Name, Price: req.Price, Stock: req.Stock,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponse(p))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	p, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponse(p))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ps, err := h.uc.List(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	out := make([]response, 0, len(ps))
	for _, p := range ps {
		out = append(out, toResponse(p))
	}
	httpx.WriteJSON(w, http.StatusOK, out)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.WriteError(w, http.StatusNotFound, err)
	case errors.Is(err, ErrEmptyName),
		errors.Is(err, ErrNegativePrice),
		errors.Is(err, ErrNegativeStock):
		httpx.WriteError(w, http.StatusBadRequest, err)
	default:
		httpx.WriteError(w, http.StatusInternalServerError, err)
	}
}
