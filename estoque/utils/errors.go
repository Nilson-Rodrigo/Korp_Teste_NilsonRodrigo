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
	ErrInvalidInput        = NewAPIError("invalid_input", "Entrada inválida", http.StatusBadRequest)
	ErrNotFound            = NewAPIError("not_found", "Recurso não encontrado", http.StatusNotFound)
	ErrInternalError       = NewAPIError("internal_error", "Erro interno do servidor", http.StatusInternalServerError)
	ErrConflict            = NewAPIError("conflict", "Recurso já existe", http.StatusConflict)
	ErrUnprocessableEntity = NewAPIError("unprocessable_entity", "Dados inválidos", http.StatusUnprocessableEntity)
)

// MapDomainErrorToHTTP mapeia erros de domínio para HTTP
func MapDomainErrorToHTTP(err error) APIError {
	if err == nil {
		return APIError{}
	}

	switch err.Error() {
	case "código do produto é obrigatório":
		return NewAPIError("missing_codigo", "Código é obrigatório", http.StatusBadRequest)
	case "descrição do produto é obrigatória":
		return NewAPIError("missing_descricao", "Descrição é obrigatória", http.StatusBadRequest)
	case "saldo não pode ser negativo":
		return NewAPIError("invalid_saldo", "Saldo não pode ser negativo", http.StatusBadRequest)
	case "saldo insuficiente para esta operação":
		return NewAPIError("insufficient_balance", "Saldo insuficiente", http.StatusBadRequest)
	case "produto não encontrado":
		return ErrNotFound
	case "código de produto já existe":
		return NewAPIError("duplicate_codigo", "Código de produto já existe", http.StatusConflict)
	default:
		return ErrInternalError
	}
}
