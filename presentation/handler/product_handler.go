package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"onion/domain/product"
	"onion/usecase"
)

type ProductHandler struct {
	uc *usecase.ProductUseCase
}

func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

type productResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toResponse(p *product.Product) productResponse {
	return productResponse{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type productRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req productRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	p, err := h.uc.Create(r.Context(), usecase.CreateProductInput{
		Name: req.Name, Price: req.Price, Stock: req.Stock,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(p))
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req productRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	p, err := h.uc.Update(r.Context(), usecase.UpdateProductInput{
		ID: id, Name: req.Name, Price: req.Price, Stock: req.Stock,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(p))
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	p, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(p))
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	ps, err := h.uc.List(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	out := make([]productResponse, 0, len(ps))
	for _, p := range ps {
		out = append(out, toResponse(p))
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func decodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, product.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, product.ErrEmptyName),
		errors.Is(err, product.ErrNegativePrice),
		errors.Is(err, product.ErrNegativeStock):
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}
