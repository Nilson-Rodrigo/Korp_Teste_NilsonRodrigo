package http

import "net/http"

// CORSMiddleware adiciona os headers CORS necessários para acesso do frontend.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SetupRoutes configura todas as rotas da API usando Go 1.22+ patterns.
func SetupRoutes(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /api/health", handler.Health)

	// Rotas de produtos
	mux.HandleFunc("GET /api/produtos", handler.ListarProdutos)
	mux.HandleFunc("POST /api/produtos", handler.CriarProduto)
	mux.HandleFunc("GET /api/produtos/{id}", handler.BuscarProdutoPorID)
	mux.HandleFunc("PATCH /api/produtos/{id}/saldo", handler.AtualizarSaldo)

	// Rotas de notas fiscais
	mux.HandleFunc("GET /api/notas", handler.ListarNotas)
	mux.HandleFunc("POST /api/notas", handler.CriarNota)
	mux.HandleFunc("GET /api/notas/{id}", handler.BuscarNotaPorID)
	mux.HandleFunc("POST /api/notas/{id}/imprimir", handler.ImprimirNota)

	return CORSMiddleware(mux)
}
