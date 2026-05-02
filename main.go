package main

import (
	"log"
	"net/http"
	"os"

	"onion/infrastructure/persistence"
	"onion/infrastructure/system"
	"onion/presentation/handler"
	"onion/presentation/router"
	"onion/usecase"
)

func main() {
	repo := persistence.NewInMemoryProductRepository()
	uc := usecase.NewProductUseCase(repo, system.RandomIDGenerator{}, system.SystemClock{})
	h := handler.NewProductHandler(uc)
	r := router.New(h)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
