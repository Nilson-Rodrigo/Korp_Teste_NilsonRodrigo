package http

import (
	"backend/internal/application"
	"backend/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
)

// Handler contém os handlers HTTP para a API.
type Handler struct {
	produtoService *application.ProdutoServiceImpl
	notaService    *application.NotaFiscalServiceImpl
}

// NovoHandler cria uma nova instância do handler HTTP.
func NovoHandler(produtoService *application.ProdutoServiceImpl, notaService *application.NotaFiscalServiceImpl) *Handler {
	return &Handler{
		produtoService: produtoService,
		notaService:    notaService,
	}
}

// writeJSON escreve uma resposta JSON com o status code especificado.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// writeError escreve uma resposta de erro em JSON.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"erro": message})
}

// mapErroParaStatus converte erros de domínio em status HTTP.
func mapErroParaStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrInvoiceAlreadyClosed):
		return http.StatusConflict
	case errors.Is(err, domain.ErrInsufficientStock):
		return http.StatusUnprocessableEntity
	case errors.Is(err, domain.ErrDuplicateCode):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Health retorna o status de saúde da API.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ListarProdutos retorna todos os produtos.
func (h *Handler) ListarProdutos(w http.ResponseWriter, r *http.Request) {
	produtos, err := h.produtoService.ListarProdutos()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, produtos)
}

// CriarProduto cria um novo produto.
func (h *Handler) CriarProduto(w http.ResponseWriter, r *http.Request) {
	var input application.CriarProdutoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	produto, err := h.produtoService.CriarProduto(input)
	if err != nil {
		writeError(w, mapErroParaStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, produto)
}

// BuscarProdutoPorID busca um produto pelo ID.
func (h *Handler) BuscarProdutoPorID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id é obrigatório")
		return
	}

	produto, err := h.produtoService.BuscarProdutoPorID(id)
	if err != nil {
		writeError(w, mapErroParaStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, produto)
}

// AtualizarSaldo atualiza o saldo de um produto.
func (h *Handler) AtualizarSaldo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id é obrigatório")
		return
	}

	var body struct {
		Saldo float64 `json:"saldo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if err := h.produtoService.AtualizarSaldo(id, body.Saldo); err != nil {
		writeError(w, mapErroParaStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensagem": "saldo atualizado com sucesso"})
}

// ListarNotas retorna todas as notas fiscais.
func (h *Handler) ListarNotas(w http.ResponseWriter, r *http.Request) {
	notas, err := h.notaService.ListarNotas()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, notas)
}

// CriarNota cria uma nova nota fiscal.
func (h *Handler) CriarNota(w http.ResponseWriter, r *http.Request) {
	var input application.CriarNotaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	nota, err := h.notaService.CriarNota(input)
	if err != nil {
		writeError(w, mapErroParaStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, nota)
}

// BuscarNotaPorID busca uma nota fiscal pelo ID.
func (h *Handler) BuscarNotaPorID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id é obrigatório")
		return
	}

	nota, err := h.notaService.BuscarNotaPorID(id)
	if err != nil {
		writeError(w, mapErroParaStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, nota)
}

// ImprimirNota fecha a nota fiscal e deduz o estoque.
func (h *Handler) ImprimirNota(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id é obrigatório")
		return
	}

	if err := h.notaService.ImprimirNota(id); err != nil {
		writeError(w, mapErroParaStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensagem": "nota fiscal impressa e fechada com sucesso"})
}
