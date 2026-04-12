package model

import "gorm.io/gorm"

type ItemNota struct {
	gorm.Model
	NotaFiscalID uint    `json:"nota_fiscal_id"`
	ProdutoID    uint    `json:"produto_id"`
	Quantidade   float64 `json:"quantidade"`
}

type NotaFiscal struct {
	gorm.Model
	Numero uint       `json:"numero" gorm:"uniqueIndex;not null"`
	Status string     `json:"status" gorm:"default:'Aberta'"`
	Itens  []ItemNota `json:"itens" gorm:"foreignKey:NotaFiscalID"`
}