package application

import (
	"backend/internal/domain"
	"fmt"
	"sync"
	"time"
)

// NotaFiscalRepositoryPort define o contrato de repositório usado pelo serviço de notas fiscais.
type NotaFiscalRepositoryPort interface {
	Listar() ([]domain.NotaFiscal, error)
	BuscarPorID(id string) (*domain.NotaFiscal, error)
	Criar(nota *domain.NotaFiscal) error
	AtualizarStatus(id string, status string) error
}

// NotaFiscalServiceImpl implementa a lógica de negócio para notas fiscais.
type NotaFiscalServiceImpl struct {
	notaRepo    NotaFiscalRepositoryPort
	produtoRepo ProdutoRepositoryPort
	mu          sync.Mutex
	nextNumero  int
}

// NovaNotaFiscalService cria uma nova instância do serviço de notas fiscais.
func NovaNotaFiscalService(notaRepo NotaFiscalRepositoryPort, produtoRepo ProdutoRepositoryPort) *NotaFiscalServiceImpl {
	return &NotaFiscalServiceImpl{
		notaRepo:    notaRepo,
		produtoRepo: produtoRepo,
		nextNumero:  1,
	}
}

// ListarNotas retorna todas as notas fiscais cadastradas.
func (s *NotaFiscalServiceImpl) ListarNotas() ([]NotaFiscalDTO, error) {
	notas, err := s.notaRepo.Listar()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar notas fiscais: %w", err)
	}

	dtos := make([]NotaFiscalDTO, len(notas))
	for i, n := range notas {
		dtos[i] = notaParaDTO(n)
	}
	return dtos, nil
}

// CriarNota valida os dados e cria uma nova nota fiscal.
func (s *NotaFiscalServiceImpl) CriarNota(input CriarNotaInput) (*NotaFiscalDTO, error) {
	if len(input.Itens) == 0 {
		return nil, fmt.Errorf("%w: a nota fiscal deve ter pelo menos um item", domain.ErrInvalidInput)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Validar itens
	for _, item := range input.Itens {
		if item.ProdutoID == "" {
			return nil, fmt.Errorf("%w: produto_id é obrigatório para cada item", domain.ErrInvalidInput)
		}
		if item.Quantidade <= 0 {
			return nil, fmt.Errorf("%w: quantidade deve ser maior que zero", domain.ErrInvalidInput)
		}
		// Verificar se o produto existe
		_, err := s.produtoRepo.BuscarPorID(item.ProdutoID)
		if err != nil {
			return nil, fmt.Errorf("produto %s: %w", item.ProdutoID, err)
		}
	}

	agora := time.Now()
	itens := make([]domain.ItemNota, len(input.Itens))
	for i, item := range input.Itens {
		itens[i] = domain.ItemNota{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
		}
	}

	nota := &domain.NotaFiscal{
		Numero:       s.nextNumero,
		Status:       "Aberta",
		Itens:        itens,
		CriadoEm:     agora,
		AtualizadoEm: agora,
	}
	s.nextNumero++

	if err := s.notaRepo.Criar(nota); err != nil {
		return nil, fmt.Errorf("erro ao criar nota fiscal: %w", err)
	}

	dto := notaParaDTO(*nota)
	return &dto, nil
}

// BuscarNotaPorID busca uma nota fiscal pelo seu ID.
func (s *NotaFiscalServiceImpl) BuscarNotaPorID(id string) (*NotaFiscalDTO, error) {
	nota, err := s.notaRepo.BuscarPorID(id)
	if err != nil {
		return nil, err
	}
	dto := notaParaDTO(*nota)
	return &dto, nil
}

// ImprimirNota fecha a nota fiscal e deduz o estoque dos produtos.
func (s *NotaFiscalServiceImpl) ImprimirNota(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	nota, err := s.notaRepo.BuscarPorID(id)
	if err != nil {
		return err
	}

	if nota.Status == "Fechada" {
		return domain.ErrInvoiceAlreadyClosed
	}

	// Verificar estoque suficiente para todos os itens antes de deduzir
	for _, item := range nota.Itens {
		produto, err := s.produtoRepo.BuscarPorID(item.ProdutoID)
		if err != nil {
			return fmt.Errorf("produto %s: %w", item.ProdutoID, err)
		}
		if produto.Saldo < item.Quantidade {
			return fmt.Errorf("%w: produto %s (saldo: %.2f, solicitado: %.2f)",
				domain.ErrInsufficientStock, produto.Codigo, produto.Saldo, item.Quantidade)
		}
	}

	// Deduzir estoque
	for _, item := range nota.Itens {
		produto, _ := s.produtoRepo.BuscarPorID(item.ProdutoID)
		novoSaldo := produto.Saldo - item.Quantidade
		if err := s.produtoRepo.AtualizarSaldo(item.ProdutoID, novoSaldo); err != nil {
			return fmt.Errorf("erro ao atualizar saldo do produto %s: %w", item.ProdutoID, err)
		}
	}

	// Atualizar status da nota
	if err := s.notaRepo.AtualizarStatus(id, "Fechada"); err != nil {
		return fmt.Errorf("erro ao fechar nota fiscal: %w", err)
	}

	return nil
}

// notaParaDTO converte uma entidade NotaFiscal para NotaFiscalDTO.
func notaParaDTO(n domain.NotaFiscal) NotaFiscalDTO {
	itens := make([]ItemNotaDTO, len(n.Itens))
	for i, item := range n.Itens {
		itens[i] = ItemNotaDTO{
			ProdutoID:  item.ProdutoID,
			Quantidade: item.Quantidade,
		}
	}
	return NotaFiscalDTO{
		ID:           n.ID,
		Numero:       n.Numero,
		Status:       n.Status,
		Itens:        itens,
		CriadoEm:     n.CriadoEm,
		AtualizadoEm: n.AtualizadoEm,
	}
}
