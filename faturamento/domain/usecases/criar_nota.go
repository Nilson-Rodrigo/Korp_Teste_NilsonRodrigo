package usecases

import (
	"context"
	"faturamento/domain/entities"
	"faturamento/domain/repositories"
)

// CriarNotaUseCase define o caso de uso para criar uma nota fiscal
type CriarNotaUseCase struct {
	repository repositories.NotaFiscalRepository
	estoque    repositories.EstoqueService
}

// NewCriarNotaUseCase cria uma nova instância
func NewCriarNotaUseCase(repo repositories.NotaFiscalRepository, estoque repositories.EstoqueService) *CriarNotaUseCase {
	return &CriarNotaUseCase{
		repository: repo,
		estoque:    estoque,
	}
}

// Execute executa o caso de uso
func (uc *CriarNotaUseCase) Execute(ctx context.Context, nota *entities.NotaFiscal) error {
	// Validar entidade
	if err := nota.CriarNota(); err != nil {
		return err
	}

	// Validar produtos no estoque
	for _, item := range nota.Itens {
		produto, err := uc.estoque.BuscarProduto(item.ProdutoID)
		if err != nil {
			return entities.ErrProdutoIndisponivel
		}
		if produto == nil {
			return entities.ErrProdutoIndisponivel
		}
	}

	// Gerar número sequencial
	numero, err := uc.repository.GerarProximoNumero()
	if err != nil {
		return err
	}
	nota.Numero = numero

	// Criar nota
	return uc.repository.Criar(nota)
}
