package _http

import (
	"net/http"
	"time"
)

type httpServer struct {
	svr *http.Server
}

type HttpHandler interface {
	RegisterRoutes(routes Routes) http.Handler
}

func NewHttpServer(handler http.Handler, addr *string) *httpServer {

	if addr == nil {
		addr = new(string)
		*addr = ":8080"
	}

	return &httpServer{
		svr: &http.Server{
			Addr:              *addr,
			Handler:           handler,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 15 * time.Second,
		},
	}
}

func (s *httpServer) Start() error {
	return s.svr.ListenAndServe()
}

func (s *httpServer) Shutdown() error {
	return s.svr.Close()
}
