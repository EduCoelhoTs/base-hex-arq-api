package main

import (
	"fmt"
	"net/http"

	_http "github.com/EduCoelhoTs/nba-predict-api/internal/adapter/http"
	_chi "github.com/EduCoelhoTs/nba-predict-api/internal/adapter/http/chi"
)

func main() {
	addr := ":8080"

	handler := _chi.NewChiHandler()

	routes := _http.Routes{
		"/teste": {
			{
				Method: http.MethodGet,
				Path:   "/teste",
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Hello, World!"))
				},
				Middlewares: []func(h http.Handler) http.Handler{},
			},
		},
	}

	httpServer := _http.NewHttpServer(handler, routes, &addr)

	if err := httpServer.Start(); err != nil {
		panic(err)
	}

	fmt.Printf("Server is running on %s\n", addr)
}
