package persistence

import (
	"faturamento/domain/entities"

	"gorm.io/gorm"
)

// NotaFiscalRepositoryImpl implementa a interface NotaFiscalRepository
type NotaFiscalRepositoryImpl struct {
	db *gorm.DB
}

// NewNotaFiscalRepository cria uma nova instância do repositório
func NewNotaFiscalRepository(db *gorm.DB) *NotaFiscalRepositoryImpl {
	return &NotaFiscalRepositoryImpl{
		db: db,
	}
}

// Criar cria uma nova nota fiscal
func (r *NotaFiscalRepositoryImpl) Criar(nota *entities.NotaFiscal) error {
	return r.db.Create(nota).Error
}

// BuscarPorID busca uma nota por ID (com itens)
func (r *NotaFiscalRepositoryImpl) BuscarPorID(id uint) (*entities.NotaFiscal, error) {
	var nota entities.NotaFiscal
	if err := r.db.Preload("Itens").First(&nota, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrNotaNaoEncontrada
		}
		return nil, err
	}
	return &nota, nil
}

// Listar lista todas as notas fiscais
func (r *NotaFiscalRepositoryImpl) Listar() ([]entities.NotaFiscal, error) {
	var notas []entities.NotaFiscal
	if err := r.db.Preload("Itens").Find(&notas).Error; err != nil {
		return nil, err
	}
	return notas, nil
}

// AtualizarStatus atualiza o status de uma nota
func (r *NotaFiscalRepositoryImpl) AtualizarStatus(id uint, status string) error {
	return r.db.Model(&entities.NotaFiscal{}).Where("id = ?", id).Update("status", status).Error
}

// GerarProximoNumero gera o próximo número sequencial
func (r *NotaFiscalRepositoryImpl) GerarProximoNumero() (uint, error) {
	var ultimaNota entities.NotaFiscal

	if err := r.db.Order("numero DESC").First(&ultimaNota).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Primeira nota
			return 1, nil
		}
		return 0, err
	}

	return ultimaNota.Numero + 1, nil
}
