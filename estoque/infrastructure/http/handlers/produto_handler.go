package handlers

import (
	"estoque/domain/entities"
	"estoque/domain/usecases"
	"estoque/infrastructure/http/dto"
	"estoque/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ProdutoHandler contém os handlers de produtos
type ProdutoHandler struct {
	criarProdutoUseCase   *usecases.CriarProdutoUseCase
	listarProdutosUseCase *usecases.ListarProdutosUseCase
	buscarProdutoUseCase  *usecases.BuscarProdutoPorIDUseCase
	atualizarSaldoUseCase *usecases.AtualizarSaldoUseCase
}

// NewProdutoHandler cria um novo handler de produtos
func NewProdutoHandler(
	criar *usecases.CriarProdutoUseCase,
	listar *usecases.ListarProdutosUseCase,
	buscar *usecases.BuscarProdutoPorIDUseCase,
	atualizar *usecases.AtualizarSaldoUseCase,
) *ProdutoHandler {
	return &ProdutoHandler{
		criarProdutoUseCase:   criar,
		listarProdutosUseCase: listar,
		buscarProdutoUseCase:  buscar,
		atualizarSaldoUseCase: atualizar,
	}
}

// Criar cria um novo produto
func (h *ProdutoHandler) Criar(c *gin.Context) {
	var req dto.CriarProdutoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Erro ao parsear request de criar produto")
		apiErr := utils.MapDomainErrorToHTTP(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	produto := &entities.Produto{
		Codigo:    req.Codigo,
		Descricao: req.Descricao,
		Saldo:     req.Saldo,
	}

	if err := h.criarProdutoUseCase.Execute(produto); err != nil {
		log.Error().Err(err).Str("codigo", req.Codigo).Msg("Erro ao criar produto")
		apiErr := utils.MapDomainErrorToHTTP(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	log.Info().Uint("id", produto.ID).Str("codigo", produto.Codigo).Msg("Produto criado com sucesso")
	c.JSON(http.StatusCreated, dto.ProdutoResponse{
		ID:        produto.ID,
		Codigo:    produto.Codigo,
		Descricao: produto.Descricao,
		Saldo:     produto.Saldo,
	})
}

// Listar lista todos os produtos
func (h *ProdutoHandler) Listar(c *gin.Context) {
	produtos, err := h.listarProdutosUseCase.Execute()
	if err != nil {
		log.Error().Err(err).Msg("Erro ao listar produtos")
		apiErr := utils.NewAPIError("database_error", "Erro ao listar produtos", http.StatusInternalServerError)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	response := make([]dto.ProdutoResponse, len(produtos))
	for i, p := range produtos {
		response[i] = dto.ProdutoResponse{
			ID:        p.ID,
			Codigo:    p.Codigo,
			Descricao: p.Descricao,
			Saldo:     p.Saldo,
		}
	}

	log.Debug().Int("total", len(produtos)).Msg("Produtos listados com sucesso")
	c.JSON(http.StatusOK, response)
}

// BuscarPorID busca um produto por ID
func (h *ProdutoHandler) BuscarPorID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("ID inválido")
		apiErr := utils.NewAPIError("invalid_id", "ID inválido", http.StatusBadRequest)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	produto, err := h.buscarProdutoUseCase.Execute(uint(id))
	if err != nil {
		log.Warn().Uint("id", uint(id)).Err(err).Msg("Produto não encontrado")
		apiErr := utils.MapDomainErrorToHTTP(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	log.Debug().Uint("id", uint(id)).Msg("Produto encontrado")
	c.JSON(http.StatusOK, dto.ProdutoResponse{
		ID:        produto.ID,
		Codigo:    produto.Codigo,
		Descricao: produto.Descricao,
		Saldo:     produto.Saldo,
	})
}

// AtualizarSaldo atualiza o saldo de um produto
func (h *ProdutoHandler) AtualizarSaldo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("ID inválido")
		apiErr := utils.NewAPIError("invalid_id", "ID inválido", http.StatusBadRequest)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	var req dto.AtualizarSaldoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Erro ao parsear request de atualizar saldo")
		apiErr := utils.MapDomainErrorToHTTP(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	if err := h.atualizarSaldoUseCase.Execute(uint(id), req.Saldo); err != nil {
		log.Error().Err(err).Uint("id", uint(id)).Float64("saldo", req.Saldo).Msg("Erro ao atualizar saldo")
		apiErr := utils.MapDomainErrorToHTTP(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	log.Info().Uint("id", uint(id)).Float64("novo_saldo", req.Saldo).Msg("Saldo atualizado com sucesso")
	c.JSON(http.StatusOK, gin.H{"message": "Saldo atualizado com sucesso"})
}
