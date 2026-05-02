package main

import (
	"log"
	"net/http"
	"os"

	"onion/internal/app"
	"onion/internal/product"
	"onion/internal/shared/system"
	"onion/internal/user"
)

// 機能を増やすときはこのスライスに 1 行足すだけ。
var moduleFactories = []app.ModuleFactory{
	product.New,
	user.New,
}

func main() {
	deps := app.Deps{
		IDGen: system.IDGenerator{},
		Clock: system.Clock{},
	}

	mux := http.NewServeMux()
	for _, factory := range moduleFactories {
		m, err := factory(deps)
		if err != nil {
			log.Fatalf("init module: %v", err)
		}
		m.RegisterRoutes(mux)
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
