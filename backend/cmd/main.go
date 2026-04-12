package main

import (
	"backend/internal/application"
	apphttp "backend/internal/infrastructure/http"
	"backend/internal/infrastructure/repository"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Banner de inicialização
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║       KORP ERP - Backend API         ║")
	fmt.Println("║       Clean Architecture in Go       ║")
	fmt.Println("╚══════════════════════════════════════╝")

	// Inicializar repositórios em memória
	produtoRepo := repository.NovoProdutoRepositoryMemoria()
	notaRepo := repository.NovaNotaFiscalRepositoryMemoria()

	// Inicializar serviços com injeção de dependências
	produtoService := application.NovoProdutoService(produtoRepo)
	notaService := application.NovaNotaFiscalService(notaRepo, produtoRepo)

	// Seed: dados iniciais de produtos
	seedProdutos(produtoService)

	// Inicializar handlers HTTP
	handler := apphttp.NovoHandler(produtoService, notaService)

	// Configurar rotas
	routes := apphttp.SetupRoutes(handler)

	// Configurar servidor HTTP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      routes,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("🚀 Servidor iniciado na porta %s", port)
		log.Printf("📡 API disponível em http://localhost:%s/api", port)
		log.Printf("❤️  Health check: http://localhost:%s/api/health", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Desligando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro no shutdown do servidor: %v", err)
	}

	log.Println("✅ Servidor desligado com sucesso")
}

// seedProdutos cria produtos iniciais para demonstração.
func seedProdutos(service *application.ProdutoServiceImpl) {
	produtos := []application.CriarProdutoInput{
		{Codigo: "PROD-001", Descricao: "Parafuso Sextavado M8x30", Saldo: 500},
		{Codigo: "PROD-002", Descricao: "Porca Sextavada M8", Saldo: 1000},
		{Codigo: "PROD-003", Descricao: "Arruela Lisa 8mm", Saldo: 2000},
		{Codigo: "PROD-004", Descricao: "Barra Roscada M10x1m", Saldo: 50},
	}

	for _, p := range produtos {
		_, err := service.CriarProduto(p)
		if err != nil {
			log.Printf("⚠️  Erro ao criar produto seed %s: %v", p.Codigo, err)
		} else {
			log.Printf("📦 Produto seed criado: %s - %s", p.Codigo, p.Descricao)
		}
	}
}
