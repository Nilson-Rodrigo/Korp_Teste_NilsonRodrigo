package repository

import (
	"estoque/config"
	"estoque/model"
)

func ListarProdutos() ([]model.Produto, error) {
	var produtos []model.Produto
	result := config.DB.Find(&produtos)
	return produtos, result.Error
}

func CriarProduto(p *model.Produto) error {
	return config.DB.Create(p).Error
}

func BuscarProdutoPorID(id uint) (*model.Produto, error) {
	var p model.Produto
	result := config.DB.First(&p, id)
	return &p, result.Error
}

func AtualizarSaldo(id uint, novoSaldo float64) error {
	return config.DB.Model(&model.Produto{}).Where("id = ?", id).Update("saldo", novoSaldo).Error
}