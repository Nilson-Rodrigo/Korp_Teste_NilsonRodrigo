package repository

import (
	"faturamento/config"
	"faturamento/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ListarNotas() ([]model.NotaFiscal, error) {
	var notas []model.NotaFiscal
	result := config.DB.Preload("Itens").Find(&notas)
	return notas, result.Error
}

func CriarNota(n *model.NotaFiscal) error {
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var ultima model.NotaFiscal
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Order("numero desc").First(&ultima).Error; err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}

	n.Numero = ultima.Numero + 1
	n.Status = "Aberta"

	if err := tx.Create(n).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func BuscarNotaPorID(id uint) (*model.NotaFiscal, error) {
	var n model.NotaFiscal
	result := config.DB.Preload("Itens").First(&n, id)
	return &n, result.Error
}

func FecharNota(id uint) error {
	return config.DB.Model(&model.NotaFiscal{}).Where("id = ?", id).Update("status", "Fechada").Error
}
