# Sistema de Emissão de Notas Fiscais - Korp ERP

## Visão Geral

Sistema completo de emissão de notas fiscais desenvolvido com arquitetura de microsserviços, seguindo boas práticas de desenvolvimento, tratamento robusto de erros e gerenciamento de dependências.

### Stack Tecnológico

**Frontend:**
- **Angular 17** (Standalone Components, TypeScript)
- **RxJS** (Reactive Programming)
- **Angular Material** (Design System)
- **SCSS** (Estilização)

**Backend:**
- **Go 1.21** (Gin Web Framework)
- **GORM** (ORM)
- **PostgreSQL** (Banco de Dados)
- **Docker & Docker Compose**

---

## Arquitetura

### Microsserviços

```
┌─────────────────────────────────────────────────────────┐
│                   Frontend Angular                       │
│              (Standalone Components)                     │
└────────┬──────────────────────────┬─────────────────────┘
         │                          │
    PORT 8080              PORT 8081
         │                          │
    ┌────▼────────┐        ┌────────▼───┐
    │   ESTOQUE   │        │FATURAMENTO │
    │ Microsserviço        │Microsserviço  
    │   (GO)      │        │   (GO)     │
    └────┬────────┘        └────────┬───┘
         │                          │
         └───────────┬──────────────┘
                     │
            ┌────────▼────────┐
            │  PostgreSQL DB  │
            │  (Containers)   │
            └─────────────────┘
```

### Serviço de Estoque (Porto 8080)

Responsável pelo gerenciamento de produtos e saldos:

- **POST /produtos** - Criar novo produto
- **GET /produtos** - Listar todos os produtos
- **GET /produtos/:id** - Buscar produto por ID
- **PATCH /produtos/:id/saldo** - Atualizar saldo de um produto

### Serviço de Faturamento (Porto 8081)

Responsável pela gestão de notas fiscais e comunicação com estoque:

- **POST /notas** - Criar nova nota fiscal
- **GET /notas** - Listar todas as notas
- **GET /notas/:id** - Buscar nota por ID
- **POST /notas/:id/imprimir** - Imprimir/Fechar nota (atualiza saldos)

---

## Funcionalidades Implementadas

### 1. Cadastro de Produtos ✅

**Campos Obrigatórios:**
- Código (único)
- Descrição 
- Saldo em estoque

**Recursos:**
- Validação de entrada
- Verificação de duplicação de código
- Feedback visual com MarialSnackBar

### 2. Cadastro de Notas Fiscais ✅

**Campos Obrigatórios:**
- Numeração sequencial (auto-gerada)
- Status (Aberta/Fechada)
- Múltiplos produtos com quantidades

**Recursos:**
- Numeração sequencial com lock pessimista no banco de dados
- Validação de quantidade
- Prevenção de duplicação de produtos na mesma nota
- Interface intuitiva com adição/remoção de itens

### 3. Impressão de Notas Fiscais ✅

**Funcionalidades:**
- Botão intuitivo com indicador de processamento (spinner)
- Validação: apenas notas "Aberta" podem ser impressas
- Atualização de status para "Fechada"
- **Atualização de saldos** com comunicação inter-serviços
- Rollback automático em caso de falha

**Fluxo de Impressão:**
1. Validar nota está "Aberta"
2. Chamar estoque para buscar saldos atuais
3. Validar se há saldo suficiente para todos itens
4. Atualizar saldos no estoque
5. Fechar nota localmente
6. Exibir feedback ao usuário

---

## Ciclos de Vida do Angular Utilizados

### OnInit
- `Produtos.ts`: Carrega lista de produtos ao montar
- `Notas.ts`: Carrega notas e produtos ao inicializar

```typescript
ngOnInit() {
  this.carregarNotas();
  this.carregarProdutos();
}
```

### OnDestroy
- Limpeza de subscriptions com RxJS
- Prevenção de memory leaks

```typescript
ngOnDestroy(): void {
  this.destroy$.next();
  this.destroy$.complete();
}
```

---

## Uso de RxJS

### Operadores Utilizados:

1. **catchError** - Tratamento de erros HTTP
2. **finalize** - Cleanup de estados (ex: remover loading states)
3. **takeUntil** - Cancelamento de subscriptions no ngOnDestroy
4. **pipe** - Composição de operadores

### Padrão de Subscriptions:

```typescript
this.notaService.criar({ itens: this.itens })
  .pipe(
    catchError((err) => {
      this.handleError(err);
      return of(null);
    }),
    finalize(() => {
      this.salvando = false;
    }),
    takeUntil(this.destroy$)
  )
  .subscribe((res) => {
    if (res) this.onSuccess();
  });
```

---

## Bibliotecas Utilizadas

### Frontend

| Biblioteca | Versão | Utilidade |
|-----------|--------|----------|
| @angular/core | 17 | Framework |
| @angular/material | 17 | Components UI |
| @angular/common | 17 | Pipes e Diretivas |
| rxjs | 7+ | Reactive Programming |
| typescript | 5+ | Type Safety |

### Backend (Go)

| Biblioteca | Utilidade |
|-----------|----------|
| github.com/gin-gonic/gin | Web Framework HTTP |
| gorm.io | ORM (Object-Relational Mapping) |
| gorm.io/driver/postgres | Driver PostgreSQL |
| github.com/joho/godotenv | Configuração via .env |

---

## Gerenciamento de Dependências

### Go (go.mod)

```bash
go mod download    # Baixar dependências
go mod tidy        # Limpar não utilizadas
go build          # Build final
```

### npm (Frontend)

```bash
npm install       # Instalar dependências
npm start        # Executar dev server
npm test         # Rodar testes
npm run build    # Build produção
```

---

## Tratamento de Erros e Exceções

### Backend (Go)

1. **Validação de Input**
   ```go
   if len(n.Itens) == 0 {
     c.JSON(http.StatusBadRequest, gin.H{"erro": "..."})
   }
   ```

2. **Timeout em Requisições HTTP**
   ```go
   client := &http.Client{Timeout: 5 * time.Second}
   ```

3. **Logging de Falhas**
   ```go
   fmt.Printf("[ERRO] Rollback falhou: %v\n", err)
   ```

4. **Rollback em Transações**
   ```go
   tx := config.DB.Begin()
   defer func() {
     if r := recover(); r != nil {
       tx.Rollback()
     }
   }()
   ```

### Frontend (Angular)

1. **Tratamento com RxJS catchError**
   ```typescript
   .pipe(
     catchError((err) => {
       const msg = err.error?.erro || 'Erro padrão';
       this.snackBar.open(msg, 'Fechar');
       return of([]);
     })
   )
   ```

2. **Feedback Visual ao Usuário**
   - SnackBar para notificações
   - Loading spinners durante operações
   - Desabilitar botões durante processamento

3. **Validações Preventivas**
   - Verificar saldo antes de desabilitar botão
   - Validar quantidade > 0
   - Preventir duplicação de produtos

---

## Como Executar

### Opção 1: Com Docker Compose (Recomendado)

```bash
# Clonar repositório
git clone <repo-url>
cd Korp_Teste_NilsonRodrigo

# Iniciar serviços
docker-compose up --build

# Acessar em: http://localhost:4200
```

### Opção 2: Local (sem containers)

**Pré-requisitos:**
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+

**Setup Banco de Dados:**
```sql
CREATE DATABASE estoque;
CREATE DATABASE faturamento;
```

**Executar Serviços:**

```bash
# Terminal 1 - Estoque
cd estoque
go run main.go
# Rodando em http://localhost:8080

# Terminal 2 - Faturamento
cd faturamento
go run main.go
# Rodando em http://localhost:8081

# Terminal 3 - Frontend
cd frontend
npm install
npm start
# Rodando em http://localhost:4200
```

---

## Arquivos de Configuração

### .env (Estoque)
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=3699
DB_NAME=estoque
PORT=8080
```

### .env (Faturamento)
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=3699
DB_NAME=faturamento
ESTOQUE_URL=http://localhost:8080
PORT=8081
```

---

## Tratamento de Falhas

### Cenário: Serviço de Estoque Indisponível

**No Frontend:**
1. Requisição falha com timeout
2. catchError captura o erro
3. SnackBar exibe mensagem amigável
4. Interface retorna com dados vazios

**No Backend (Faturamento):**
1. PATCH para atualizar saldo falha
2. Função `rollbackEstoque()` reverte mudanças anteriores
3. Resposta HTTP 503 Service Unavailable
4. Logs informam qual produto causou falha

```go
for _, update := range updates {
  resp, err := client.Do(reqPatch)
  if err != nil {
    rollbackEstoque(estoqueURL, successful)
    c.JSON(http.StatusServiceUnavailable, ...)
    return
  }
}
```

---

## Recursos Opcionais Implementáveis

### ✅ Tratamento de Concorrência
Implementado com **row-level locking** no GORM:
```go
tx.Clauses(clause.Locking{Strength: "UPDATE"})
  .Order("numero desc")
  .First(&ultima)
```

### ✅ Idempotência
Numeração única de notas garante que mesma requisição não cria duplicatas.

### Inteligência Artificial (Futuro)
- Previsão de demanda por produto
- Recomendação de saldo mínimo
- Detecção de fraudes

---

## Testes

### Testes Angular (Karma/Jasmine)

```bash
npm test
```

### Testes Go

```bash
cd estoque && go test ./...
cd faturamento && go test ./...
```

---

## Estrutura de Pastas

```
Korp_Teste_NilsonRodrigo/
├── estoque/                    (Serviço Estoque)
│   ├── main.go                 
│   ├── go.mod
│   ├── config/database.go      
│   ├── handler/produto_handler.go
│   ├── model/produto.go
│   ├── repository/produto_repository.go
│   └── Dockerfile
│
├── faturamento/                (Serviço Faturamento)
│   ├── main.go                 
│   ├── go.mod
│   ├── config/database.go      
│   ├── handler/nota_handler.go 
│   ├── model/nota_fiscal.go
│   ├── repository/nota_repository.go
│   └── Dockerfile
│
├── frontend/                   (Aplicação Angular)
│   ├── src/
│   │   ├── app/
│   │   │   ├── app.ts
│   │   │   ├── app.routes.ts
│   │   │   ├── services/
│   │   │   │   ├── produto.ts
│   │   │   │   └── nota.ts
│   │   │   └── components/
│   │   │       ├── produtos/
│   │   │       └── notas/
│   │   ├── environments/
│   │   │   ├── environment.ts      (dev)
│   │   │   ├── environment.dev.ts  (dev)
│   │   │   └── environment.prod.ts (prod)
│   │   └── index.html
│   ├── package.json
│   ├── angular.json
│   └── tsconfig.json
│
├── docker-compose.yml
├── init-db.sql
└── README.md
```

---

## Próximos Passos / Melhorias

- [ ] Adicionar autenticação JWT
- [ ] Implementar histórico de alterações (audit log)
- [ ] Dashboard com gráficos de estoque
- [ ] Relatórios em PDF
- [ ] Notificações em tempo real (WebSocket)
- [ ] Paginação na tabela de notas
- [ ] Filtros avançados
- [ ] Testes E2E com Cypress
- [ ] CI/CD com GitHub Actions

---

## Contato & Informações

**Desenvolvido para:** Teste de Seleção Korp
**Data:** 2026
**Tecnologias:** Angular + Go + PostgreSQL
