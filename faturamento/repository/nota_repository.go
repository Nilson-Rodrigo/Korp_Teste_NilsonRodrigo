package repository

import (
	"faturamento/config"
	"faturamento/model"
)

func ListarNotas() ([]model.NotaFiscal, error) {
	var notas []model.NotaFiscal
	result := config.DB.Preload("Itens").Find(&notas)
	return notas, result.Error
}

func CriarNota(n *model.NotaFiscal) error {
	var ultima model.NotaFiscal
	config.DB.Order("numero desc").First(&ultima)
	n.Numero = ultima.Numero + 1
	n.Status = "Aberta"
	return config.DB.Create(n).Error
}

func BuscarNotaPorID(id uint) (*model.NotaFiscal, error) {
	var n model.NotaFiscal
	result := config.DB.Preload("Itens").First(&n, id)
	return &n, result.Error
}

func FecharNota(id uint) error {
	return config.DB.Model(&model.NotaFiscal{}).Where("id = ?", id).Update("status", "Fechada").Error
}