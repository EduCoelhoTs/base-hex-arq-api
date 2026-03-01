package _http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockHttpHandler implementa a interface HttpHandler para testes
type MockHttpHandler struct {
	callCount int
	routes    Routes
}

func (m *MockHttpHandler) RegisterRoutes(routes Routes) http.Handler {
	m.callCount++
	m.routes = routes
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock response"))
	})
}

func TestNewHttpServer(t *testing.T) {
	tests := []struct {
		name           string
		handler        HttpHandler
		routes         Routes
		addr           *string
		expectedAddr   string
		expectedServer bool
	}{
		{
			name:         "should create server with provided address",
			handler:      &MockHttpHandler{},
			routes:       Routes{},
			addr:         StringPtr(":9090"),
			expectedAddr: ":9090",
		},
		{
			name:         "should use default address when nil",
			handler:      &MockHttpHandler{},
			routes:       Routes{},
			addr:         nil,
			expectedAddr: ":8080",
		},
		{
			name:         "should configure timeouts correctly",
			handler:      &MockHttpHandler{},
			routes:       Routes{},
			addr:         StringPtr(":3000"),
			expectedAddr: ":3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewHttpServer(tt.handler, tt.routes, tt.addr)

			if server == nil {
				t.Error("expected server to be not nil")
			}

			if server.svr == nil {
				t.Error("expected http.Server to be not nil")
			}

			if server.svr.Addr != tt.expectedAddr {
				t.Errorf("expected address %s, got %s", tt.expectedAddr, server.svr.Addr)
			}

			// Validar timeouts
			if server.svr.ReadTimeout != 30*time.Second {
				t.Errorf("expected ReadTimeout 30s, got %v", server.svr.ReadTimeout)
			}

			if server.svr.WriteTimeout != 30*time.Second {
				t.Errorf("expected WriteTimeout 30s, got %v", server.svr.WriteTimeout)
			}

			if server.svr.IdleTimeout != 60*time.Second {
				t.Errorf("expected IdleTimeout 60s, got %v", server.svr.IdleTimeout)
			}

			if server.svr.ReadHeaderTimeout != 15*time.Second {
				t.Errorf("expected ReadHeaderTimeout 15s, got %v", server.svr.ReadHeaderTimeout)
			}
		})
	}
}

func TestNewHttpServer_RegisterRoutes(t *testing.T) {
	mock := &MockHttpHandler{}
	routes := Routes{
		"/api": []Route{
			{
				Method:      http.MethodGet,
				Path:        "/users",
				HandlerFunc: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			},
		},
	}

	server := NewHttpServer(mock, routes, StringPtr(":8080"))

	// Validar que RegisterRoutes foi chamado
	if mock.callCount != 1 {
		t.Errorf("expected RegisterRoutes to be called once, got %d", mock.callCount)
	}

	// Validar que as rotas foram passadas corretamente
	if len(mock.routes) != len(routes) {
		t.Errorf("expected %d routes, got %d", len(routes), len(mock.routes))
	}

	if server.svr.Handler == nil {
		t.Error("expected handler to be set")
	}
}

func TestHttpServer_Shutdown(t *testing.T) {
	mock := &MockHttpHandler{}
	server := NewHttpServer(mock, Routes{}, StringPtr(":0"))

	// Não iniciamos o servidor, apenas testamos a chamada de Shutdown
	// Em um teste real, você faria: go server.Start()
	err := server.Shutdown()
	if err != nil {
		t.Logf("Shutdown error (esperado pois servidor não está rodando): %v", err)
	}

	if server.svr == nil {
		t.Error("expected server to still exist after shutdown")
	}
}

func TestHttpServer_HandlerFunctionality(t *testing.T) {
	mock := &MockHttpHandler{}
	server := NewHttpServer(mock, Routes{}, StringPtr(":8080"))

	// Usar httptest para validar o handler sem iniciar um servidor real
	if server.svr.Handler == nil {
		t.Fatal("handler should not be nil")
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	server.svr.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "mock response" {
		t.Errorf("expected body 'mock response', got %s", w.Body.String())
	}
}

// Helper para criar ponteiro de string
func StringPtr(s string) *string {
	return &s
}
