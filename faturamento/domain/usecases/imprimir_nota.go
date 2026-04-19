package usecases

import (
	"context"
	"faturamento/domain/entities"
	"faturamento/domain/repositories"

	"gorm.io/gorm"
)

// ImprimirNotaUseCase define o caso de uso para imprimir uma nota fiscal
type ImprimirNotaUseCase struct {
	repository repositories.NotaFiscalRepository
	estoque    repositories.EstoqueService
	db         *gorm.DB
}

// NewImprimirNotaUseCase cria uma nova instância
func NewImprimirNotaUseCase(
	repo repositories.NotaFiscalRepository,
	estoque repositories.EstoqueService,
	db *gorm.DB,
) *ImprimirNotaUseCase {
	return &ImprimirNotaUseCase{
		repository: repo,
		estoque:    estoque,
		db:         db,
	}
}

// Execute executa o caso de uso com transação
func (uc *ImprimirNotaUseCase) Execute(ctx context.Context, notaID uint) error {
	// Iniciar transação
	tx := uc.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Buscar nota
	nota, err := uc.repository.BuscarPorID(notaID)
	if err != nil {
		tx.Rollback()
		return entities.ErrNotaNaoEncontrada
	}

	// Validar status
	if err := nota.ImprimirNota(); err != nil {
		tx.Rollback()
		return err
	}

	// Validar saldos no estoque ANTES de qualquer atualização
	for _, item := range nota.Itens {
		saldo, err := uc.estoque.BuscarSaldo(item.ProdutoID)
		if err != nil {
			tx.Rollback()
			return entities.ErrEstoqueIndisponivel
		}

		if saldo < item.Quantidade {
			tx.Rollback()
			return entities.ErrSaldoInsuficiente
		}
	}

	// Atualizar saldos no estoque
	for _, item := range nota.Itens {
		saldo, _ := uc.estoque.BuscarSaldo(item.ProdutoID)
		novoSaldo := saldo - item.Quantidade

		if err := uc.estoque.AtualizarSaldo(item.ProdutoID, novoSaldo); err != nil {
			tx.Rollback()
			return entities.ErrEstoqueIndisponivel
		}
	}

	// Atualizar status da nota dentro da transação
	if err := tx.Model(&entities.NotaFiscal{}).Where("id = ?", notaID).Update("status", "Fechada").Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commitar transação
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
