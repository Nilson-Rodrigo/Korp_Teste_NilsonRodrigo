package application

import (
	"backend/internal/domain"
	"fmt"
	"strings"
	"sync"
	"time"
)

// ProdutoRepositoryPort define o contrato de repositório usado pelo serviço de produtos.
type ProdutoRepositoryPort interface {
	Listar() ([]domain.Produto, error)
	BuscarPorID(id string) (*domain.Produto, error)
	Criar(produto *domain.Produto) error
	AtualizarSaldo(id string, novoSaldo float64) error
}

// ProdutoServiceImpl implementa a lógica de negócio para produtos.
type ProdutoServiceImpl struct {
	repo ProdutoRepositoryPort
	mu   sync.Mutex
}

// NovoProdutoService cria uma nova instância do serviço de produtos.
func NovoProdutoService(repo ProdutoRepositoryPort) *ProdutoServiceImpl {
	return &ProdutoServiceImpl{repo: repo}
}

// ListarProdutos retorna todos os produtos cadastrados.
func (s *ProdutoServiceImpl) ListarProdutos() ([]ProdutoDTO, error) {
	produtos, err := s.repo.Listar()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar produtos: %w", err)
	}

	dtos := make([]ProdutoDTO, len(produtos))
	for i, p := range produtos {
		dtos[i] = produtoParaDTO(p)
	}
	return dtos, nil
}

// CriarProduto valida os dados e cria um novo produto.
func (s *ProdutoServiceImpl) CriarProduto(input CriarProdutoInput) (*ProdutoDTO, error) {
	if strings.TrimSpace(input.Codigo) == "" {
		return nil, fmt.Errorf("%w: código é obrigatório", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(input.Descricao) == "" {
		return nil, fmt.Errorf("%w: descrição é obrigatória", domain.ErrInvalidInput)
	}
	if input.Saldo < 0 {
		return nil, fmt.Errorf("%w: saldo não pode ser negativo", domain.ErrInvalidInput)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Verificar código duplicado
	produtos, err := s.repo.Listar()
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar código duplicado: %w", err)
	}
	for _, p := range produtos {
		if p.Codigo == input.Codigo {
			return nil, domain.ErrDuplicateCode
		}
	}

	agora := time.Now()
	produto := &domain.Produto{
		Codigo:       input.Codigo,
		Descricao:    input.Descricao,
		Saldo:        input.Saldo,
		CriadoEm:     agora,
		AtualizadoEm: agora,
	}

	if err := s.repo.Criar(produto); err != nil {
		return nil, fmt.Errorf("erro ao criar produto: %w", err)
	}

	dto := produtoParaDTO(*produto)
	return &dto, nil
}

// BuscarProdutoPorID busca um produto pelo seu ID.
func (s *ProdutoServiceImpl) BuscarProdutoPorID(id string) (*ProdutoDTO, error) {
	produto, err := s.repo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}
	dto := produtoParaDTO(*produto)
	return &dto, nil
}

// AtualizarSaldo atualiza o saldo de um produto.
func (s *ProdutoServiceImpl) AtualizarSaldo(id string, novoSaldo float64) error {
	if novoSaldo < 0 {
		return fmt.Errorf("%w: saldo não pode ser negativo", domain.ErrInvalidInput)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.repo.BuscarPorID(id)
	if err != nil {
		return err
	}

	return s.repo.AtualizarSaldo(id, novoSaldo)
}

// produtoParaDTO converte uma entidade Produto para ProdutoDTO.
func produtoParaDTO(p domain.Produto) ProdutoDTO {
	return ProdutoDTO{
		ID:           p.ID,
		Codigo:       p.Codigo,
		Descricao:    p.Descricao,
		Saldo:        p.Saldo,
		CriadoEm:     p.CriadoEm,
		AtualizadoEm: p.AtualizadoEm,
	}
}
