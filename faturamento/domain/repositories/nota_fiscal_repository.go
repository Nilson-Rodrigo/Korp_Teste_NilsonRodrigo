package repositories

import "faturamento/domain/entities"

// NotaFiscalRepository define as operações de persistência para notas fiscais
type NotaFiscalRepository interface {
	// Criar cria uma nova nota fiscal
	Criar(nota *entities.NotaFiscal) error

	// BuscarPorID busca uma nota por ID
	BuscarPorID(id uint) (*entities.NotaFiscal, error)

	// Listar lista todas as notas
	Listar() ([]entities.NotaFiscal, error)

	// AtualizarStatus atualiza o status de uma nota
	AtualizarStatus(id uint, status string) error

	// GerarProximoNumero gera o próximo número de nota sequencial
	GerarProximoNumero() (uint, error)
}

// EstoqueService define a interface para comunicação com o serviço de estoque
type EstoqueService interface {
	// BuscarProduto busca um produto no estoque
	BuscarProduto(id uint) (interface{}, error)

	// AtualizarSaldo atualiza o saldo de um produto no estoque
	AtualizarSaldo(id uint, novoSaldo float64) error

	// BuscarSaldo busca o saldo de um produto
	BuscarSaldo(id uint) (float64, error)
}
