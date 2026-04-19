package dto

// CriarProdutoRequest é o DTO para criar um produto
type CriarProdutoRequest struct {
	Codigo    string  `json:"codigo" binding:"required,min=1,max=50"`
	Descricao string  `json:"descricao" binding:"required,min=3,max=500"`
	Saldo     float64 `json:"saldo" binding:"required,min=0,max=999999.99"`
}

// AtualizarSaldoRequest é o DTO para atualizar o saldo
type AtualizarSaldoRequest struct {
	Saldo float64 `json:"saldo" binding:"required,min=0,max=999999.99"`
}

// ProdutoResponse é o DTO de resposta
type ProdutoResponse struct {
	ID        uint    `json:"id"`
	Codigo    string  `json:"codigo"`
	Descricao string  `json:"descricao"`
	Saldo     float64 `json:"saldo"`
}
