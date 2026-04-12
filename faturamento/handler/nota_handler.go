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

type estoqueProduto struct {
	ID    uint    `json:"ID"`
	Saldo float64 `json:"saldo"`
}

func buscarProdutoEstoque(estoqueURL string, id uint) (*estoqueProduto, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/produtos/%d", estoqueURL, id)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("[ERRO] Falha ao buscar produto %d do estoque: %v\n", id, err)
		return nil, fmt.Errorf("serviço de estoque indisponível")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[AVISO] Produto %d retornou status %d do estoque\n", id, resp.StatusCode)
		return nil, fmt.Errorf("produto não encontrado ou serviço indisponível")
	}

	var produto estoqueProduto
	if err := json.NewDecoder(resp.Body).Decode(&produto); err != nil {
		fmt.Printf("[ERRO] Erro ao desserializar produto %d: %v\n", id, err)
		return nil, fmt.Errorf("erro ao processar dados do produto")
	}
	return &produto, nil
}

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

	if len(n.Itens) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "A nota fiscal precisa conter ao menos um item"})
		return
	}

	for _, item := range n.Itens {
		if item.Quantidade <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"erro": "A quantidade de cada item deve ser maior que zero"})
			return
		}
	}

	estoqueURL := os.Getenv("ESTOQUE_URL")
	if estoqueURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Variável ESTOQUE_URL não configurada"})
		return
	}

	for _, item := range n.Itens {
		if _, err := buscarProdutoEstoque(estoqueURL, item.ProdutoID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erro": fmt.Sprintf("Produto ID %d inválido ou indisponível", item.ProdutoID)})
			return
		}
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
	if estoqueURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Variável ESTOQUE_URL não configurada"})
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	updates := make([]struct {
		ProdutoID   uint
		SaldoAntigo float64
		SaldoNovo   float64
	}, 0, len(nota.Itens))

	for _, item := range nota.Itens {
		produto, err := buscarProdutoEstoque(estoqueURL, item.ProdutoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erro": fmt.Sprintf("Produto ID %d inválido ou indisponível", item.ProdutoID)})
			return
		}

		if produto.Saldo < item.Quantidade {
			c.JSON(http.StatusBadRequest, gin.H{"erro": fmt.Sprintf("Saldo insuficiente para o produto ID %d", item.ProdutoID)})
			return
		}

		updates = append(updates, struct {
			ProdutoID   uint
			SaldoAntigo float64
			SaldoNovo   float64
		}{ProdutoID: item.ProdutoID, SaldoAntigo: produto.Saldo, SaldoNovo: produto.Saldo - item.Quantidade})
	}

	successful := make([]struct {
		ProdutoID   uint
		SaldoAntigo float64
	}, 0, len(updates))

	for _, update := range updates {
		body, _ := json.Marshal(map[string]float64{"saldo": update.SaldoNovo})
		reqURL := fmt.Sprintf("%s/produtos/%d/saldo", estoqueURL, update.ProdutoID)
		reqPatch, _ := http.NewRequest(http.MethodPatch, reqURL, bytes.NewBuffer(body))
		reqPatch.Header.Set("Content-Type", "application/json")

		respPatch, err := client.Do(reqPatch)
		if err != nil {
			rollbackEstoque(estoqueURL, successful)
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "Falha ao atualizar saldo no estoque"})
			return
		}
		respPatch.Body.Close()

		if respPatch.StatusCode != http.StatusOK {
			rollbackEstoque(estoqueURL, successful)
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "Falha ao atualizar saldo no estoque"})
			return
		}

		successful = append(successful, struct {
			ProdutoID   uint
			SaldoAntigo float64
		}{ProdutoID: update.ProdutoID, SaldoAntigo: update.SaldoAntigo})
	}

	if err := repository.FecharNota(uint(id)); err != nil {
		rollbackEstoque(estoqueURL, successful)
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "Nota impressa e fechada com sucesso"})
}

func rollbackEstoque(estoqueURL string, updates []struct {
	ProdutoID   uint
	SaldoAntigo float64
}) {
	client := &http.Client{Timeout: 5 * time.Second}
	for _, update := range updates {
		body, err := json.Marshal(map[string]float64{"saldo": update.SaldoAntigo})
		if err != nil {
			fmt.Printf("[ERRO] Rollback falhou ao serializar saldo para produto %d: %v\n", update.ProdutoID, err)
			continue
		}
		
		reqURL := fmt.Sprintf("%s/produtos/%d/saldo", estoqueURL, update.ProdutoID)
		reqPatch, err := http.NewRequest(http.MethodPatch, reqURL, bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("[ERRO] Rollback falhou ao criar request para produto %d: %v\n", update.ProdutoID, err)
			continue
		}
		
		reqPatch.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(reqPatch)
		if err != nil {
			fmt.Printf("[ERRO] Rollback falhou ao atualizar produto %d no estoque: %v\n", update.ProdutoID, err)
			continue
		}
		
		if resp != nil {
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("[ERRO] Rollback recebeu status %d ao atualizar produto %d\n", resp.StatusCode, update.ProdutoID)
			} else {
				fmt.Printf("[INFO] Rollback bem-sucedido para produto %d\n", update.ProdutoID)
			}
			resp.Body.Close()
		}
	}
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
