package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Handler contém todos os handlers de notas fiscais
type Handler struct {
	db         *gorm.DB
	estoqueURL string
}

// NewHandler cria um novo handler
func NewHandler(db *gorm.DB, estoqueURL string) *Handler {
	return &Handler{
		db:         db,
		estoqueURL: estoqueURL,
	}
}

// Criar cria uma nova nota fiscal
func (h *Handler) Criar(c *gin.Context) {
	// Aceitar com tipos flexíveis
	var req struct {
		Itens []map[string]interface{} `json:"itens"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Erro ao parsear request")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "invalid_input",
			"message": err.Error(),
		})
		return
	}

	if len(req.Itens) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "no_items",
			"message": "Pelo menos um item é obrigatório",
		})
		return
	}

	// Converter e validar itens
	type ItemReq struct {
		ProdutoID  uint
		Quantidade float64
	}
	var itens []ItemReq

	for _, itemRaw := range req.Itens {
		item := ItemReq{}

		// Parse produto_id
		if prodIdVal, ok := itemRaw["produto_id"]; ok {
			switch v := prodIdVal.(type) {
			case float64:
				item.ProdutoID = uint(v)
			case string:
				id, err := strconv.ParseUint(v, 10, 32)
				if err != nil {
					log.Warn().Err(err).Str("produto_id", v).Msg("Erro ao converter produto_id")
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    "invalid_produto_id",
						"message": "ID do produto inválido",
					})
					return
				}
				item.ProdutoID = uint(id)
			}
		}

		// Parse quantidade
		if qtyVal, ok := itemRaw["quantidade"]; ok {
			switch v := qtyVal.(type) {
			case float64:
				item.Quantidade = v
			case string:
				qty, err := strconv.ParseFloat(v, 64)
				if err != nil {
					log.Warn().Err(err).Str("quantidade", v).Msg("Erro ao converter quantidade")
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    "invalid_quantidade",
						"message": "Quantidade inválida",
					})
					return
				}
				item.Quantidade = qty
			}
		}

		if item.ProdutoID == 0 || item.Quantidade == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "invalid_item",
				"message": "Produto ID e quantidade são obrigatórios",
			})
			return
		}

		itens = append(itens, item)
	}

	log.Info().Interface("itens", itens).Msg("Request de criação de nota recebida")

	// Validar se produtos existem no estoque
	for _, item := range itens {
		if err := h.verificarProdutoEstoque(item.ProdutoID); err != nil {
			log.Warn().Err(err).Uint("produto_id", item.ProdutoID).Msg("Produto não encontrado")
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "produto_indisponivel",
				"message": "Produto não encontrado no estoque",
			})
			return
		}
	}

	// Gerar próximo número
	var ultimaNota NotaFiscal
	numero := uint(1)
	if err := h.db.Order("numero DESC").First(&ultimaNota).Error; err == nil {
		numero = ultimaNota.Numero + 1
	}

	// Criar nota
	nota := NotaFiscal{
		Numero:    numero,
		Status:    "Aberta",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.db.Create(&nota).Error; err != nil {
		log.Error().Err(err).Msg("Erro ao criar nota")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "database_error",
			"message": "Erro ao criar nota",
		})
		return
	}

	// Criar itens
	notasItens := make([]ItemNota, len(itens))
	for i, item := range itens {
		itemNota := ItemNota{
			NotaFiscalID: nota.ID,
			ProdutoID:    item.ProdutoID,
			Quantidade:   item.Quantidade,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := h.db.Create(&itemNota).Error; err != nil {
			log.Error().Err(err).Msg("Erro ao criar item de nota")
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "database_error",
				"message": "Erro ao criar nota",
			})
			return
		}
		notasItens[i] = itemNota
	}

	log.Info().Uint("id", nota.ID).Uint("numero", nota.Numero).Msg("Nota fiscal criada")
	c.JSON(http.StatusCreated, gin.H{
		"id":     nota.ID,
		"numero": nota.Numero,
		"status": nota.Status,
	})
}

// Listar lista todas as notas
func (h *Handler) Listar(c *gin.Context) {
	var notas []NotaFiscal

	if err := h.db.Preload("Itens").Find(&notas).Error; err != nil {
		log.Error().Err(err).Msg("Erro ao listar notas")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "database_error",
			"message": "Erro ao listar notas",
		})
		return
	}

	log.Debug().Int("total", len(notas)).Msg("Notas listadas")
	c.JSON(http.StatusOK, notas)
}

// BuscarPorID busca uma nota por ID
func (h *Handler) BuscarPorID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("ID inválido")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "invalid_id",
			"message": "ID inválido",
		})
		return
	}

	var nota NotaFiscal
	if err := h.db.Preload("Itens").First(&nota, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "not_found",
				"message": "Nota não encontrada",
			})
			return
		}
		log.Error().Err(err).Msg("Erro ao buscar nota")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "database_error",
			"message": "Erro ao buscar nota",
		})
		return
	}

	log.Debug().Uint("id", uint(id)).Msg("Nota encontrada")
	c.JSON(http.StatusOK, nota)
}

// Imprimir imprime uma nota e atualiza saldos
func (h *Handler) Imprimir(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Warn().Err(err).Msg("ID inválido")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "invalid_id",
			"message": "ID inválido",
		})
		return
	}

	// Iniciar transação
	tx := h.db.Begin()

	// Buscar nota
	var nota NotaFiscal
	if err := tx.Preload("Itens").First(&nota, id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "not_found",
				"message": "Nota não encontrada",
			})
			return
		}
		log.Error().Err(err).Msg("Erro ao buscar nota")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "database_error",
			"message": "Erro ao imprimir nota",
		})
		return
	}

	// Validar se está aberta
	if nota.Status != "Aberta" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "invalid_status",
			"message": "Apenas notas com status Aberta podem ser impressas",
		})
		return
	}

	// Validar saldos no estoque ANTES de qualquer atualização
	for _, item := range nota.Itens {
		saldo, err := h.buscarSaldoEstoque(item.ProdutoID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "estoque_indisponivel",
				"message": "Serviço de estoque indisponível",
			})
			return
		}

		if saldo < item.Quantidade {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "saldo_insuficiente",
				"message": "Saldo insuficiente no estoque",
			})
			return
		}
	}

	// Atualizar saldos no estoque
	for _, item := range nota.Itens {
		saldo, _ := h.buscarSaldoEstoque(item.ProdutoID)
		novoSaldo := saldo - item.Quantidade

		if err := h.atualizarSaldoEstoque(item.ProdutoID, novoSaldo); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "estoque_indisponivel",
				"message": "Erro ao atualizar estoque",
			})
			return
		}
	}

	// Atualizar status da nota dentro da transação
	if err := tx.Model(&NotaFiscal{}).Where("id = ?", id).Update("status", "Fechada").Update("updated_at", time.Now()).Error; err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Erro ao atualizar status da nota")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "database_error",
			"message": "Erro ao imprimir nota",
		})
		return
	}

	// Commitar transação
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Erro ao confirmar transação")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "database_error",
			"message": "Erro ao imprimir nota",
		})
		return
	}

	log.Info().Uint("id", uint(id)).Msg("Nota impressa com sucesso")
	c.JSON(http.StatusOK, gin.H{"message": "Nota impressa com sucesso"})
}

// verificarProdutoEstoque verifica se um produto existe no estoque
func (h *Handler) verificarProdutoEstoque(produtoID uint) error {
	url := fmt.Sprintf("%s/produtos/%d", h.estoqueURL, produtoID)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Error().Err(err).Uint("produto_id", produtoID).Msg("Erro ao conectar ao serviço de estoque")
		return fmt.Errorf("estoque indisponível")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().Int("status", resp.StatusCode).Uint("produto_id", produtoID).Msg("Produto não encontrado")
		return fmt.Errorf("produto não encontrado")
	}

	return nil
}

// buscarSaldoEstoque busca o saldo de um produto no estoque
func (h *Handler) buscarSaldoEstoque(produtoID uint) (float64, error) {
	url := fmt.Sprintf("%s/produtos/%d", h.estoqueURL, produtoID)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("produto não encontrado")
	}

	var produto struct {
		ID    uint    `json:"id"`
		Saldo float64 `json:"saldo"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&produto); err != nil {
		return 0, err
	}

	return produto.Saldo, nil
}

// atualizarSaldoEstoque atualiza o saldo de um produto no estoque
func (h *Handler) atualizarSaldoEstoque(produtoID uint, novoSaldo float64) error {
	url := fmt.Sprintf("%s/produtos/%d/saldo", h.estoqueURL, produtoID)

	payload := map[string]float64{"saldo": novoSaldo}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Erro ao conectar ao serviço de estoque")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().Int("status", resp.StatusCode).Msg("Erro ao atualizar saldo")
		return fmt.Errorf("erro ao atualizar saldo no estoque")
	}

	return nil
}
