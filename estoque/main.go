package main

import (
	"estoque/config"
	"estoque/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("Arquivo .env não encontrado, usando variáveis do sistema")
	}

	// Inicializar logger
	utils.InitLogger()

	// Conectar ao banco de dados
	config.ConnectDB()
	config.DB.AutoMigrate(&Produto{})

	// Inicializar handler
	handler := NewHandler(config.DB)

	// Configurar router
	r := gin.Default()

	// Middleware CORS
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
	r.GET("/produtos", handler.Listar)
	r.POST("/produtos", handler.Criar)
	r.GET("/produtos/:id", handler.BuscarPorID)
	r.PATCH("/produtos/:id/saldo", handler.AtualizarSaldo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Str("port", port).Msg("Iniciando servidor de estoque")
	if err := r.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Erro ao iniciar servidor")
	}
}
