package _chi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	_http "github.com/EduCoelhoTs/nba-predict-api/internal/adapter/http"
	"github.com/stretchr/testify/assert"
)

// AuthMiddleware cria um middleware simples de autenticação
func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

// SimpleHandler é um handler simples para testes
func SimpleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

// NotFoundHandler é um handler que retorna 404
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

// TestNewChiHandler testa a criação de um novo chiHandler
func TestNewChiHandler(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "deve criar uma nova instância do chiHandler",
			test: func(t *testing.T) {
				h := NewChiHandler()

				assert.NotNil(t, h)
				assert.NotNil(t, h.handler)
			},
		},
		{
			name: "deve retornar diferentes instâncias",
			test: func(t *testing.T) {
				h1 := NewChiHandler()
				h2 := NewChiHandler()

				assert.NotEqual(t, h1, h2)
				assert.NotEqual(t, h1.handler, h2.handler)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestRegisterRoutes testa o registro de rotas
func TestRegisterRoutes(t *testing.T) {
	tests := []struct {
		name          string
		routes        _http.Routes
		testMethod    string
		testPath      string
		expectedCode  int
		expectedBody  string
	}{
		{
			name: "deve registrar uma rota GET simples",
			routes: _http.Routes{
				"/api": {
					{
						Method:      "GET",
						Path:        "/users",
						HandlerFunc: SimpleHandler,
						Middlewares: []func(h http.Handler) http.Handler{},
					},
				},
			},
			testMethod:   "GET",
			testPath:     "/api/users",
			expectedCode: http.StatusOK,
			expectedBody: "success",
		},
		{
			name: "deve registrar uma rota POST simples",
			routes: _http.Routes{
				"/api": {
					{
						Method:      "POST",
						Path:        "/users",
						HandlerFunc: SimpleHandler,
						Middlewares: []func(h http.Handler) http.Handler{},
					},
				},
			},
			testMethod:   "POST",
			testPath:     "/api/users",
			expectedCode: http.StatusOK,
			expectedBody: "success",
		},
		{
			name: "deve retornar 404 para rota não registrada",
			routes: _http.Routes{
				"/api": {
					{
						Method:      "GET",
						Path:        "/users",
						HandlerFunc: SimpleHandler,
						Middlewares: []func(h http.Handler) http.Handler{},
					},
				},
			},
			testMethod:   "GET",
			testPath:     "/api/products",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewChiHandler()
			handler := h.RegisterRoutes(tt.routes)

			req := httptest.NewRequest(tt.testMethod, tt.testPath, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

// TestRegisterRoutesWithMiddleware testa rotas com middlewares
func TestRegisterRoutesWithMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		routes       _http.Routes
		testMethod   string
		testPath     string
		expectedCode int
	}{
		{
			name: "deve aplicar middleware na rota",
			routes: _http.Routes{
				"/api": {
					{
						Method:      "GET",
						Path:        "/users",
						HandlerFunc: SimpleHandler,
						Middlewares: []func(h http.Handler) http.Handler{
							AuthMiddleware(),
						},
					},
				},
			},
			testMethod:   "GET",
			testPath:     "/api/users",
			expectedCode: http.StatusOK,
		},
		{
			name: "deve aplicar múltiplos middlewares em ordem",
			routes: _http.Routes{
				"/api": {
					{
						Method:      "GET",
						Path:        "/users",
						HandlerFunc: SimpleHandler,
						Middlewares: []func(h http.Handler) http.Handler{
							AuthMiddleware(),
							AuthMiddleware(),
						},
					},
				},
			},
			testMethod:   "GET",
			testPath:     "/api/users",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewChiHandler()
			handler := h.RegisterRoutes(tt.routes)

			req := httptest.NewRequest(tt.testMethod, tt.testPath, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

// TestRegisterMultipleRouteGroups testa múltiplos grupos de rotas
func TestRegisterMultipleRouteGroups(t *testing.T) {
	routes := _http.Routes{
		"/api/v1": {
			{
				Method:      "GET",
				Path:        "/users",
				HandlerFunc: SimpleHandler,
				Middlewares: []func(h http.Handler) http.Handler{},
			},
			{
				Method:      "POST",
				Path:        "/users",
				HandlerFunc: SimpleHandler,
				Middlewares: []func(h http.Handler) http.Handler{},
			},
		},
		"/api/v2": {
			{
				Method:      "GET",
				Path:        "/products",
				HandlerFunc: SimpleHandler,
				Middlewares: []func(h http.Handler) http.Handler{},
			},
		},
	}

	h := NewChiHandler()
	handler := h.RegisterRoutes(routes)

	testCases := []struct {
		method   string
		path     string
		expected int
	}{
		{"GET", "/api/v1/users", http.StatusOK},
		{"POST", "/api/v1/users", http.StatusOK},
		{"GET", "/api/v2/products", http.StatusOK},
		{"GET", "/api/v1/products", http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code)
		})
	}
}

// TestRegisterRoutesReturnsHttpHandler testa se RegisterRoutes retorna um http.Handler
func TestRegisterRoutesReturnsHttpHandler(t *testing.T) {
	h := NewChiHandler()
	routes := _http.Routes{
		"/api": {
			{
				Method:      "GET",
				Path:        "/test",
				HandlerFunc: SimpleHandler,
				Middlewares: []func(h http.Handler) http.Handler{},
			},
		},
	}

	handler := h.RegisterRoutes(routes)

	// Verificar que é um http.Handler
	_, ok := interface{}(handler).(http.Handler)
	assert.True(t, ok, "RegisterRoutes deve retornar um http.Handler")
}

// TestRegisterRoutesWithEmptyRoutes testa com rotas vazias
func TestRegisterRoutesWithEmptyRoutes(t *testing.T) {
	h := NewChiHandler()
	routes := _http.Routes{}

	handler := h.RegisterRoutes(routes)

	assert.NotNil(t, handler)

	req := httptest.NewRequest("GET", "/any-path", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestRegisterRoutesWithDifferentHTTPMethods testa diferentes métodos HTTP
func TestRegisterRoutesWithDifferentHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	routes := _http.Routes{}

	for _, method := range methods {
		routes["/api"] = append(routes["/api"], _http.Route{
			Method:      method,
			Path:        "/resource",
			HandlerFunc: SimpleHandler,
			Middlewares: []func(h http.Handler) http.Handler{},
		})
	}

	h := NewChiHandler()
	handler := h.RegisterRoutes(routes)

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/resource", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "success", w.Body.String())
		})
	}
}

// TestIntegrationRegisterRoutesAndServe testa integração completa
func TestIntegrationRegisterRoutesAndServe(t *testing.T) {
	routes := _http.Routes{
		"/api": {
			{
				Method:      "GET",
				Path:        "/health",
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					io.WriteString(w, `{"status":"healthy"}`)
				},
				Middlewares: []func(h http.Handler) http.Handler{},
			},
		},
	}

	h := NewChiHandler()
	handler := h.RegisterRoutes(routes)

	// Simular servidor HTTP
	server := httptest.NewServer(handler)
	defer server.Close()

	// Fazer requisição
	resp, err := http.Get(server.URL + "/api/health")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"status":"healthy"}`, string(body))
}
