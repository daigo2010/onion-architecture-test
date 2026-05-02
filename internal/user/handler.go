package user

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
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toResponse(u *User) response {
	return response{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type request struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}
	u, err := h.uc.Create(r.Context(), CreateInput{
		Name: req.Name, Email: req.Email,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, toResponse(u))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req request
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}
	u, err := h.uc.Update(r.Context(), UpdateInput{
		ID: id, Name: req.Name, Email: req.Email,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponse(u))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	u, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponse(u))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	us, err := h.uc.List(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	out := make([]response, 0, len(us))
	for _, u := range us {
		out = append(out, toResponse(u))
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
	case errors.Is(err, ErrEmptyName), errors.Is(err, ErrInvalidEmail):
		httpx.WriteError(w, http.StatusBadRequest, err)
	default:
		httpx.WriteError(w, http.StatusInternalServerError, err)
	}
}
