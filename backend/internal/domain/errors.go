package domain

import "errors"

// Erros de domínio utilizados em toda a aplicação.
var (
	ErrNotFound             = errors.New("recurso não encontrado")
	ErrInvalidInput         = errors.New("dados inválidos")
	ErrInvoiceAlreadyClosed = errors.New("nota fiscal já está fechada")
	ErrInsufficientStock    = errors.New("estoque insuficiente")
	ErrDuplicateCode        = errors.New("código já existe")
)
