package main

import "time"

// NotaFiscal é a entidade que representa uma nota fiscal
type NotaFiscal struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Numero    uint       `json:"numero" gorm:"uniqueIndex"`
	Status    string     `json:"status"` // "Aberta" ou "Fechada"
	Itens     []ItemNota `json:"itens" gorm:"foreignKey:NotaFiscalID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ItemNota representa um item dentro de uma nota fiscal
type ItemNota struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	NotaFiscalID uint      `json:"nota_fiscal_id"`
	ProdutoID    uint      `json:"produto_id"`
	Quantidade   float64   `json:"quantidade"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName especifica o nome da tabela para NotaFiscal
func (NotaFiscal) TableName() string {
	return "notas_fiscais"
}

// TableName especifica o nome da tabela para ItemNota
func (ItemNota) TableName() string {
	return "itens_nota"
}
