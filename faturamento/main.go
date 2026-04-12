package main

import (
	"faturamento/config"
	"faturamento/handler"
	"faturamento/model"
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
	config.DB.AutoMigrate(&model.NotaFiscal{}, &model.ItemNota{})

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

	r.GET("/notas", handler.ListarNotas)
	r.POST("/notas", handler.CriarNota)
	r.GET("/notas/:id", handler.BuscarNotaPorID)
	r.POST("/notas/:id/imprimir", handler.ImprimirNota)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}