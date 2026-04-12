# DETALHAMENTO TÉCNICO - Sistema de Notas Fiscais Korp

Documento de suporte para apresentação em vídeo com todos os detalhes técnicos solicitados.

---

## 1. CICLOS DE VIDA DO ANGULAR UTILIZADOS

### 1.1 OnInit
**Classe:** `Produtos`, `Notas`

Executado uma vez após o componente ser inicializado:

```typescript
ngOnInit() {
  this.carregarNotas();
  this.carregarProdutos();
}
```

**Por que usamos:**
- Carregar dados iniciais da API
- Inicializar formulários
- Configurar listeners globais

### 1.2 OnDestroy
**Classe:** `Produtos`, `Notas`

Executado antes do componente ser destruído:

```typescript
ngOnDestroy(): void {
  this.destroy$.next();
  this.destroy$.complete();
}
```

**Por que usamos:**
- Limpar subscriptions (evitar memory leaks)
- Cancelar requisições HTTP em progresso
- Liberar recursos (timers, listeners)

### 1.3 OnChanges (Não utilizado, mas disponível)
Detecta mudanças em @Input() properties

---

## 2. BIBLIOTECA RXJS E COMO FOI UTILIZADA

### 2.1 Operadores Principais

#### catchError - Tratamento de Erros
```typescript
this.produtoService.listar()
  .pipe(
    catchError((err) => {
      // Tratamento do erro
      this.snackBar.open('Erro ao carregar', 'Fechar', {
        duration: 5000,
        panelClass: 'snack-error',
      });
      return of([]); // Retorna array vazio como fallback
    })
  )
  .subscribe((data) => {
    this.produtos = data;
  });
```

#### finalize - Cleanup Garantido
```typescript
this.notaService.criar(data)
  .pipe(
    finalize(() => {
      // Executa INDEPENDENTE de sucesso ou erro
      this.salvando = false; // Remover loading
    })
  )
  .subscribe(res => { /* ... */ });
```

#### takeUntil - Cancelamento de Subscription
```typescript
private destroy$ = new Subject<void>();

// Em cada subscription:
this.produtoService.listar()
  .pipe(
    takeUntil(this.destroy$) // Cancela quando destroy$ emite
  )
  .subscribe(data => { /* ... */ });

// Quando componente é destruído:
ngOnDestroy() {
  this.destroy$.next();
}
```

### 2.2 Padrão de Composition Pipe
Todos os operadores são compostos com **pipe()**:

```typescript
observable$
  .pipe(
    operator1(),
    operator2(),
    operator3()
  )
  .subscribe(result => {});
```

**Benefícios:**
- Cadeia de transformações declarativa
- Lazy evaluation (executa apenas na subscription)
- Fácil de testar e componentizar

### 2.3 Uso em Componentes

**Pattern Geral (Reativo):**
```typescript
// Componente está sempre "escutando" atualizações
this.produtoService.listar()
  .pipe(
    takeUntil(this.destroy$)
  )
  .subscribe(produtos => {
    this.produtos = produtos; // Template se atualiza automaticamente
  });
```

---

## 3. OUTRAS BIBLIOTECAS UTILIZADAS

### Frontend

#### @angular/material
**Versão:** 17
**Utilidade:** Design System / Components visuais

**Componentes usados:**
- `MatTableModule` - Tabelas de dados
- `MatFormFieldModule` - Labels e campos de formulário  
- `MatInputModule` - Inputs de texto
- `MatSelectModule` - Seleção com dropdown
- `MatButtonModule` - Botões estilizados
- `MatCardModule` - Cards/Containers
- `MatSnackBarModule` - Notificações (toasts)
- `MatIconModule` - Ícones vetoriais
- `MatProgressSpinnerModule` - Loading spinners

**Exemplo:**
```html
<mat-form-field appearance="outline">
  <mat-label>Descrição</mat-label>
  <input matInput [(ngModel)]="novoProduto.descricao" />
</mat-form-field>
```

#### @angular/common
**Utilidade:** Pipes e Diretivas

**Diretivas usadas:**
- `*ngIf` - Renderização condicional
- `*ngFor` - Loops
- `(click)` - Event binding
- `[(ngModel)]` - Two-way data binding

**Pipes usados:**
- `| number:'1.0-2'` - Formatação de números com 2 casas decimais

#### @angular/forms
**Utilidade:** Validação e binding de formulários

- `FormsModule` - ngModel e two-way binding
- `FormBuilder` - (não usado neste projeto, mas disponível)

---

## 4. BIBLIOTECAS VISUAIS / COMPONENTES

### Material Design System
Todo o design visual implementado com **Angular Material**.

**Estrutura Visual:**
- **Navbar** - Header sticky com navegação
- **Cards** - Containers para formulários e tabelas
- **Tables** - Listagem de dados com linhas interativas
- **Buttons** - Botões primários, secundários, com icons
- **Badges** - Status indicators (Aberta/Fechada)
- **SnackBar** - Notificações de sucesso/erro
- **Spinner** - Indicador de processamento

**Tema:** Customizado com SCSS
```scss
// Cores customizadas
$primary: #2563eb;
$accent: #f97316;
$warn: #ef4444;
```

---

## 5. GERENCIAMENTO DE DEPENDÊNCIAS NO GOLANG

### 5.1 go.mod
Arquivo de manifesto de dependências (similar a package.json no npm).

**Exemplo (estoque/go.mod):**
```
module estoque

go 1.21

require (
  github.com/gin-gonic/gin v1.9.1
  gorm.io/driver/postgres v1.5.2
  gorm.io/gorm v1.25.4
  github.com/joho/godotenv v1.5.1
)
```

### 5.2 Comandos principais
```bash
go mod download     # Baixa dependências especificadas em go.sum
go mod tidy         # Remove não-utilizadas, adiciona faltantes
go mod graph        # Mostra árvore de dependências
go get -u ./...     # Atualiza para latest version
```

### 5.3 Resolução de Versões
Go usa **semantic versioning** (v1.2.3):
- `v1` = Major version (breaking changes)
- `.2` = Minor (new features, backwards compatible)
- `.3` = Patch (bug fixes)

---

## 6. FRAMEWORKS UTILIZADOS NO GOLANG

### 6.1 Gin Web Framework
**URL:** github.com/gin-gonic/gin

**Utilidade:** HTTP Web Server / API REST

**Recursos utilizados:**

#### Roteamento
```go
r := gin.Default()

r.GET("/produtos", handler.ListarProdutos)
r.POST("/produtos", handler.CriarProduto)
r.PATCH("/produtos/:id/saldo", handler.AtualizarSaldo)
```

#### Middleware CORS
```go
r.Use(func(c *gin.Context) {
  c.Header("Access-Control-Allow-Origin", "*")
  c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
  if c.Request.Method == "OPTIONS" {
    c.AbortWithStatus(204)
    return
  }
  c.Next()
})
```

#### Binding de JSON
```go
var produto Produto
if err := c.ShouldBindJSON(&produto); err != nil {
  c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
  return
}
```

#### Respostas HTTP
```go
c.JSON(http.StatusOK, gin.H{"mensagem": "Sucesso"})
c.JSON(http.StatusCreated, produto)
c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro interno"})
```

### 6.2 GORM - ORM (Object-Relational Mapping)
**URL:** gorm.io/gorm

**Utilidade:** Interação com banco de dados

**Recursos utilizados:**

#### Model definition
```go
type Produto struct {
  gorm.Model
  Codigo    string  `json:"codigo" gorm:"uniqueIndex"`
  Descricao string  `json:"descricao"`
  Saldo     float64 `json:"saldo"`
}
```

#### AutoMigration (Schema auto-create)
```go
config.DB.AutoMigrate(&model.Produto{})
config.DB.AutoMigrate(&model.NotaFiscal{}, &model.ItemNota{})
```

#### CRUD Operations
```go
// Create
config.DB.Create(&produto)

// Read
var produto Produto
config.DB.First(&produto, id)

// List
var produtos []Produto
config.DB.Find(&produtos)

// Update
config.DB.Model(&Produto{}).
  Where("id = ?", id).
  Update("saldo", novoSaldo)
```

#### Transações com Lock
```go
tx := config.DB.Begin()

// Row-level locking (pessimistic)
tx.Clauses(clause.Locking{Strength: "UPDATE"}).
  Order("numero desc").
  First(&ultima)

tx.Create(&nota)
tx.Commit()
```

#### Associações (Relationships)
```go
type NotaFiscal struct {
  Itens []ItemNota `gorm:"foreignKey:NotaFiscalID"`
}

// Preload com associações
config.DB.Preload("Itens").Find(&notas)
```

---

## 7. TRATAMENTO DE ERROS E EXCEÇÕES NO BACKEND

### 7.1 em Handlers HTTP

#### Validação de Input
```go
func CriarNota(c *gin.Context) {
  var n model.NotaFiscal
  if err := c.ShouldBindJSON(&n); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
    return
  }

  if len(n.Itens) == 0 {
    c.JSON(http.StatusBadRequest, 
      gin.H{"erro": "A nota precisa de pelo menos um item"})
    return
  }
}
```

#### Validação de Regras de Negócio
```go
if produto.Saldo < item.Quantidade {
  c.JSON(http.StatusBadRequest, 
    gin.H{"erro": "Saldo insuficiente"})
  return
}
```

#### Timeout em HTTP Client
```go
client := &http.Client{Timeout: 5 * time.Second}
resp, err := client.Get(url) // Timeout após 5s
if err != nil {
  c.JSON(http.StatusServiceUnavailable, 
    gin.H{"erro": "Serviço indisponível"})
  return
}
```

#### Tratamento de Network Errors
```go
if err != nil {
  fmt.Printf("[ERRO] Falha ao buscar produto: %v\n", err)
  c.JSON(http.StatusInternalServerError, 
    gin.H{"erro": "Erro de comunicação"})
  return
}
```

### 7.2 Rollback em Falhas

**Cenário:** Atualizar múltiplos saldos, um falha no meio

```go
successful := []Update{} // Rastreia sucesso

for _, update := range updates {
  resp, err := client.Do(request)
  if err != nil {
    // Desfazer as mudanças anteriores
    rollbackEstoque(estoqueURL, successful)
    
    c.JSON(http.StatusServiceUnavailable, 
      gin.H{"erro": "Falha ao atualizar estoque"})
    return
  }
  successful = append(successful, update)
}
```

### 7.3 Logging de Erros

```go
fmt.Printf("[ERRO] Rollback falhou para produto %d: %v\n", 
  update.ProdutoID, err)

fmt.Printf("[AVISO] Status %d retornado para produto %d\n", 
  resp.StatusCode, product ID)
```

---

## 8. ARQUITETURA DE MICROSSERVIÇOS

### 8.1 Separação de Responsabilidades

**Estoque Service:**
- Gerencia CRUD de produtos
- Controla saldos
- **Não conhece** notas fiscais

**Faturamento Service:**
- Gerencia CRUD de notas
- Coordena impressão
- **Chama** Estoque para atualizar saldos

### 8.2 Comunicação Inter-serviços

HTTP REST com proper error handling:

```go
// Faturamento.handler
resp, err := client.Get(fmt.Sprintf("%s/produtos/%d", estoqueURL, id))

if resp.StatusCode != http.StatusOK {
  return nil, fmt.Errorf("produto não encontrado")
}

json.NewDecoder(resp.Body).Decode(&produto)
```

### 8.3 Docker Networking

```yaml
# docker-compose.yml
services:
  estoque:
    environment:
      PORT: 8080
    
  faturamento:
    environment:
      ESTOQUE_URL: http://estoque:8080  # Usa nome do serviço
```

---

## 9. PADRÕES E BOAS PRÁTICAS

### 9.1 Repository Pattern (Go)
Abstrai acesso a dados em métodos específicos:

```go
// repository/produto_repository.go
func ListarProdutos() ([]model.Produto, error)
func CriarProduto(p *model.Produto) error
func BuscarProdutoPorID(id uint) (*model.Produto, error)
func AtualizarSaldo(id uint, saldo float64) error
```

**Vantagens:**
- Fácil mudar banco de dados depois
- Testes mais simples (pode mockar)
- Lógica centralizada

### 9.2 Handler Pattern (Go)
Cada endpoint = uma função handler:

```go
func ListarNotas(c *gin.Context) { /* ... */ }
func CriarNota(c *gin.Context) { /* ... */ }
func ImprimirNota(c *gin.Context) { /* ... */ }
```

### 9.3 Service Pattern (Angular)
Lógica de HTTP centralizada:

```typescript
// services/nota.service.ts
@Injectable({ providedIn: 'root' })
export class NotaService {
  listar(): Observable<NotaFiscal[]> { /* ... */ }
  criar(nota): Observable<NotaFiscal> { /* ... */ }
  imprimir(id): Observable<any> { /* ... */ }
}
```

### 9.4 Numeração Sequencial com Concorrência

**Problema:** Dois usuários criando notas simultâneos → números duplicados

**Solução:** Row-level locking pessimista (GORM + PostgreSQL)

```go
tx := config.DB.Begin()

// Busca ÚLTIMA nota com lock UPDATE
tx.Clauses(clause.Locking{Strength: "UPDATE"}).
  Order("numero desc").
  First(&ultima)

novaNumero := ultima.Numero + 1
nota.Numero = novaNumero

tx.Create(&nota)
tx.Commit()
```

**Como funciona:**
1. Lock adquirido pelo usuário A
2. Usuário B espera pelo lock
3. Usuário A incrementa número e libera
4. Usuário B pega lock, vê número maior

---

## 10. AMBIENTE DE DESENVOLVIMENTO

### Docker Compose Setup

```yaml
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_PASSWORD: 3699
    ports:
      - "5432:5432"

  estoque:
    build: ./estoque
    environment:
      DB_HOST: postgres
      ESTOQUE_URL: http://estoque:8080
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy

  faturamento:
    build: ./faturamento
    environment:
      ESTOQUE_URL: http://estoque:8080
    ports:
      - "8081:8081"
```

**Iniciar:**
```bash
docker-compose up --build
```

---

## 11. TESTES (Como Estruturar)

### Testes Angular (Jasmine/Karma)
```typescript
describe('ProdutoService', () => {
  it('deve listar produtos', (done) => {
    service.listar().subscribe(produtos => {
      expect(produtos.length).toBeGreaterThan(0);
      done();
    });
  });
});
```

### Testes Go
```go
func TestListarProdutos(t *testing.T) {
  produtos, err := repository.ListarProdutos()
  assert.NoError(t, err)
  assert.Greater(t, len(produtos), 0)
}
```

---

## 12. FLUXO COMPLETO DA APLICAÇÃO

### Use Case: Imprimir Nota Fiscal

**1. Usuario clica botão "Imprimir"**
   - Frontend: `imprimir(notaId)`
   - Estado: `imprimindo = notaId`

**2. Frontend faz POST para /notas/:id/imprimir**
   ```
   POST http://localhost:8081/notas/5/imprimir
   ```

**3. Backend (Faturamento) processa:**
   - Busca nota (verifica status = "Aberta")
   - Para cada item, busca produto no estoque
   - Valida saldo
   - Atualiza saldo (PATCH /produtos/:id/saldo)
   - Fecha nota (UPDATE status = "Fechada")

**4. Volta para Frontend:**
   - Status 200 OK
   - `imprimindo = null` (remove spinner)

**5. Frontend carrega dados novamente:**
   - Nova lista de notas (status agora "Fechada")
   - Nova lista de produtos (saldos atualizados)

**6. Usuário vê:**
   - Spinner desaparece
   - Nota com status "Fechada"
   - SnackBar de sucesso
   - Saldos dos produtos atualizados

---

## 13. RESUMO TÉCNICO PARA VÍDEO

**Pontos-chave a mencionar:**

1. **Ciclos de vida:** OnInit (dados iniciais), OnDestroy (cleanup)
2. **RxJS:** Reactive programming com pipe, catchError, finalize, takeUntil
3. **Material:** Design system para UI consistente
4. **GORM:** ORM que simplifica queries, suporta transações e locks
5. **Gin:** HTTP framework minimalista mas poderoso
6. **Microsserviços:**  Estoque e Faturamento independentes, comunicação HTTP
7. **Tratamento de erros:** Validações, timeouts, rollback automático
8. **Concorrência:** Row-level locking garante numeração única
9. **Docker:** Orquestração de serviços com docker-compose

