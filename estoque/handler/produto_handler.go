package handler

import (
	"estoque/model"
	"estoque/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListarProdutos(c *gin.Context) {
	produtos, err := repository.ListarProdutos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, produtos)
}

func CriarProduto(c *gin.Context) {
	var p model.Produto
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	if err := repository.CriarProduto(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func AtualizarSaldo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	var body struct {
		Saldo float64 `json:"saldo"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	if err := repository.AtualizarSaldo(uint(id), body.Saldo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensagem": "Saldo atualizado"})
}
func BuscarProdutoPorID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}
	p, err := repository.BuscarProdutoPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Produto não encontrado"})
		return
	}
	c.JSON(http.StatusOK, p)
}