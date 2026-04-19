package usecases

import (
	"faturamento/domain/entities"
	"faturamento/domain/repositories"
)

// ListarNotasUseCase define o caso de uso para listar notas fiscais
type ListarNotasUseCase struct {
	repository repositories.NotaFiscalRepository
}

// NewListarNotasUseCase cria uma nova instância
func NewListarNotasUseCase(repo repositories.NotaFiscalRepository) *ListarNotasUseCase {
	return &ListarNotasUseCase{
		repository: repo,
	}
}

// Execute executa o caso de uso
func (uc *ListarNotasUseCase) Execute() ([]entities.NotaFiscal, error) {
	return uc.repository.Listar()
}
