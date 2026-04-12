# KORP ERP — Sistema de Gestão

Aplicação **ERP profissional** com arquitetura escalável seguindo **Clean Architecture** e **SOLID principles**.

## 🏗️ Arquitetura

### Backend (Go) — Clean Architecture
```
backend/
├── cmd/main.go                         # Entry point com graceful shutdown
├── internal/
│   ├── domain/                         # Entidades de negócio
│   │   ├── invoice.go                  # NotaFiscal, ItemNota
│   │   ├── inventory.go               # Produto
│   │   └── errors.go                  # Erros de domínio
│   ├── application/                    # Casos de uso (lógica de negócio)
│   │   ├── invoice_service.go         # Serviço de notas fiscais
│   │   ├── inventory_service.go       # Serviço de estoque
│   │   └── dto.go                     # Data Transfer Objects
│   ├── infrastructure/                 # Dependências externas
│   │   ├── repository/                # Repositórios em memória
│   │   ├── http/                      # Handlers e rotas HTTP
│   │   └── database/                  # Placeholder PostgreSQL
│   └── ports/                          # Interfaces (contratos)
│       ├── repository.go              # Contratos de repositório
│       └── service.go                 # Contratos de serviço
└── go.mod
```

### Frontend (Angular 21) — Feature-based Architecture
```
frontend/src/app/
├── core/                               # Serviços singleton
│   ├── services/                      # ApiService, AuthService
│   ├── guards/                        # AuthGuard
│   └── interceptors/                  # HTTP Error Interceptor
├── shared/                             # Componentes reutilizáveis
│   ├── components/                    # Header, Sidebar, ConfirmDialog
│   └── models/                        # Interfaces comuns
├── features/                           # Módulos de funcionalidade
│   ├── inventory/                     # Gestão de Estoque
│   │   ├── pages/                     # List + Detail
│   │   ├── services/                  # InventoryService
│   │   └── models/                    # Produto, CriarProdutoInput
│   └── invoice/                       # Faturamento
│       ├── pages/                     # List + Detail
│       ├── services/                  # InvoiceService
│       └── models/                    # NotaFiscal, ItemNota
├── app.config.ts                       # Configuração (zoneless CD)
├── app.routes.ts                       # Rotas com lazy loading
└── app.ts                             # Componente raiz
```

## 🚀 Como Executar

### Pré-requisitos
- **Go** 1.22+ ([download](https://go.dev/dl/))
- **Node.js** 20+ ([download](https://nodejs.org/))
- **Angular CLI** (`npm install -g @angular/cli`)

### 1. Backend
```bash
cd backend
go run ./cmd/main.go
```
O servidor inicia em `http://localhost:8080`.

### 2. Frontend
```bash
cd frontend
npm install
ng serve
```
A aplicação abre em `http://localhost:4200`.

## 📡 API Endpoints

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/health` | Health check |
| `GET` | `/api/produtos` | Listar produtos |
| `POST` | `/api/produtos` | Criar produto |
| `GET` | `/api/produtos/{id}` | Buscar produto por ID |
| `PATCH` | `/api/produtos/{id}/saldo` | Atualizar saldo |
| `GET` | `/api/notas` | Listar notas fiscais |
| `POST` | `/api/notas` | Criar nota fiscal |
| `GET` | `/api/notas/{id}` | Buscar nota por ID |
| `POST` | `/api/notas/{id}/imprimir` | Imprimir nota (deduz estoque) |

## 🧪 Testes

### Backend
```bash
cd backend
go test ./... -v
```

### Frontend
```bash
cd frontend
ng test
```

## 🎯 Princípios Aplicados

- **Clean Architecture** — Separação clara de responsabilidades
- **SOLID** — Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, Dependency Inversion
- **Dependency Injection** — Inversão de controle via interfaces
- **Feature-based Modules** — Organização por funcionalidade
- **Lazy Loading** — Carregamento sob demanda
- **Angular Signals** — Reatividade moderna (Angular 21+)
- **Thread Safety** — Repositórios com `sync.RWMutex`
- **CORS** — Middleware para acesso cross-origin

## 📦 Stack Tecnológico

| Camada | Tecnologia |
|--------|-----------|
| Backend | Go 1.24 (standard library) |
| Frontend | Angular 21 + Angular Material |
| Estilização | CSS custom + Material Design 3 |
| Tipografia | Inter + Roboto |
| Estado | Angular Signals |
| HTTP | HttpClient + interceptors |
