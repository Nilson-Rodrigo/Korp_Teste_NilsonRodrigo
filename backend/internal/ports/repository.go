package ports

import "backend/internal/domain"

// ProdutoRepository define o contrato para persistência de produtos.
type ProdutoRepository interface {
	Listar() ([]domain.Produto, error)
	BuscarPorID(id string) (*domain.Produto, error)
	Criar(produto *domain.Produto) error
	AtualizarSaldo(id string, novoSaldo float64) error
}

// NotaFiscalRepository define o contrato para persistência de notas fiscais.
type NotaFiscalRepository interface {
	Listar() ([]domain.NotaFiscal, error)
	BuscarPorID(id string) (*domain.NotaFiscal, error)
	Criar(nota *domain.NotaFiscal) error
	AtualizarStatus(id string, status string) error
}
