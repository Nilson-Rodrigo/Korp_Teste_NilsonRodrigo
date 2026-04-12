package application_test

import (
	"backend/internal/application"
	"backend/internal/domain"
	"backend/internal/infrastructure/repository"
	"testing"
)

func setupProdutoService() *application.ProdutoServiceImpl {
	repo := repository.NovoProdutoRepositoryMemoria()
	return application.NovoProdutoService(repo)
}

func TestCriarProdutoSucesso(t *testing.T) {
	service := setupProdutoService()

	input := application.CriarProdutoInput{
		Codigo:    "TEST-001",
		Descricao: "Produto de Teste",
		Saldo:     100,
	}

	produto, err := service.CriarProduto(input)
	if err != nil {
		t.Fatalf("erro inesperado ao criar produto: %v", err)
	}

	if produto.Codigo != "TEST-001" {
		t.Errorf("esperado código TEST-001, obteve %s", produto.Codigo)
	}
	if produto.Descricao != "Produto de Teste" {
		t.Errorf("esperado descrição 'Produto de Teste', obteve %s", produto.Descricao)
	}
	if produto.Saldo != 100 {
		t.Errorf("esperado saldo 100, obteve %f", produto.Saldo)
	}
	if produto.ID == "" {
		t.Error("esperado ID não vazio")
	}
}

func TestCriarProdutoCodigoVazio(t *testing.T) {
	service := setupProdutoService()

	input := application.CriarProdutoInput{
		Codigo:    "",
		Descricao: "Produto",
		Saldo:     10,
	}

	_, err := service.CriarProduto(input)
	if err == nil {
		t.Fatal("esperado erro ao criar produto com código vazio")
	}
}

func TestCriarProdutoDescricaoVazia(t *testing.T) {
	service := setupProdutoService()

	input := application.CriarProdutoInput{
		Codigo:    "TEST-002",
		Descricao: "",
		Saldo:     10,
	}

	_, err := service.CriarProduto(input)
	if err == nil {
		t.Fatal("esperado erro ao criar produto com descrição vazia")
	}
}

func TestCriarProdutoSaldoNegativo(t *testing.T) {
	service := setupProdutoService()

	input := application.CriarProdutoInput{
		Codigo:    "TEST-003",
		Descricao: "Produto",
		Saldo:     -10,
	}

	_, err := service.CriarProduto(input)
	if err == nil {
		t.Fatal("esperado erro ao criar produto com saldo negativo")
	}
}

func TestCriarProdutoCodigoDuplicado(t *testing.T) {
	service := setupProdutoService()

	input := application.CriarProdutoInput{
		Codigo:    "TEST-DUP",
		Descricao: "Produto 1",
		Saldo:     10,
	}
	_, err := service.CriarProduto(input)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	input2 := application.CriarProdutoInput{
		Codigo:    "TEST-DUP",
		Descricao: "Produto 2",
		Saldo:     20,
	}
	_, err = service.CriarProduto(input2)
	if err == nil {
		t.Fatal("esperado erro ao criar produto com código duplicado")
	}
}

func TestListarProdutos(t *testing.T) {
	service := setupProdutoService()

	// Criar dois produtos
	service.CriarProduto(application.CriarProdutoInput{Codigo: "A", Descricao: "Produto A", Saldo: 10})
	service.CriarProduto(application.CriarProdutoInput{Codigo: "B", Descricao: "Produto B", Saldo: 20})

	produtos, err := service.ListarProdutos()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(produtos) != 2 {
		t.Errorf("esperado 2 produtos, obteve %d", len(produtos))
	}
}

func TestBuscarProdutoPorID(t *testing.T) {
	service := setupProdutoService()

	criado, _ := service.CriarProduto(application.CriarProdutoInput{
		Codigo: "BUSCA-001", Descricao: "Produto Busca", Saldo: 50,
	})

	encontrado, err := service.BuscarProdutoPorID(criado.ID)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if encontrado.ID != criado.ID {
		t.Errorf("esperado ID %s, obteve %s", criado.ID, encontrado.ID)
	}
}

func TestBuscarProdutoNaoExistente(t *testing.T) {
	service := setupProdutoService()

	_, err := service.BuscarProdutoPorID("nao-existe")
	if err == nil {
		t.Fatal("esperado erro ao buscar produto inexistente")
	}
	if err != domain.ErrNotFound {
		t.Errorf("esperado ErrNotFound, obteve %v", err)
	}
}

func TestAtualizarSaldo(t *testing.T) {
	service := setupProdutoService()

	criado, _ := service.CriarProduto(application.CriarProdutoInput{
		Codigo: "SALDO-001", Descricao: "Produto Saldo", Saldo: 100,
	})

	err := service.AtualizarSaldo(criado.ID, 75)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	atualizado, _ := service.BuscarProdutoPorID(criado.ID)
	if atualizado.Saldo != 75 {
		t.Errorf("esperado saldo 75, obteve %f", atualizado.Saldo)
	}
}

func TestAtualizarSaldoNegativo(t *testing.T) {
	service := setupProdutoService()

	criado, _ := service.CriarProduto(application.CriarProdutoInput{
		Codigo: "SALDO-002", Descricao: "Produto", Saldo: 100,
	})

	err := service.AtualizarSaldo(criado.ID, -10)
	if err == nil {
		t.Fatal("esperado erro ao atualizar saldo para valor negativo")
	}
}

func TestAtualizarSaldoProdutoInexistente(t *testing.T) {
	service := setupProdutoService()

	err := service.AtualizarSaldo("nao-existe", 50)
	if err == nil {
		t.Fatal("esperado erro ao atualizar saldo de produto inexistente")
	}
}
