package application

import "time"

// CriarProdutoInput representa os dados necessários para criar um produto.
type CriarProdutoInput struct {
	Codigo    string  `json:"codigo"`
	Descricao string  `json:"descricao"`
	Saldo     float64 `json:"saldo"`
}

// ProdutoDTO representa os dados de um produto retornados pela API.
type ProdutoDTO struct {
	ID           string    `json:"id"`
	Codigo       string    `json:"codigo"`
	Descricao    string    `json:"descricao"`
	Saldo        float64   `json:"saldo"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// CriarNotaInput representa os dados necessários para criar uma nota fiscal.
type CriarNotaInput struct {
	Itens []ItemNotaInput `json:"itens"`
}

// ItemNotaInput representa os dados de um item ao criar uma nota fiscal.
type ItemNotaInput struct {
	ProdutoID  string  `json:"produto_id"`
	Quantidade float64 `json:"quantidade"`
}

// NotaFiscalDTO representa os dados de uma nota fiscal retornados pela API.
type NotaFiscalDTO struct {
	ID           string        `json:"id"`
	Numero       int           `json:"numero"`
	Status       string        `json:"status"`
	Itens        []ItemNotaDTO `json:"itens"`
	CriadoEm     time.Time     `json:"criado_em"`
	AtualizadoEm time.Time     `json:"atualizado_em"`
}

// ItemNotaDTO representa os dados de um item de nota fiscal retornados pela API.
type ItemNotaDTO struct {
	ProdutoID  string  `json:"produto_id"`
	Quantidade float64 `json:"quantidade"`
}
