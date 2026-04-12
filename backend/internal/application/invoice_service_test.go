package application_test

import (
	"backend/internal/application"
	"backend/internal/domain"
	"backend/internal/infrastructure/repository"
	"errors"
	"testing"
)

func setupNotaService() (*application.NotaFiscalServiceImpl, *application.ProdutoServiceImpl) {
	produtoRepo := repository.NovoProdutoRepositoryMemoria()
	notaRepo := repository.NovaNotaFiscalRepositoryMemoria()
	produtoService := application.NovoProdutoService(produtoRepo)
	notaService := application.NovaNotaFiscalService(notaRepo, produtoRepo)
	return notaService, produtoService
}

func TestCriarNotaSucesso(t *testing.T) {
	notaService, produtoService := setupNotaService()

	// Criar produto primeiro
	produto, _ := produtoService.CriarProduto(application.CriarProdutoInput{
		Codigo: "NF-PROD-001", Descricao: "Produto NF", Saldo: 100,
	})

	input := application.CriarNotaInput{
		Itens: []application.ItemNotaInput{
			{ProdutoID: produto.ID, Quantidade: 10},
		},
	}

	nota, err := notaService.CriarNota(input)
	if err != nil {
		t.Fatalf("erro inesperado ao criar nota: %v", err)
	}

	if nota.Status != "Aberta" {
		t.Errorf("esperado status 'Aberta', obteve '%s'", nota.Status)
	}
	if nota.Numero != 1 {
		t.Errorf("esperado número 1, obteve %d", nota.Numero)
	}
	if len(nota.Itens) != 1 {
		t.Errorf("esperado 1 item, obteve %d", len(nota.Itens))
	}
}

func TestCriarNotaSemItens(t *testing.T) {
	notaService, _ := setupNotaService()

	input := application.CriarNotaInput{
		Itens: []application.ItemNotaInput{},
	}

	_, err := notaService.CriarNota(input)
	if err == nil {
		t.Fatal("esperado erro ao criar nota sem itens")
	}
}

func TestCriarNotaProdutoInexistente(t *testing.T) {
	notaService, _ := setupNotaService()

	input := application.CriarNotaInput{
		Itens: []application.ItemNotaInput{
			{ProdutoID: "nao-existe", Quantidade: 10},
		},
	}

	_, err := notaService.CriarNota(input)
	if err == nil {
		t.Fatal("esperado erro ao criar nota com produto inexistente")
	}
}

func TestCriarNotaQuantidadeZero(t *testing.T) {
	notaService, produtoService := setupNotaService()

	produto, _ := produtoService.CriarProduto(application.CriarProdutoInput{
		Codigo: "NF-PROD-002", Descricao: "Produto", Saldo: 50,
	})

	input := application.CriarNotaInput{
		Itens: []application.ItemNotaInput{
			{ProdutoID: produto.ID, Quantidade: 0},
		},
	}

	_, err := notaService.CriarNota(input)
	if err == nil {
		t.Fatal("esperado erro ao criar nota com quantidade zero")
	}
}

func TestAutoIncrementoNumero(t *testing.T) {
	notaService, produtoService := setupNotaService()

	produto, _ := produtoService.CriarProduto(application.CriarProdutoInput{
		Codigo: "NF-PROD-003", Descricao: "Produto", Saldo: 200,
	})

	item := []application.ItemNotaInput{{ProdutoID: produto.ID, Quantidade: 5}}

	nota1, _ := notaService.CriarNota(application.CriarNotaInput{Itens: item})
	nota2, _ := notaService.CriarNota(application.CriarNotaInput{Itens: item})

	if nota1.Numero != 1 {
		t.Errorf("esperado número 1 para primeira nota, obteve %d", nota1.Numero)
	}
	if nota2.Numero != 2 {
		t.Errorf("esperado número 2 para segunda nota, obteve %d", nota2.Numero)
	}
}

func TestImprimirNotaSucesso(t *testing.T) {
	notaService, produtoService := setupNotaService()

	produto, _ := produtoService.CriarProduto(application.CriarProdutoInput{
		Codigo: "IMP-001", Descricao: "Produto Impressão", Saldo: 100,
	})

	nota, _ := notaService.CriarNota(application.CriarNotaInput{
		Itens: []application.ItemNotaInput{
			{ProdutoID: produto.ID, Quantidade: 30},
		},
	})

	err := notaService.ImprimirNota(nota.ID)
	if err != nil {
		t.Fatalf("erro inesperado ao imprimir nota: %v", err)
	}

	// Verificar que a nota foi fechada
	notaAtualizada, _ := notaService.BuscarNotaPorID(nota.ID)
	if notaAtualizada.Status != "Fechada" {
		t.Errorf("esperado status 'Fechada', obteve '%s'", notaAtualizada.Status)
	}

	// Verificar que o estoque foi deduzido
	produtoAtualizado, _ := produtoService.BuscarProdutoPorID(produto.ID)
	if produtoAtualizado.Saldo != 70 {
		t.Errorf("esperado saldo 70 após dedução, obteve %f", produtoAtualizado.Saldo)
	}
}

func TestImprimirNotaJaFechada(t *testing.T) {
	notaService, produtoService := setupNotaService()

	produto, _ := produtoService.CriarProduto(application.CriarProdutoInput{
		Codigo: "IMP-002", Descricao: "Produto", Saldo: 100,
	})

	nota, _ := notaService.CriarNota(application.CriarNotaInput{
		Itens: []application.ItemNotaInput{
			{ProdutoID: produto.ID, Quantidade: 10},
		},
	})

	// Primeira impressão
	notaService.ImprimirNota(nota.ID)

	// Segunda impressão deve falhar
	err := notaService.ImprimirNota(nota.ID)
	if err == nil {
		t.Fatal("esperado erro ao imprimir nota já fechada")
	}
	if !errors.Is(err, domain.ErrInvoiceAlreadyClosed) {
		t.Errorf("esperado ErrInvoiceAlreadyClosed, obteve %v", err)
	}
}

func TestImprimirNotaEstoqueInsuficiente(t *testing.T) {
	notaService, produtoService := setupNotaService()

	produto, _ := produtoService.CriarProduto(application.CriarProdutoInput{
		Codigo: "IMP-003", Descricao: "Produto", Saldo: 5,
	})

	nota, _ := notaService.CriarNota(application.CriarNotaInput{
		Itens: []application.ItemNotaInput{
			{ProdutoID: produto.ID, Quantidade: 10},
		},
	})

	err := notaService.ImprimirNota(nota.ID)
	if err == nil {
		t.Fatal("esperado erro ao imprimir nota com estoque insuficiente")
	}
	if !errors.Is(err, domain.ErrInsufficientStock) {
		t.Errorf("esperado ErrInsufficientStock, obteve %v", err)
	}
}

func TestImprimirNotaInexistente(t *testing.T) {
	notaService, _ := setupNotaService()

	err := notaService.ImprimirNota("nao-existe")
	if err == nil {
		t.Fatal("esperado erro ao imprimir nota inexistente")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("esperado ErrNotFound, obteve %v", err)
	}
}

func TestListarNotas(t *testing.T) {
	notaService, produtoService := setupNotaService()

	produto, _ := produtoService.CriarProduto(application.CriarProdutoInput{
		Codigo: "LIST-001", Descricao: "Produto", Saldo: 200,
	})

	item := []application.ItemNotaInput{{ProdutoID: produto.ID, Quantidade: 5}}
	notaService.CriarNota(application.CriarNotaInput{Itens: item})
	notaService.CriarNota(application.CriarNotaInput{Itens: item})

	notas, err := notaService.ListarNotas()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(notas) != 2 {
		t.Errorf("esperado 2 notas, obteve %d", len(notas))
	}
}

func TestBuscarNotaPorID(t *testing.T) {
	notaService, produtoService := setupNotaService()

	produto, _ := produtoService.CriarProduto(application.CriarProdutoInput{
		Codigo: "BUSCA-NF-001", Descricao: "Produto", Saldo: 100,
	})

	criada, _ := notaService.CriarNota(application.CriarNotaInput{
		Itens: []application.ItemNotaInput{
			{ProdutoID: produto.ID, Quantidade: 5},
		},
	})

	encontrada, err := notaService.BuscarNotaPorID(criada.ID)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if encontrada.ID != criada.ID {
		t.Errorf("esperado ID %s, obteve %s", criada.ID, encontrada.ID)
	}
}

func TestImprimirNotaMultiplosItens(t *testing.T) {
	notaService, produtoService := setupNotaService()

	prod1, _ := produtoService.CriarProduto(application.CriarProdutoInput{
		Codigo: "MULTI-001", Descricao: "Produto 1", Saldo: 100,
	})
	prod2, _ := produtoService.CriarProduto(application.CriarProdutoInput{
		Codigo: "MULTI-002", Descricao: "Produto 2", Saldo: 50,
	})

	nota, _ := notaService.CriarNota(application.CriarNotaInput{
		Itens: []application.ItemNotaInput{
			{ProdutoID: prod1.ID, Quantidade: 20},
			{ProdutoID: prod2.ID, Quantidade: 15},
		},
	})

	err := notaService.ImprimirNota(nota.ID)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Verificar saldos atualizados
	p1, _ := produtoService.BuscarProdutoPorID(prod1.ID)
	p2, _ := produtoService.BuscarProdutoPorID(prod2.ID)

	if p1.Saldo != 80 {
		t.Errorf("esperado saldo 80 para produto 1, obteve %f", p1.Saldo)
	}
	if p2.Saldo != 35 {
		t.Errorf("esperado saldo 35 para produto 2, obteve %f", p2.Saldo)
	}
}
