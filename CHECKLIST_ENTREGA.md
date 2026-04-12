# CHECKLIST DE ENTREGA - Sistema Korp Notas Fiscais

Use este checklist para garantir que tudo está pronto antes de enviar por email.

---

## ✅ VERIFICAÇÃO DE CÓDIGO

### Frontend (Angular)

- [ ] Todos os componentes importados corretamente
- [ ] Não há erros de TypeScript/tslint
- [ ] Testes passando: `npm test`
- [ ] Build para produção: `npm run build` (sem erros)
- [ ] Material Design temas aplicados
- [ ] RxJS operators sendo usados (catchError, finalize, takeUntil)
- [ ] OnInit e OnDestroy implementados

### Backend - Estoque (Go)

- [ ] Sem erros compilação: `go build`
- [ ] Rotas HTTP: GET, POST, PATCH funcionando
- [ ] Validações de input presentes
- [ ] Tratamento de erro implementado
- [ ] Logs adicionados

### Backend - Faturamento (Go)

- [ ] Sem erros compilação: `go build`
- [ ] Rotas HTTP funcionando
- [ ] Inter-service communication (chama Estoque API)
- [ ] Rollback em caso de falha implementado
- [ ] Transaction com locking implementado

### Banco de Dados

- [ ] PostgreSQL rodando localmente ou em Docker
- [ ] Databases criadas: estoque, faturamento
- [ ] AutoMigrate cria tabelas automaticamente
- [ ] Dados persistem após restart

---

## ✅ FUNCIONALIDADES IMPLEMENTADAS

### Requisitos Obrigatórios

- [ ] **Cadastro de Produtos**
  - [ ] Código (único)
  - [ ] Descrição
  - [ ] Saldo (quantidade em estoque)
  - [ ] Validação de duplicação

- [ ] **Cadastro de Notas Fiscais**
  - [ ] Numeração sequencial
  - [ ] Status (Aberta/Fechada)
  - [ ] Múltiplos produtos
  - [ ] Quantidades

- [ ] **Impressão de Notas**
  - [ ] Botão visível e intuitivo
  - [ ] Indicador de processamento
  - [ ] Atualiza status para Fechada
  - [ ] Bloqueia impressão se não "Aberta"
  - [ ] Atualiza saldos corretamente
  - [ ] Rollback automático em falha

- [ ] **Microsserviços**
  - [ ] Serviço Estoque (porta 8080)
  - [ ] Serviço Faturamento (porta 8081)
  - [ ] Comunicação HTTP entre eles
  - [ ] Independência de deploy

- [ ] **Banco de Dados Real**
  - [ ] PostgreSQL
  - [ ] Persistência física
  - [ ] Transações

- [ ] **Tratamento de Falhas**
  - [ ] Serviço indisponível (timeout)
  - [ ] Erro de validação
  - [ ] Rollback automático
  - [ ] Feedback ao usuário

### Requisitos Opcionais

- [ ] **Concorrência** - Row-level locking implementado
- [ ] **IA** - (Opcional, não implementado é OK)
- [ ] **Idempotência** - Numeração única garante isso

---

## ✅ QUALIDADE DE CÓDIGO

### Frontend

- [ ] Sem console.log desnecessários
- [ ] Sem commented code
- [ ] Imports organizados
- [ ] Variáveis nomeadas claramente
- [ ] Componentes menores e reusáveis
- [ ] Estilos SCSS bem organizados
- [ ] Sem CSS inline (exceto inline required)
- [ ] Accessibility considerada (alt text, labels, etc)

### Backend

- [ ] Sem commented code
- [ ] Error handling em todos endpoints
- [ ] Validação de input
- [ ] Nomes de funções em inglês (standars Go)
- [ ] Estrutura de pasta clara
- [ ] Interface bem definida entre serviços

---

## ✅ DOCUMENTAÇÃO

- [ ] **README_COMPLETO.md**
  - [ ] Stack tecnológico listado
  - [ ] Arquitetura diagramada
  - [ ] Como rodar (Docker e local)
  - [ ] Estrutura de pastas

- [ ] **DETALHAMENTO_TECNICO.md**
  - [ ] Ciclos de vida mencionados
  - [ ] RxJS operators explicados
  - [ ] Bibliotecas e suas funções
  - [ ] Tratamento robusto de erros
  - [ ] Go frameworks explicados
  - [ ] Microsserviços arquitetura

- [ ] **ROTEIRO_VIDEO.md**
  - [ ] Script de apresentação
  - [ ] Fluxo de demonstração
  - [ ] Pontos técnicos a mencionar

- [ ] **GUIA_EXECUCAO.md**
  - [ ] Pré-requisitos
  - [ ] Setup com Docker
  - [ ] Setup local
  - [ ] Troubleshooting

---

## ✅ TESTES E EXECUÇÃO

### Testes Locais

- [ ] Cadastrar 5+ produtos
- [ ] Criar 3+ notas com múltiplos itens
- [ ] Imprimir 2+ notas e verificar saldos
- [ ] Tentar operações inválidas (código duplicado, etc)
- [ ] Verificar feedback visual (snackbars, spinners)
- [ ] Simular falta de conexão (parar backend)

### Build

- [ ] `npm run build` sem erros (frontend)
- [ ] `go build ./...` sem erros (backend)
- [ ] `docker-compose up --build` sem erros

---

## ✅ CONFIGURAÇÕES

### .env Files

- [ ] `estoque/.env` exist e configurado
- [ ] `faturamento/.env` existe e configurado
- [ ] `.env` NÃO está em .gitignore (deve ser versionado para fácil setup)
- [ ] Credentials são genéricas (postgres/3699) para teste

### Docker

- [ ] `docker-compose.yml` existe
- [ ] `Dockerfile` em estoque/
- [ ] `Dockerfile` em faturamento/
- [ ] `init-db.sql` existe

### Environment Angular

- [ ] `environment.ts` (dev) - localhost:8080/8081
- [ ] `environment.prod.ts` - URLs de produção

---

## ✅ ENTREGA

### GitHub Repository

- [ ] Nome: **Korp_Teste_SeuNome**
- [ ] Público (público, não privado!)
- [ ] Todos os arquivos inclusos
- [ ] .env versionado
- [ ] README.md no root (pode ser simples link para README_COMPLETO.md)
- [ ] Commits com mensagens significativas

### Vídeo de Apresentação

- [ ] Duração: 8-10 minutos (máximo 15)
- [ ] Qualidade áudio e vídeo ok
- [ ] Script seguido (não muito frenético)
- [ ] Mostra todas as 3 funcionalidades
- [ ] Menciona tecnologias utilizadas
- [ ] Detalhamento técnico conforme requisito

### Email de Submissão

Enviar para: **julia.canever@korp.com.br**

Assunto: `Teste Korp - Desenvolvimento Notas Fiscais`

Corpo do email:

```
Olá,

Segue em anexo/ links a submissão do teste técnico de desenvolvimento.

**Link Repositório GitHub:**
https://github.com/seu-usuario/Korp_Teste_SeuNome

**Link Vídeo de Apresentação:**
https://drive.google.com/file/d/... (ou One Drive, etc)

**Resumo das Funcionalidades:**
- ✅ Cadastro de Produtos
- ✅ Cadastro de Notas Fiscais com numeração sequencial
- ✅ Impressão de Notas com atualização de saldos
- ✅ Microsserviços (Estoque + Faturamento)
- ✅ Banco de dados PostgreSQL
- ✅ Tratamento de falhas e rollback automático

**Stack Técnico:**
- Frontend: Angular 17 + RxJS + Material Design
- Backend: Go 1.21 + Gin + GORM
- Database: PostgreSQL
- DevOps: Docker + Docker Compose

Obrigado!

[Seu Nome]
```

---

## ✅ ÚLTIMAS VERIFICAÇÕES (24h antes)

- [ ] Clonou o repositório em pasta nova e testou (simula receptor)
- [ ] Docker compose inteiro roda sem erros
- [ ] Frontend acessa em http://localhost:4200
- [ ] Consegue criar produto → nota → imprimir
- [ ] Saldos atualizam corretamente
- [ ] Todos 3 readmes estão bem escritos
- [ ] Vídeo foi gravado com qualidade
- [ ] Email está pronto pra enviar

---

## 🚀 PRONTO PARA ENVIAR!

Se todos os checkboxes estão marcados, você está pronto!

**Última checagem antes de enviar:**

1. Copia link do repositório
2. Cola no navegador - existe e é público?
3. Verifica README (está legível?)
4. Verifica código (sem TO DOs, comentários?)
5. Vídeo tem som de qualidade?
6. Email escrito corretamente?

**Sucesso na submissão! 🎉**

