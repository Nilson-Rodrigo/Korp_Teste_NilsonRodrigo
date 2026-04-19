package utils

import (
	"net/http"
	"time"
)

// APIError é a estrutura padrão de erro da API
type APIError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Status    int    `json:"-"`
}

func (e APIError) Error() string {
	return e.Message
}

// NewAPIError cria um novo erro de API
func NewAPIError(code, message string, status int) APIError {
	return APIError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Status:    status,
	}
}

// Erros comuns
var (
	ErrInvalidInput  = NewAPIError("invalid_input", "Entrada inválida", http.StatusBadRequest)
	ErrNotFound      = NewAPIError("not_found", "Nota não encontrada", http.StatusNotFound)
	ErrInternalError = NewAPIError("internal_error", "Erro interno do servidor", http.StatusInternalServerError)
	ErrBadRequest    = NewAPIError("bad_request", "Solicitação inválida", http.StatusBadRequest)
)

// MapDomainErrorToHTTP mapeia erros de domínio para HTTP
func MapDomainErrorToHTTP(err error) APIError {
	if err == nil {
		return APIError{}
	}

	switch err.Error() {
	case "nota fiscal deve conter ao menos um item":
		return NewAPIError("empty_items", "Nota deve conter itens", http.StatusBadRequest)
	case "quantidade deve ser maior que zero":
		return NewAPIError("invalid_quantity", "Quantidade deve ser maior que zero", http.StatusBadRequest)
	case "apenas notas com status Aberta podem ser impressas":
		return NewAPIError("invalid_status", "Apenas notas abertas podem ser impressas", http.StatusBadRequest)
	case "nota fiscal não encontrada":
		return ErrNotFound
	case "produto indisponível no estoque":
		return NewAPIError("product_unavailable", "Produto indisponível no estoque", http.StatusBadRequest)
	case "saldo insuficiente no estoque":
		return NewAPIError("insufficient_balance", "Saldo insuficiente no estoque", http.StatusBadRequest)
	case "serviço de estoque indisponível":
		return NewAPIError("estoque_service_unavailable", "Serviço de estoque indisponível", http.StatusServiceUnavailable)
	default:
		return ErrInternalError
	}
}
