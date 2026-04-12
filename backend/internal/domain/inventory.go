package domain

import "time"

// Produto representa um produto no sistema de estoque.
type Produto struct {
	ID           string    `json:"id"`
	Codigo       string    `json:"codigo"`
	Descricao    string    `json:"descricao"`
	Saldo        float64   `json:"saldo"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}
