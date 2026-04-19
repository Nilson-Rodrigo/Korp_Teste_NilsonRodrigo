package usecases

import (
	"estoque/domain/entities"
	"estoque/domain/repositories"
)

// CriarProdutoUseCase define o caso de uso para criar um novo produto
type CriarProdutoUseCase struct {
	repository repositories.ProdutoRepository
}

// NewCriarProdutoUseCase cria uma nova instância
func NewCriarProdutoUseCase(repo repositories.ProdutoRepository) *CriarProdutoUseCase {
	return &CriarProdutoUseCase{
		repository: repo,
	}
}

// Execute executa o caso de uso
func (uc *CriarProdutoUseCase) Execute(produto *entities.Produto) error {
	// Validar entidade de domínio
	if err := produto.Validar(); err != nil {
		return err
	}

	// Verificar se código já existe
	_, err := uc.repository.BuscarPorCodigo(produto.Codigo)
	if err == nil {
		return entities.ErrCodigoDuplicado
	}

	// Criar produto
	return uc.repository.Criar(produto)
}
