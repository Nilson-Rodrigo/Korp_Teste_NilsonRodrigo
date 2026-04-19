package services

import (
	"bytes"
	"encoding/json"
	"faturamento/domain/entities"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// EstoqueServiceImpl implementa a interface EstoqueService
type EstoqueServiceImpl struct {
	baseURL string
	client  *http.Client
}

// NewEstoqueService cria uma nova instância do serviço de estoque
func NewEstoqueService(baseURL string) *EstoqueServiceImpl {
	return &EstoqueServiceImpl{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type estoqueProduto struct {
	ID    uint    `json:"ID"`
	Saldo float64 `json:"saldo"`
}

// BuscarProduto busca um produto no estoque
func (s *EstoqueServiceImpl) BuscarProduto(id uint) (interface{}, error) {
	url := fmt.Sprintf("%s/produtos/%d", s.baseURL, id)

	resp, err := s.client.Get(url)
	if err != nil {
		log.Error().Err(err).Uint("produto_id", id).Msg("Erro ao conectar ao serviço de estoque")
		return nil, entities.ErrEstoqueIndisponivel
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().Int("status_code", resp.StatusCode).Uint("produto_id", id).Msg("Produto não encontrado no estoque")
		return nil, entities.ErrProdutoIndisponivel
	}

	var produto estoqueProduto
	if err := json.NewDecoder(resp.Body).Decode(&produto); err != nil {
		log.Error().Err(err).Msg("Erro ao desserializar produto")
		return nil, entities.ErrEstoqueIndisponivel
	}

	log.Debug().Uint("produto_id", id).Float64("saldo", produto.Saldo).Msg("Produto encontrado no estoque")
	return &produto, nil
}

// BuscarSaldo busca o saldo de um produto
func (s *EstoqueServiceImpl) BuscarSaldo(id uint) (float64, error) {
	produto, err := s.BuscarProduto(id)
	if err != nil {
		return 0, err
	}

	if p, ok := produto.(*estoqueProduto); ok {
		return p.Saldo, nil
	}

	return 0, entities.ErrEstoqueIndisponivel
}

// AtualizarSaldo atualiza o saldo de um produto no estoque
func (s *EstoqueServiceImpl) AtualizarSaldo(id uint, novoSaldo float64) error {
	url := fmt.Sprintf("%s/produtos/%d/saldo", s.baseURL, id)

	payload := map[string]float64{"saldo": novoSaldo}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("Erro ao serializar payload")
		return entities.ErrEstoqueIndisponivel
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		log.Error().Err(err).Msg("Erro ao criar request")
		return entities.ErrEstoqueIndisponivel
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		log.Error().Err(err).Uint("produto_id", id).Msg("Erro ao conectar ao serviço de estoque")
		return entities.ErrEstoqueIndisponivel
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status_code", resp.StatusCode).Uint("produto_id", id).Msg("Erro ao atualizar saldo no estoque")
		return entities.ErrEstoqueIndisponivel
	}

	log.Info().Uint("produto_id", id).Float64("novo_saldo", novoSaldo).Msg("Saldo atualizado no estoque")
	return nil
}
