# Korp Teste - Sistema de Emissão de Notas Fiscais

Sistema completo de emissão de notas fiscais com arquitetura de microsserviços, desenvolvido com **Go**, **Angular**, **PostgreSQL** e **Docker**. Implementa **Clean Architecture** em todas as camadas.

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend Angular 21+                      │
│              (Standalone Components + Material)              │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTP/REST
        ┌──────────────┴──────────────┐
        │                             │
┌───────▼──────────┐         ┌───────▼──────────┐
│  Microsserviço   │         │  Microsserviço   │
│     Estoque      │         │   Faturamento    │
│    (Go + Gin)    │         │    (Go + Gin)    │
│ Clean Arch       │         │ Clean Arch       │
└───────┬──────────┘         └───────┬──────────┘
        │                             │
        └──────────────┬──────────────┘
                       │
                ┌──────▼──────┐
                │  PostgreSQL  │
                │  (Único BD)  │
                └─────────────┘
```

## 📁 Estrutura do Projeto

### Backend - Clean Architecture

```
estoque/
├── domain/                    # Camada de Domínio (sem dependências)
│   ├── entities/             #  Entidades de negócio
│   ├── repositories/         # Interfaces abstratas
│   ├── usecases/            # Lógica de negócio
│   └── errors.go            # Erros de domínio
├── infrastructure/           # Camada de Infraestrutura
│   ├── persistence/         # Implementação de repositórios
│   ├── http/
│   │   ├── handlers/        # Handlers HTTP
│   │   └── dto/             # Data Transfer Objects
│   └── services/            # Integração com outros serviços
├── utils/                    # Utilitários (Logger, Errors)
├── config/                   # Configuração (BD, etc)
└── main.go                   # Entry point
```

### Frontend - Feature-Based Architecture

```
frontend/src/app/
├── core/                     # Funcionalidades principais
│   ├── api/                 # Serviços de API
│   ├── models/              # Interfaces TypeScript
│   └── interceptors/        # Interceptores HTTP
├── features/                # Módulos de negócio
│   ├── produtos/           # Feature de Produtos
│   └── notas/              # Feature de Notas Fiscais
├── shared/                  # Componentes compartilhados
└── app.config.ts           # Configuração global
```

## 🚀 Como Executar

### Pré-requisitos

- Docker e Docker Compose
- (Alternativo: Go 1.21+, Node.js 18+, PostgreSQL 15)

### Com Docker (Recomendado)

```bash
# 1. Clone o repositório
git clone https://github.com/seu-usuario/Korp_Teste_SeuNome.git
cd Korp_Teste_SeuNome

# 2. Inicie todos os serviços
docker-compose up -d

# 3. Aguarde os serviços ficarem saudáveis
docker-compose ps

# 4. Acesse a aplicação
# Frontend: http://localhost:4200
# API Estoque: http://localhost:8080/health
# API Faturamento: http://localhost:8081/health
```

### Localmente (sem Docker)

#### Backend - Estoque

```bash
cd estoque

# 1. Copie o arquivo de ambiente
cp .env.example .env

# 2. Instale dependências Go
go mod download

# 3. Execute o servidor
go run main.go
# Server escutando em http://localhost:8080
```

#### Backend - Faturamento

```bash
cd faturamento

# 1. Copie o arquivo de ambiente
cp .env.example .env

# 2. Instale dependências Go
go mod download

# 3. Execute o servidor
go run main.go
# Server escutando em http://localhost:8081
```

#### Frontend

```bash
cd frontend

# 1. Instale dependências
npm install

# 2. Inicie o servidor de desenvolvimento
npm start

#3. Abra http://localhost:4200 no navegador
```

## 📚 Tecnologias Utilizadas

### Backend

| Tecnologia | Versão | Finalidade |
|-----------|--------|-----------|
| Go | 1.26+ | Linguagem principal |
| Gin Framework | 1.12+ | Web framework |
| GORM | 1.31+ | ORM para banco de dados |
| PostgreSQL Driver | 1.6+ | Driver PostgreSQL |
| Zerolog | 1.31+ | Logging estruturado |

### Frontend

| Tecnologia | Versão | Finalidade |
|-----------|--------|-----------|
| Angular | 21+ | Framework principal |
| TypeScript | 5.9+ | Linguagem |
| Angular Material | 21+ | Componentes visuais |
| RxJS | 7.8+ | Programação reativa |
| Reactive Forms | 21+ | Validação de formulários |

### DevOps

| Tecnologia | Finalidade |
|-----------|-----------|
| Docker | Containerização |
| Docker Compose | Orquestração local |
| PostgreSQL | Banco de dados |

## 🎯 Funcionalidades Implementadas

### ✅ Obrigatórias

- [x] **Cadastro de Produtos**
  - Campos: Código, Descrição, Saldo
  - Validação de duplicação de código
  - Restrições: Saldo não negativo

- [x] **Cadastro de Notas Fiscais**
  - Numeração sequencial automática
  - Status: Aberta ou Fechada
  - Múltiplos produtos por nota
  - Validação de quantidade > 0

- [x] **Impressão de Notas Fiscais**
  - Atualização de status para Fechada
  - Atualização de saldo de produtos
  - Validação: Apenas notas Abertas podem ser impressas
  - Indicador de processamento

- [x] **Arquitetura de Microsserviços**
  - Serviço de Estoque (Produto Management)
  - Serviço de Faturamento (Nota Fiscal Management)
  - Comunicação via HTTP REST

- [x] **Tratamento de Falhas**
  - Retry automático (1 tentativa)
  - Health checks em ambos os serviços
  - Tratamento centralizado de erros

- [x] **Persistência em Banco Real**
  - PostgreSQL com GORM
  - Migrações automáticas
  - Transações nas operações críticas

### 🎁 Opcionais Implementados

- [x] **Tratamento de Concorrência**
  - Transações database com lock
  - Validação de saldos antes de atualização

- [x] **Observabilidade**
  - Logging estruturado com zerolog
  - Health checks
  - Tratamento de erro padronizado

## 🏛️ Clean Architecture

### Princípios Implementados

1. **Independência de Frameworks**
   - A lógica de negócio não depende de Gin, GORM ou Angular Material
   - Fácil de testar e reutilizar

2. **Independência de Interface de Usuário**
   - Frontend e backend são separados
   - Mudança na UI não afeta a lógica

3. **Independência de Banco de Dados**
   - Repositórios abstratos
   - Fácil trocar PostgreSQL por outro BD

4. **Independência de Agentes Externos**
   - Serviços abstratos
   - Fácil mockar em testes

### Camadas

#### Domain (Núcleo)
- **Entities**: `Produto`, `NotaFiscal`, `ItemNota`
- **Use Cases**: `CriarProduto`, `ListarProdutos`, `ImprimirNota`, etc
- **Interfaces**: `ProdutoRepository`, `EstoqueService`
- **Erros de Domínio**: Mensagens de erro de negócio

#### Application (Infra)
- **Repositories Impl**: `ProdutoRepositoryImpl`, `NotaFiscalRepositoryImpl`
- **HTTP Handlers**: `ProdutoHandler`, `NotaFiscalHandler`
- **DTOs**: `CriarProdutoRequest`, `NotaFiscalResponse`
- **Services**: `EstoqueServiceImpl`

#### Interface External (Apresentação)
- **HTTP Routers**: Configuração de rotas
- **Main**: Inicialização da aplicação
- **Config**: Configurações de ambiente

## 🔒 Segurança

- ✅ CORS restrito (apenas localhost:4200 em desenvolvimento)
- ✅ Validação de entrada em formulários (Frontend + Backend)
- ✅ Erros genéricos em produção (sem stacktrace)
- ✅ Não há credenciais hardcoded (via .env)
- ✅ HTTPS pronto para produção (CA certificates)

## 📊 Endpoints da API

### Estoque Service (Port 8080)

```bash
# Saúde
GET /health

# Produtos
GET    /produtos              # Listar todos
POST   /produtos              # Criar novo
GET    /produtos/:id          # Buscar por ID
PATCH  /produtos/:id/saldo    # Atualizar saldo
```

### Faturamento Service (Port 8081)

```bash
# Saúde
GET /health

# Notas Fiscais
GET    /notas                 # Listar todas
POST   /notas                 # Criar nova
GET    /notas/:id             # Buscar por ID
POST   /notas/:id/imprimir    # Imprimir nota
```

## 📝 Exemplos de Uso

### Criar Produto

```bash
curl -X POST http://localhost:8080/produtos \
  -H "Content-Type: application/json" \
  -d '{
    "codigo": "PROD001",
    "descricao": "Notebook Dell XPS",
    "saldo": 10
  }'
```

### Listar Produtos

```bash
curl http://localhost:8080/produtos
```

### Criar Nota Fiscal

```bash
curl -X POST http://localhost:8081/notas \
  -H "Content-Type: application/json" \
  -d '{
    "itens": [
      {"produto_id": 1, "quantidade": 2},
      {"produto_id": 3, "quantidade": 1}
    ]
  }'
```

### Imprimir Nota

```bash
curl -X POST http://localhost:8081/notas/1/imprimir
```

## 🧪 Testes

### Backend (Go)

```bash
# Executar testes
go test ./...

# Com cobertura
go test -cover ./...
```

### Frontend (Angular)

```bash
# Executar testes
npm test

# Build para produção
npm run build
```

## 📋 Ciclos de Vida Angular Utilizados

- ✅ `OnInit` - Carregamento inicial de dados
- ✅ `OnDestroy` - Cleanup e unsubscribe
- ✅ Reactive Forms - Gerenciamento de estado de formulário
- ✅ RxJS Operators - `takeUntil`, `tap`, `catchError`, `shareReplay`

## 🔗 Bibliotecas RxJS Utilizadas

- ✅ `Observable` - Fluxo de dados reativo
- ✅ `BehaviorSubject` - Estado gerenciável
- ✅ `takeUntil` - Prevenção de memory leak
- ✅ `tap` - Efeitos colaterais
- ✅ `catchError` - Tratamento de erro
- ✅ `shareReplay` - Cache de requisições
- ✅ `retryWhen` - Retry automático com delay

## 📦 Dependências do Frontend

- Angular Material - Componentes de UI profissional
- RxJS - Programação reativa
- TypeScript - Type safety

## ⚙️ Configuração de Ambiente

### Backend (.env)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=senha_segura
DB_NAME=estoque
PORT=8080
GIN_MODE=debug
ESTOQUE_URL=http://localhost:8080  # (apenas faturamento)
```

### Frontend

Configurado em `src/app/core/api/`:
- `ProdutoApiService` - Endpoint: `http://localhost:8080`
- `NotaApiService` - Endpoint: `http://localhost:8081`

## 🐛 Troubleshooting

### Serviço não conecta ao banco

```bash
# Verifique se PostgreSQL está rodando
docker-compose ps

# Reinicie os serviços
docker-compose restart
```

### Frontend não carrega dados

```bash
# Verifique se os backends estão saudáveis
curl http://localhost:8080/health
curl http://localhost:8081/health

# Verifique o console do navegador (F12)
# Pode estar bloqueado por CORS
```

### Porta já em uso

```bash
# Encontre o processo usando a porta
lsof -i :8080

# Libere a porta ou use outra em docker-compose.yml
```

## 📞 Suporte

Para dúvidas sobre a implementação, consulte:
- [CHECKLIST_PROFISSIONALISMO.md](./CHECKLIST_PROFISSIONALISMO.md)
- [IMPLEMENTACAO_PRATICA.md](./IMPLEMENTACAO_PRATICA.md)
- [TOP_10_PRIORIDADES.md](./TOP_10_PRIORIDADES.md)

## 📄 Licença

Projeto desenvolvido para fins educacionais como parte do processo seletivo da Korp ERP.

---

**Desenvolvido por**: [Seu Nome]  
**Data**: Abril 2026  
**Versão**: 1.0.0


### Frontend

```bash
cd frontend
npm install
npm start
```

A aplicação Angular será servida em `http://localhost:4200`.

## Configuração de ambiente

- `frontend/src/environments/environment.ts` define as URLs das APIs.
- `estoque/.env.example` e `faturamento/.env.example` aceitam `DB_DSN` ou variáveis separadas.

## Notas técnicas

- Angular: uso de `OnInit` nos componentes para carregar dados inicialmente.
- RxJS: `HttpClient` com `Observable`, `pipe` e `catchError` para tratamento de falhas.
- UI: `@angular/material` para tabelas, formulários, botões e notificações.
- Go: `gin-gonic/gin` para rotas HTTP e `gorm.io/gorm` para persistência PostgreSQL.
- Cada microsserviço tem seu próprio `go.mod` e dependências isoladas.

## Próximos incrementos recomendados

- Adicionar tratamento de concorrência para o saldo de produtos.
- Implementar idempotência na impressão de notas fiscais.
- Criar testes automatizados de integração entre os microsserviços.
- Melhorar a experiência de falha com retry e fallback no frontend.
