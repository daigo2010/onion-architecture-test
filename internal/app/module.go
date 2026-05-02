package app

import (
	"net/http"
	"time"
)

type IDGenerator interface {
	NewID() string
}

type Clock interface {
	Now() time.Time
}

type Deps struct {
	IDGen IDGenerator
	Clock Clock
}

type Module interface {
	RegisterRoutes(mux *http.ServeMux)
}

type ModuleFactory func(Deps) (Module, error)
