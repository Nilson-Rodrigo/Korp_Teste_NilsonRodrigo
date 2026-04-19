package repositories

import "estoque/domain/entities"

// ProdutoRepository define as operações de persistência para produtos
type ProdutoRepository interface {
	// Criar cria um novo produto
	Criar(produto *entities.Produto) error

	// BuscarPorID busca um produto por ID
	BuscarPorID(id uint) (*entities.Produto, error)

	// BuscarPorCodigo busca um produto por código
	BuscarPorCodigo(codigo string) (*entities.Produto, error)

	// Listar lista todos os produtos
	Listar() ([]entities.Produto, error)

	// AtualizarSaldo atualiza o saldo de um produto
	AtualizarSaldo(id uint, novoSaldo float64) error

	// Deletar deleta um produto (opcional)
	Deletar(id uint) error
}
