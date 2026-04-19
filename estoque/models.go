package main

import "time"

// Produto é a entidade que representa um produto em estoque
type Produto struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Codigo    string    `json:"codigo" gorm:"uniqueIndex"`
	Descricao string    `json:"descricao"`
	Saldo     float64   `json:"saldo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName especifica o nome da tabela
func (Produto) TableName() string {
	return "produtos"
}
