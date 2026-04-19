package dto

// ItemNotaRequest é o DTO para item de nota
type ItemNotaRequest struct {
	ProdutoID  uint    `json:"produto_id" binding:"required"`
	Quantidade float64 `json:"quantidade" binding:"required,gt=0,max=999999.99"`
}

// CriarNotaRequest é o DTO para criar uma nota
type CriarNotaRequest struct {
	Itens []ItemNotaRequest `json:"itens" binding:"required,min=1"`
}

// ItemNotaResponse é o DTO de resposta para item
type ItemNotaResponse struct {
	ID         uint    `json:"id"`
	ProdutoID  uint    `json:"produto_id"`
	Quantidade float64 `json:"quantidade"`
}

// NotaFiscalResponse é o DTO de resposta para nota
type NotaFiscalResponse struct {
	ID     uint               `json:"id"`
	Numero uint               `json:"numero"`
	Status string             `json:"status"`
	Itens  []ItemNotaResponse `json:"itens"`
}
