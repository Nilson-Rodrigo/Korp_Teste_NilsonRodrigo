package repository

import (
	"backend/internal/domain"
	"fmt"
	"sync"
	"time"
)

// ProdutoRepositoryMemoria implementa o repositório de produtos em memória.
type ProdutoRepositoryMemoria struct {
	mu        sync.RWMutex
	produtos  map[string]domain.Produto
	nextID    int
}

// NovoProdutoRepositoryMemoria cria uma nova instância do repositório em memória.
func NovoProdutoRepositoryMemoria() *ProdutoRepositoryMemoria {
	return &ProdutoRepositoryMemoria{
		produtos: make(map[string]domain.Produto),
		nextID:   1,
	}
}

// gerarID gera um ID único baseado em contador.
func (r *ProdutoRepositoryMemoria) gerarID() string {
	id := fmt.Sprintf("prod-%d", r.nextID)
	r.nextID++
	return id
}

// Listar retorna todos os produtos armazenados.
func (r *ProdutoRepositoryMemoria) Listar() ([]domain.Produto, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	produtos := make([]domain.Produto, 0, len(r.produtos))
	for _, p := range r.produtos {
		produtos = append(produtos, p)
	}
	return produtos, nil
}

// BuscarPorID busca um produto pelo ID.
func (r *ProdutoRepositoryMemoria) BuscarPorID(id string) (*domain.Produto, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	produto, ok := r.produtos[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &produto, nil
}

// Criar insere um novo produto no repositório.
func (r *ProdutoRepositoryMemoria) Criar(produto *domain.Produto) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	produto.ID = r.gerarID()
	r.produtos[produto.ID] = *produto
	return nil
}

// AtualizarSaldo atualiza o saldo de um produto existente.
func (r *ProdutoRepositoryMemoria) AtualizarSaldo(id string, novoSaldo float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	produto, ok := r.produtos[id]
	if !ok {
		return domain.ErrNotFound
	}

	produto.Saldo = novoSaldo
	produto.AtualizadoEm = time.Now()
	r.produtos[id] = produto
	return nil
}
