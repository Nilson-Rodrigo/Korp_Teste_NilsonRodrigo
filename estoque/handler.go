package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Handler contém todos os handlers de produtos
type Handler struct {
	db *gorm.DB
}

// NewHandler cria um novo handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// Criar cria um novo produto
func (h *Handler) Criar(c *gin.Context) {
	var req struct {
		Codigo    string  `json:"codigo" binding:"required,min=1,max=50"`
		Descricao string  `json:"descricao" binding:"required,min=3,max=500"`
		Saldo     float64 `json:"saldo" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Erro ao parsear request")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "invalid_input",
			"message": err.Error(),
		})
		return
	}

	// Validar se código já existe
	var count int64
	if err := h.db.Model(&Produto{}).Where("codigo = ?", req.Codigo).Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Erro ao verificar código duplicado")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "database_error",
			"message": "Erro ao criar produto",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"code":    "duplicate_codigo",
			"message": "Código de produto já existe",
		})
		return
	}

	produto := Produto{
		Codigo:    req.Codigo,
		Descricao: req.Descricao,
		Saldo:     req.Saldo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.db.Create(&produto).Error; err != nil {
		log.Error().Err(err).Str("codigo", req.Codigo).Msg("Erro ao criar produto")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "database_error",
			"message": "Erro ao criar produto",
		})
		return
	}

	log.Info().Uint("id", produto.ID).Str("codigo", produto.Codigo).Msg("Produto criado com sucesso")
	c.JSON(http.StatusCreated, gin.H{
		"id":        produto.ID,
		"codigo":    produto.Codigo,
		"descricao": produto.Descricao,
		"saldo":     produto.Saldo,
	})
}

// Listar lista todos os produtos
func (h *Handler) Listar(c *gin.Context) {
	var produtos []Produto

	if err := h.db.Find(&produtos).Error; err != nil {
		log.Error().Err(err).Msg("Erro ao listar produtos")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "database_error",
			"message": "Erro ao listar produtos",
		})
		return
	}

	log.Debug().Int("total", len(produtos)).Msg("Produtos listados com sucesso")
	c.JSON(http.StatusOK, produtos)
}

// BuscarPorID busca um produto por ID
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

	var produto Produto
	if err := h.db.First(&produto, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Uint("id", uint(id)).Msg("Produto não encontrado")
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "not_found",
				"message": "Produto não encontrado",
			})
			return
		}
		log.Error().Err(err).Msg("Erro ao buscar produto")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "database_error",
			"message": "Erro ao buscar produto",
		})
		return
	}

	log.Debug().Uint("id", uint(id)).Msg("Produto encontrado")
	c.JSON(http.StatusOK, produto)
}

// AtualizarSaldo atualiza o saldo de um produto
func (h *Handler) AtualizarSaldo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("ID inválido")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "invalid_id",
			"message": "ID inválido",
		})
		return
	}

	var req struct {
		Saldo float64 `json:"saldo" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Erro ao parsear request")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "invalid_input",
			"message": err.Error(),
		})
		return
	}

	var produto Produto
	if err := h.db.First(&produto, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "not_found",
				"message": "Produto não encontrado",
			})
			return
		}
		log.Error().Err(err).Msg("Erro ao buscar produto")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "database_error",
			"message": "Erro ao atualizar saldo",
		})
		return
	}

	if err := h.db.Model(&produto).Update("saldo", req.Saldo).Update("updated_at", time.Now()).Error; err != nil {
		log.Error().Err(err).Uint("id", uint(id)).Float64("saldo", req.Saldo).Msg("Erro ao atualizar saldo")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "database_error",
			"message": "Erro ao atualizar saldo",
		})
		return
	}

	log.Info().Uint("id", uint(id)).Float64("novo_saldo", req.Saldo).Msg("Saldo atualizado com sucesso")
	c.JSON(http.StatusOK, gin.H{"message": "Saldo atualizado com sucesso"})
}
