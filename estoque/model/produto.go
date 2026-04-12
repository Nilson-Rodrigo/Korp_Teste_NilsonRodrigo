package model

import "gorm.io/gorm"

type Produto struct {
	gorm.Model
	Codigo    string  `json:"codigo" gorm:"uniqueIndex;not null"`
	Descricao string  `json:"descricao" gorm:"not null"`
	Saldo     float64 `json:"saldo" gorm:"not null;default:0"`
}