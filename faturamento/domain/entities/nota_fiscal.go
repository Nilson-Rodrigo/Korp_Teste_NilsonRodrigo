package entities

import "time"

// ItemNota representa um item dentro de uma nota fiscal
type ItemNota struct {
	ID           uint      `json:"id"`
	NotaFiscalID uint      `json:"nota_fiscal_id"`
	ProdutoID    uint      `json:"produto_id"`
	Quantidade   float64   `json:"quantidade"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NotaFiscal é a entidade de domínio que representa uma nota fiscal
type NotaFiscal struct {
	ID        uint       `json:"id"`
	Numero    uint       `json:"numero"`
	Status    string     `json:"status"` // "Aberta" ou "Fechada"
	Itens     []ItemNota `json:"itens"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CriarNota cria uma nova nota fiscal
func (n *NotaFiscal) CriarNota() error {
	if len(n.Itens) == 0 {
		return ErrNotaSemItens
	}

	for _, item := range n.Itens {
		if item.Quantidade <= 0 {
			return ErrQuantidadeInvalida
		}
	}

	n.Status = "Aberta"
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	return nil
}

// ImprimirNota marca a nota como fechada (impressa)
func (n *NotaFiscal) ImprimirNota() error {
	if n.Status != "Aberta" {
		return ErrNotaNaoAberta
	}

	n.Status = "Fechada"
	n.UpdatedAt = time.Now()

	return nil
}
