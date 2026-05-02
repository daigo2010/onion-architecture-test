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
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "products.db"
	}
	repo, err := persistence.NewSQLiteProductRepository(dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer repo.Close()

	uc := usecase.NewProductUseCase(repo, system.RandomIDGenerator{}, system.SystemClock{})
	h := handler.NewProductHandler(uc)
	r := router.New(h)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("listening on %s (db=%s)", addr, dsn)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
