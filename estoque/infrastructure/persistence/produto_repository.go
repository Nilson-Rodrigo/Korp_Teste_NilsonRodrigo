package persistence

import (
	"estoque/domain/entities"

	"gorm.io/gorm"
)

// ProdutoRepositoryImpl implementa a interface ProdutoRepository
type ProdutoRepositoryImpl struct {
	db *gorm.DB
}

// NewProdutoRepository cria uma nova instância do repositório
func NewProdutoRepository(db *gorm.DB) *ProdutoRepositoryImpl {
	return &ProdutoRepositoryImpl{
		db: db,
	}
}

// Criar cria um novo produto no banco de dados
func (r *ProdutoRepositoryImpl) Criar(produto *entities.Produto) error {
	return r.db.Create(produto).Error
}

// BuscarPorID busca um produto por ID
func (r *ProdutoRepositoryImpl) BuscarPorID(id uint) (*entities.Produto, error) {
	var produto entities.Produto
	if err := r.db.First(&produto, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrProdutoNaoEncontrado
		}
		return nil, err
	}
	return &produto, nil
}

// BuscarPorCodigo busca um produto por código
func (r *ProdutoRepositoryImpl) BuscarPorCodigo(codigo string) (*entities.Produto, error) {
	var produto entities.Produto
	if err := r.db.Where("codigo = ?", codigo).First(&produto).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrProdutoNaoEncontrado
		}
		return nil, err
	}
	return &produto, nil
}

// Listar lista todos os produtos
func (r *ProdutoRepositoryImpl) Listar() ([]entities.Produto, error) {
	var produtos []entities.Produto
	if err := r.db.Find(&produtos).Error; err != nil {
		return nil, err
	}
	return produtos, nil
}

// AtualizarSaldo atualiza o saldo de um produto
func (r *ProdutoRepositoryImpl) AtualizarSaldo(id uint, novoSaldo float64) error {
	return r.db.Model(&entities.Produto{}).Where("id = ?", id).Update("saldo", novoSaldo).Error
}

// Deletar deleta um produto
func (r *ProdutoRepositoryImpl) Deletar(id uint) error {
	return r.db.Delete(&entities.Produto{}, id).Error
}
