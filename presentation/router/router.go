package router

import (
	"net/http"

	"onion/presentation/handler"
)

func New(h *handler.ProductHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /products", h.Create)
	mux.HandleFunc("GET /products", h.List)
	mux.HandleFunc("GET /products/{id}", h.Get)
	mux.HandleFunc("PUT /products/{id}", h.Update)
	mux.HandleFunc("DELETE /products/{id}", h.Delete)
	return mux
}
