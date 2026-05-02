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
		dsn = "app.db"
	}
	db, err := persistence.OpenSQLite(dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	productRepo, err := persistence.NewSQLiteProductRepository(db)
	if err != nil {
		log.Fatalf("init product repo: %v", err)
	}
	userRepo, err := persistence.NewSQLiteUserRepository(db)
	if err != nil {
		log.Fatalf("init user repo: %v", err)
	}

	idGen := system.RandomIDGenerator{}
	clock := system.SystemClock{}

	productUC := usecase.NewProductUseCase(productRepo, idGen, clock)
	userUC := usecase.NewUserUseCase(userRepo, idGen, clock)

	r := router.New(
		handler.NewProductHandler(productUC),
		handler.NewUserHandler(userUC),
	)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("listening on %s (db=%s)", addr, dsn)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
