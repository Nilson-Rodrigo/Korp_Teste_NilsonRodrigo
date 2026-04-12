package ports

import "backend/internal/application"

// ProdutoService define o contrato para operações de negócio com produtos.
type ProdutoService interface {
	ListarProdutos() ([]application.ProdutoDTO, error)
	CriarProduto(input application.CriarProdutoInput) (*application.ProdutoDTO, error)
	BuscarProdutoPorID(id string) (*application.ProdutoDTO, error)
	AtualizarSaldo(id string, novoSaldo float64) error
}

// NotaFiscalService define o contrato para operações de negócio com notas fiscais.
type NotaFiscalService interface {
	ListarNotas() ([]application.NotaFiscalDTO, error)
	CriarNota(input application.CriarNotaInput) (*application.NotaFiscalDTO, error)
	BuscarNotaPorID(id string) (*application.NotaFiscalDTO, error)
	ImprimirNota(id string) error
}
