package usecases

import (
	"faturamento/domain/entities"
	"faturamento/domain/repositories"
)

// BuscarNotaPorIDUseCase define o caso de uso para buscar uma nota por ID
type BuscarNotaPorIDUseCase struct {
	repository repositories.NotaFiscalRepository
}

// NewBuscarNotaPorIDUseCase cria uma nova instância
func NewBuscarNotaPorIDUseCase(repo repositories.NotaFiscalRepository) *BuscarNotaPorIDUseCase {
	return &BuscarNotaPorIDUseCase{
		repository: repo,
	}
}

// Execute executa o caso de uso
func (uc *BuscarNotaPorIDUseCase) Execute(id uint) (*entities.NotaFiscal, error) {
	nota, err := uc.repository.BuscarPorID(id)
	if err != nil {
		return nil, entities.ErrNotaNaoEncontrada
	}
	return nota, nil
}
