package entities

import "errors"

var (
	ErrCodigoObrigatorio    = errors.New("código do produto é obrigatório")
	ErrDescricaoObrigatoria = errors.New("descrição do produto é obrigatória")
	ErrSaldoNegativo        = errors.New("saldo não pode ser negativo")
	ErrSaldoInsuficiente    = errors.New("saldo insuficiente para esta operação")
	ErrProdutoNaoEncontrado = errors.New("produto não encontrado")
	ErrCodigoDuplicado      = errors.New("código de produto já existe")
)
