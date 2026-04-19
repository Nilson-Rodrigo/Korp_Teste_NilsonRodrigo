package usecases

import (
	"estoque/domain/entities"
	"estoque/domain/repositories"
)

// ListarProdutosUseCase define o caso de uso para listar todos os produtos
type ListarProdutosUseCase struct {
	repository repositories.ProdutoRepository
}

// NewListarProdutosUseCase cria uma nova instância
func NewListarProdutosUseCase(repo repositories.ProdutoRepository) *ListarProdutosUseCase {
	return &ListarProdutosUseCase{
		repository: repo,
	}
}

// Execute executa o caso de uso
func (uc *ListarProdutosUseCase) Execute() ([]entities.Produto, error) {
	return uc.repository.Listar()
}
