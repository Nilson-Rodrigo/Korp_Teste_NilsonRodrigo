package domain

import "time"

// NotaFiscal representa uma nota fiscal no sistema de faturamento.
type NotaFiscal struct {
	ID           string     `json:"id"`
	Numero       int        `json:"numero"`
	Status       string     `json:"status"` // "Aberta" ou "Fechada"
	Itens        []ItemNota `json:"itens"`
	CriadoEm     time.Time  `json:"criado_em"`
	AtualizadoEm time.Time  `json:"atualizado_em"`
}

// ItemNota representa um item dentro de uma nota fiscal.
type ItemNota struct {
	ProdutoID  string  `json:"produto_id"`
	Quantidade float64 `json:"quantidade"`
}
