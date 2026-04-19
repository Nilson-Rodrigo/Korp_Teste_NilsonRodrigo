package entities

import "errors"

var (
	ErrNotaSemItens        = errors.New("nota fiscal deve conter ao menos um item")
	ErrQuantidadeInvalida  = errors.New("quantidade deve ser maior que zero")
	ErrNotaNaoAberta       = errors.New("apenas notas com status Aberta podem ser impressas")
	ErrNotaNaoEncontrada   = errors.New("nota fiscal não encontrada")
	ErrProdutoIndisponivel = errors.New("produto indisponível no estoque")
	ErrSaldoInsuficiente   = errors.New("saldo insuficiente no estoque")
	ErrEstoqueIndisponivel = errors.New("serviço de estoque indisponível")
)
