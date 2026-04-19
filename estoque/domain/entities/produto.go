package entities

import "time"

// Produto é a entidade de domínio que representa um produto em estoque
type Produto struct {
	ID        uint      `json:"id"`
	Codigo    string    `json:"codigo"`
	Descricao string    `json:"descricao"`
	Saldo     float64   `json:"saldo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validar valida se o produto tem campos obrigatórios
func (p *Produto) Validar() error {
	if p.Codigo == "" {
		return ErrCodigoObrigatorio
	}
	if p.Descricao == "" {
		return ErrDescricaoObrigatoria
	}
	if p.Saldo < 0 {
		return ErrSaldoNegativo
	}
	return nil
}

// AtualizarSaldo atualiza o saldo do produto
func (p *Produto) AtualizarSaldo(quantidade float64) error {
	novoSaldo := p.Saldo - quantidade
	if novoSaldo < 0 {
		return ErrSaldoInsuficiente
	}
	p.Saldo = novoSaldo
	p.UpdatedAt = time.Now()
	return nil
}
