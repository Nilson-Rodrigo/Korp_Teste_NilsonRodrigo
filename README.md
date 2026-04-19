# Korp Teste - Sistema de Emissão de Notas Fiscais

Sistema de emissão de notas fiscais com arquitetura de microsserviços, desenvolvido com **Go**, **Angular**, **PostgreSQL** e **Docker**.

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

### Backend

```
estoque/
├── config/                   # Configuração do banco
├── utils/                    # Logger e erros
├── handler.go                # Handlers HTTP
├── models.go                 # Modelos GORM
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

## 🎁 Opcionais

- [ ] Tratamento de concorrencia
- [ ] Idempotencia
- [ ] Uso de IA

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

## 📄 Licenca

Projeto desenvolvido para fins educacionais como parte do processo seletivo da Korp ERP.

---

**Desenvolvido por**: [Seu Nome]
**Data**: Abril 2026
**Versao**: 1.0.0
