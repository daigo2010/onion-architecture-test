package handler

import (
	"errors"
	"net/http"
	"time"

	"onion/domain/user"
	"onion/usecase"
)

type UserHandler struct {
	uc *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

type userResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toUserResponse(u *user.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type userRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	u, err := h.uc.Create(r.Context(), usecase.CreateUserInput{
		Name: req.Name, Email: req.Email,
	})
	if err != nil {
		writeUserError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toUserResponse(u))
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req userRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	u, err := h.uc.Update(r.Context(), usecase.UpdateUserInput{
		ID: id, Name: req.Name, Email: req.Email,
	})
	if err != nil {
		writeUserError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(u))
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	u, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeUserError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(u))
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	us, err := h.uc.List(r.Context())
	if err != nil {
		writeUserError(w, err)
		return
	}
	out := make([]userResponse, 0, len(us))
	for _, u := range us {
		out = append(out, toUserResponse(u))
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeUserError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeUserError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, user.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, user.ErrEmptyName), errors.Is(err, user.ErrInvalidEmail):
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}
