# Korp Teste - Sistema de Emissão de Notas Fiscais

Projeto backend/frontend para o desafio técnico de emissão de notas fiscais.

## Arquitetura

- `frontend/` - aplicação Angular 21 com componentes standalone e Angular Material.
- `estoque/` - microsserviço Go responsável pelo cadastro de produtos e controle de saldo.
- `faturamento/` - microsserviço Go responsável pela criação, listagem e impressão de notas fiscais.

## Funcionalidades implementadas

- Cadastro de produtos com código, descrição e saldo.
- Cadastro de notas fiscais com itens e quantidade.
- Impressão de nota fiscal que atualiza o status para `Fechada`.
- Atualização de saldo de produtos ao imprimir nota.
- Validação de itens na criação de nota.
- Uso de banco PostgreSQL via GORM.
- Processamento de erros no backend com respostas HTTP apropriadas.

## Como executar

### Pré-requisitos

- PostgreSQL
- Go 1.26+
- Node.js 18+ / npm 10+

### Backend

Para cada serviço, crie um arquivo `.env` a partir de `.env.example` e ajuste os valores.

#### Estoque

```bash
cd estoque
copy .env.example .env
# ajustar as variáveis de conexão
go run main.go
```

#### Faturamento

```bash
cd faturamento
copy .env.example .env
# ajustar as variáveis de conexão e ESTOQUE_URL
go run main.go
```

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
