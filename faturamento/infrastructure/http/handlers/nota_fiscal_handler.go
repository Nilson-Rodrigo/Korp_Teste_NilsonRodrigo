package handlers

import (
	"context"
	"faturamento/domain/entities"
	"faturamento/domain/usecases"
	"faturamento/infrastructure/http/dto"
	"faturamento/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// NotaFiscalHandler contém os handlers de notas fiscais
type NotaFiscalHandler struct {
	criarNotaUseCase    *usecases.CriarNotaUseCase
	listarNotasUseCase  *usecases.ListarNotasUseCase
	buscarNotaUseCase   *usecases.BuscarNotaPorIDUseCase
	imprimirNotaUseCase *usecases.ImprimirNotaUseCase
}

// NewNotaFiscalHandler cria um novo handler
func NewNotaFiscalHandler(
	criar *usecases.CriarNotaUseCase,
	listar *usecases.ListarNotasUseCase,
	buscar *usecases.BuscarNotaPorIDUseCase,
	imprimir *usecases.ImprimirNotaUseCase,
) *NotaFiscalHandler {
	return &NotaFiscalHandler{
		criarNotaUseCase:    criar,
		listarNotasUseCase:  listar,
		buscarNotaUseCase:   buscar,
		imprimirNotaUseCase: imprimir,
	}
}

// Criar cria uma nova nota fiscal
func (h *NotaFiscalHandler) Criar(c *gin.Context) {
	var req dto.CriarNotaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Erro ao parsear request de criar nota")
		apiErr := utils.MapDomainErrorToHTTP(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	nota := &entities.NotaFiscal{
		Itens: make([]entities.ItemNota, len(req.Itens)),
	}

	for i, item := range req.Itens {
		nota.Itens[i] = entities.ItemNota{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
		}
	}

	if err := h.criarNotaUseCase.Execute(context.Background(), nota); err != nil {
		log.Error().Err(err).Msg("Erro ao criar nota")
		apiErr := utils.MapDomainErrorToHTTP(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	log.Info().Uint("id", nota.ID).Uint("numero", nota.Numero).Msg("Nota fiscal criada")
	c.JSON(http.StatusCreated, dto.NotaFiscalResponse{
		ID:     nota.ID,
		Numero: nota.Numero,
		Status: nota.Status,
	})
}

// Listar lista todas as notas
func (h *NotaFiscalHandler) Listar(c *gin.Context) {
	notas, err := h.listarNotasUseCase.Execute()
	if err != nil {
		log.Error().Err(err).Msg("Erro ao listar notas")
		apiErr := utils.NewAPIError("database_error", "Erro ao listar notas", http.StatusInternalServerError)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	response := make([]dto.NotaFiscalResponse, len(notas))
	for i, nota := range notas {
		itens := make([]dto.ItemNotaResponse, len(nota.Itens))
		for j, item := range nota.Itens {
			itens[j] = dto.ItemNotaResponse{
				ID:         item.ID,
				ProdutoID:  item.ProdutoID,
				Quantidade: item.Quantidade,
			}
		}

		response[i] = dto.NotaFiscalResponse{
			ID:     nota.ID,
			Numero: nota.Numero,
			Status: nota.Status,
			Itens:  itens,
		}
	}

	log.Debug().Int("total", len(notas)).Msg("Notas listadas")
	c.JSON(http.StatusOK, response)
}

// BuscarPorID busca uma nota por ID
func (h *NotaFiscalHandler) BuscarPorID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Warn().Err(err).Msg("ID inválido")
		apiErr := utils.NewAPIError("invalid_id", "ID inválido", http.StatusBadRequest)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	nota, err := h.buscarNotaUseCase.Execute(uint(id))
	if err != nil {
		log.Warn().Uint("id", uint(id)).Err(err).Msg("Nota não encontrada")
		apiErr := utils.MapDomainErrorToHTTP(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	itens := make([]dto.ItemNotaResponse, len(nota.Itens))
	for j, item := range nota.Itens {
		itens[j] = dto.ItemNotaResponse{
			ID:         item.ID,
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
		}
	}

	log.Debug().Uint("id", uint(id)).Msg("Nota encontrada")
	c.JSON(http.StatusOK, dto.NotaFiscalResponse{
		ID:     nota.ID,
		Numero: nota.Numero,
		Status: nota.Status,
		Itens:  itens,
	})
}

// Imprimir imprime uma nota e atualiza saldos
func (h *NotaFiscalHandler) Imprimir(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Warn().Err(err).Msg("ID inválido")
		apiErr := utils.NewAPIError("invalid_id", "ID inválido", http.StatusBadRequest)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	if err := h.imprimirNotaUseCase.Execute(context.Background(), uint(id)); err != nil {
		log.Error().Err(err).Uint("id", uint(id)).Msg("Erro ao imprimir nota")
		apiErr := utils.MapDomainErrorToHTTP(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	log.Info().Uint("id", uint(id)).Msg("Nota impressa com sucesso")
	c.JSON(http.StatusOK, gin.H{"message": "Nota impressa com sucesso"})
}
