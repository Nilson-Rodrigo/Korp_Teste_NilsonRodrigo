package repository

import (
	"backend/internal/domain"
	"fmt"
	"sync"
	"time"
)

// NotaFiscalRepositoryMemoria implementa o repositório de notas fiscais em memória.
type NotaFiscalRepositoryMemoria struct {
	mu     sync.RWMutex
	notas  map[string]domain.NotaFiscal
	nextID int
}

// NovaNotaFiscalRepositoryMemoria cria uma nova instância do repositório em memória.
func NovaNotaFiscalRepositoryMemoria() *NotaFiscalRepositoryMemoria {
	return &NotaFiscalRepositoryMemoria{
		notas:  make(map[string]domain.NotaFiscal),
		nextID: 1,
	}
}

// gerarID gera um ID único baseado em contador.
func (r *NotaFiscalRepositoryMemoria) gerarID() string {
	id := fmt.Sprintf("nf-%d", r.nextID)
	r.nextID++
	return id
}

// Listar retorna todas as notas fiscais armazenadas.
func (r *NotaFiscalRepositoryMemoria) Listar() ([]domain.NotaFiscal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	notas := make([]domain.NotaFiscal, 0, len(r.notas))
	for _, n := range r.notas {
		notas = append(notas, n)
	}
	return notas, nil
}

// BuscarPorID busca uma nota fiscal pelo ID.
func (r *NotaFiscalRepositoryMemoria) BuscarPorID(id string) (*domain.NotaFiscal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	nota, ok := r.notas[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &nota, nil
}

// Criar insere uma nova nota fiscal no repositório.
func (r *NotaFiscalRepositoryMemoria) Criar(nota *domain.NotaFiscal) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	nota.ID = r.gerarID()
	r.notas[nota.ID] = *nota
	return nil
}

// AtualizarStatus atualiza o status de uma nota fiscal existente.
func (r *NotaFiscalRepositoryMemoria) AtualizarStatus(id string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	nota, ok := r.notas[id]
	if !ok {
		return domain.ErrNotFound
	}

	nota.Status = status
	nota.AtualizadoEm = time.Now()
	r.notas[id] = nota
	return nil
}
