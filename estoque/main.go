package main

import (
	"estoque/config"
	"estoque/handler"
	"estoque/model"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	config.ConnectDB()
	config.DB.AutoMigrate(&model.Produto{})

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/produtos", handler.ListarProdutos)
	r.POST("/produtos", handler.CriarProduto)
	r.GET("/produtos/:id", handler.BuscarProdutoPorID)
	r.PATCH("/produtos/:id/saldo", handler.AtualizarSaldo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}