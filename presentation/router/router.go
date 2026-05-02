package router

import (
	"net/http"

	"onion/presentation/handler"
)

func New(ph *handler.ProductHandler, uh *handler.UserHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /products", ph.Create)
	mux.HandleFunc("GET /products", ph.List)
	mux.HandleFunc("GET /products/{id}", ph.Get)
	mux.HandleFunc("PUT /products/{id}", ph.Update)
	mux.HandleFunc("DELETE /products/{id}", ph.Delete)

	mux.HandleFunc("POST /users", uh.Create)
	mux.HandleFunc("GET /users", uh.List)
	mux.HandleFunc("GET /users/{id}", uh.Get)
	mux.HandleFunc("PUT /users/{id}", uh.Update)
	mux.HandleFunc("DELETE /users/{id}", uh.Delete)

	return mux
}
