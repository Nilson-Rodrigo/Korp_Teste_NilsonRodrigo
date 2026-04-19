package usecases

import (
	"estoque/domain/entities"
	"estoque/domain/repositories"
)

// AtualizarSaldoUseCase define o caso de uso para atualizar o saldo de um produto
type AtualizarSaldoUseCase struct {
	repository repositories.ProdutoRepository
}

// NewAtualizarSaldoUseCase cria uma nova instância
func NewAtualizarSaldoUseCase(repo repositories.ProdutoRepository) *AtualizarSaldoUseCase {
	return &AtualizarSaldoUseCase{
		repository: repo,
	}
}

// Execute executa o caso de uso
func (uc *AtualizarSaldoUseCase) Execute(id uint, novoSaldo float64) error {
	produto, err := uc.repository.BuscarPorID(id)
	if err != nil {
		return entities.ErrProdutoNaoEncontrado
	}

	if novoSaldo < 0 {
		return entities.ErrSaldoNegativo
	}

	return uc.repository.AtualizarSaldo(produto.ID, novoSaldo)
}
