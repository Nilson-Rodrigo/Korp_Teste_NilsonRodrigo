package handler

import (
	"bytes"
	"encoding/json"
	"faturamento/config"
	"faturamento/model"
	"faturamento/repository"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ListarNotas(c *gin.Context) {
	notas, err := repository.ListarNotas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notas)
}

func CriarNota(c *gin.Context) {
	var n model.NotaFiscal
	if err := c.ShouldBindJSON(&n); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	if err := repository.CriarNota(&n); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, n)
}

func ImprimirNota(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	nota, err := repository.BuscarNotaPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Nota não encontrada"})
		return
	}

	if nota.Status != "Aberta" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Apenas notas com status Aberta podem ser impressas"})
		return
	}

	estoqueURL := os.Getenv("ESTOQUE_URL")
	client := &http.Client{Timeout: 5 * time.Second}

	for _, item := range nota.Itens {
		resp, err := client.Get(fmt.Sprintf("%s/produtos/%d", estoqueURL, item.ProdutoID))
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "Serviço de estoque indisponível"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"erro": fmt.Sprintf("Produto ID %d não encontrado no estoque", item.ProdutoID)})
			return
		}

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "Erro ao buscar produto no estoque"})
			return
		}

		var produto struct {
			ID    uint    `json:"ID"`
			Saldo float64 `json:"saldo"`
		}
		json.NewDecoder(resp.Body).Decode(&produto)

		novoSaldo := produto.Saldo - item.Quantidade

		body, _ := json.Marshal(map[string]float64{"saldo": novoSaldo})
		reqURL := fmt.Sprintf("%s/produtos/%d/saldo", estoqueURL, item.ProdutoID)
		reqPatch, _ := http.NewRequest(http.MethodPatch, reqURL, bytes.NewBuffer(body))
		reqPatch.Header.Set("Content-Type", "application/json")

		respPatch, err := client.Do(reqPatch)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "Falha ao atualizar saldo no estoque"})
			return
		}
		respPatch.Body.Close()

		if respPatch.StatusCode != http.StatusOK {
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "Falha ao atualizar saldo no estoque"})
			return
		}
	}

	if err := repository.FecharNota(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "Nota impressa e fechada com sucesso"})
}

func BuscarNotaPorID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}
	nota, err := repository.BuscarNotaPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Nota não encontrada"})
		return
	}
	c.JSON(http.StatusOK, nota)
}

func init() {
	_ = config.DB
}