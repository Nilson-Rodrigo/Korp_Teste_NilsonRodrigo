package usecases

import (
	"estoque/domain/entities"
	"estoque/domain/repositories"
)

// BuscarProdutoPorIDUseCase define o caso de uso para buscar um produto por ID
type BuscarProdutoPorIDUseCase struct {
	repository repositories.ProdutoRepository
}

// NewBuscarProdutoPorIDUseCase cria uma nova instância
func NewBuscarProdutoPorIDUseCase(repo repositories.ProdutoRepository) *BuscarProdutoPorIDUseCase {
	return &BuscarProdutoPorIDUseCase{
		repository: repo,
	}
}

// Execute executa o caso de uso
func (uc *BuscarProdutoPorIDUseCase) Execute(id uint) (*entities.Produto, error) {
	produto, err := uc.repository.BuscarPorID(id)
	if err != nil {
		return nil, entities.ErrProdutoNaoEncontrado
	}
	return produto, nil
}
