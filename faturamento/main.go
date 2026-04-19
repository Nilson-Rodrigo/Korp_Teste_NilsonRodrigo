package main

import (
	"faturamento/config"
	"faturamento/domain/entities"
	"faturamento/domain/usecases"
	"faturamento/infrastructure/http/handlers"
	"faturamento/infrastructure/persistence"
	"faturamento/infrastructure/services"
	"faturamento/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("Arquivo .env não encontrado")
	}

	// Inicializar logger
	utils.InitLogger()

	// Conectar ao banco de dados
	config.ConnectDB()
	config.DB.AutoMigrate(&entities.NotaFiscal{}, &entities.ItemNota{})

	// Configurar URL do estoque
	estoqueURL := os.Getenv("ESTOQUE_URL")
	if estoqueURL == "" {
		estoqueURL = "http://localhost:8080"
	}

	// Inicializar repositórios
	notaRepo := persistence.NewNotaFiscalRepository(config.DB)

	// Inicializar serviços
	estoqueService := services.NewEstoqueService(estoqueURL)

	// Inicializar use cases
	criarNotaUC := usecases.NewCriarNotaUseCase(notaRepo, estoqueService)
	listarNotasUC := usecases.NewListarNotasUseCase(notaRepo)
	buscarNotaUC := usecases.NewBuscarNotaPorIDUseCase(notaRepo)
	imprimirNotaUC := usecases.NewImprimirNotaUseCase(notaRepo, estoqueService, config.DB)

	// Inicializar handlers
	notaHandler := handlers.NewNotaFiscalHandler(criarNotaUC, listarNotasUC, buscarNotaUC, imprimirNotaUC)

	// Configurar router
	r := gin.Default()

	// Middleware CORS seguro
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:4200")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Rotas
	r.GET("/notas", notaHandler.Listar)
	r.POST("/notas", notaHandler.Criar)
	r.GET("/notas/:id", notaHandler.BuscarPorID)
	r.POST("/notas/:id/imprimir", notaHandler.Imprimir)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Info().Str("port", port).Msg("Iniciando servidor de faturamento")
	if err := r.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Erro ao iniciar servidor")
	}
}
