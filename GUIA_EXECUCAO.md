# GUIA DE EXECUÇÃO - Sistema Korp ERP Notas Fiscais

## Pré-requisitos

Você terá que ter instalado:

### Opção A: Apenas Docker (RECOMENDADO)
- Docker Desktop (Windows/Mac) ou Docker + Docker Compose (Linux)
- [Download Docker](https://www.docker.com/products/docker-desktop)

### Opção B: Desenvolvimento Local
- Go 1.21+ ([Download](https://golang.org/dl))
- Node.js 18+ + npm ([Download](https://nodejs.org))
- PostgreSQL 15+ ([Download](https://www.postgresql.org/download))

---

## Configuração com Docker (MAIS FÁCIL)

### 1. Clonar o Repositório
```bash
git clone https://github.com/seu-usuario/Korp_Teste_NomeAqui.git
cd Korp_Teste_NomeAqui
```

### 2. Iniciar com Docker Compose
```bash
docker-compose up --build
```

**Isso vai:**
- ✅ Criar container PostgreSQL
- ✅ Build e rodar serviço Estoque (porta 8080)
- ✅ Build e rodar serviço Faturamento (porta 8081)
- ✅ Criar databases (estoque, faturamento)

### 3. Em OUTRO terminal, rodar Frontend
```bash
cd frontend
npm install  # na primeira vez
npm start
```

### 4. Acessar Aplicação
- Frontend: **http://localhost:4200**
- API Estoque: http://localhost:8080/produtos
- API Faturamento: http://localhost:8081/notas

### 5. Parar Tudo
```bash
# Terminal 1
Ctrl+C  (docker-compose)

# Terminal 2
Ctrl+C  (npm)
```

---

## Configuração Local (Desenvolvimento)

### 1. Setup Banco de Dados

#### Windows (PostgreSQL)
```sql
-- Abra PgAdmin ou algum cliente
CREATE DATABASE estoque;
CREATE DATABASE faturamento;
```

#### macOS (com Homebrew)
```bash
brew install postgresql@15
brew services start postgresql@15

psql postgres
CREATE DATABASE estoque;
CREATE DATABASE faturamento;
\q
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get install postgresql postgresql-contrib
sudo -u postgres psql

CREATE DATABASE estoque;
CREATE DATABASE faturamento;
\q
```

### 2. Terminal 1 - Serviço Estoque
```bash
cd estoque

# Instalar dependências (primeira vez)
go mod download

# Executar
go run main.go
```

**Outputs esperado:**
```
[GIN-debug] engine.Run() listening on 0.0.0.0:8080
```

### 3. Terminal 2 - Serviço Faturamento
```bash
cd faturamento

# Instalar dependências (primeira vez)
go mod download

# Executar
go run main.go
```

**Output esperado:**
```
[GIN-debug] engine.Run() listening on 0.0.0.0:8081
```

### 4. Terminal 3 - Frontend
```bash
cd frontend

# Instalar dependências (primeira vez)
npm install

# Executar dev server
npm start
```

**Output esperado:**
```
✔ Compiled successfully
✔ 4 bundles generated
✔ Build at: 2026-04-12 10:30:00
```

### 5. Acessar
- Frontend: **http://localhost:4200**

Se tudo deu certo, você verá logo azul "Korp ERP" com menu de navegação.

---

## Testando a Aplicação

### Teste 1: Cadastrar Produtos

1. Clique em **"Produtos"** no menu
2. Preencha:
   - Código: **PROD-001**
   - Descrição: **Teclado**
   - Saldo: **50**
3. Clique **"Cadastrar Produto"**
4. Deve aparecer na tabela abaixo

### Teste 2: Criar Nota Fiscal

1. Clique em **"Notas Fiscais"**
2. No dropdown "Produto", selecione **PROD-001**
3. Quantidade: **5**
4. Clique **"Adicionar"**
5. Clique **"Criar Nota Fiscal"**
6. Deve aparecer NF #1 na tabela com status "Aberta"

### Teste 3: Imprimir Nota

1. Na tabela de notas, clique **"Imprimir"** em NF #1
2. Deve mostrar spinner/loading por alguns segundos
3. Status deve mudar para **"Fechada"**
4. Volta em Produtos e verifica: saldo de PROD-001 deve ser **45** (era 50, usou 5)

---

## Troubleshooting

### ❌ "Cannot connect to database"

**Solução:**
```bash
# Verifica se PostgreSQL está rodando
# Windows: Procure "Services" → PostgreSQL

# macOS:
brew services list

# Linux:
sudo systemctl status postgresql
```

### ❌ "port 8080/8081 already in use"

**Uma das portas já está ocupada.**

```bash
# Windows
netstat -ano | findstr :8080

# macOS/Linux
lsof -i :8080

# Matar processo (ex: PID 1234)
# Windows: taskkill /PID 1234 /F
# macOS/Linux: kill -9 1234
```

### ❌ "Frontend não conecta na API"

**Verifica environment.ts**

Arquivo: `frontend/src/environments/environment.ts`

Deve ter:
```typescript
estoqueApiUrl: 'http://localhost:8080',
faturamentoApiUrl: 'http://localhost:8081',
```

### ❌ "npm install falha com erro de permissão"

```bash
# Windows (rodar como Admin)
npm install -g npm  # atualizar npm

# macOS/Linux
sudo npm install
```

### ❌ "Componentes Material não aparecem"

```bash
cd frontend
npm install @angular/material
ng add @angular/material
```

### ❌ Erro ao rodar `npm start` (Angular compilation error)

```bash
# Limpar cache
rm -rf node_modules package-lock.json
npm install
npm start
```

---

## Verificar se Tudo Está Funcionando

### Teste de Conectividade

#### Teste Estoque API
```bash
# Windows (PowerShell)
Invoke-WebRequest -Uri "http://localhost:8080/produtos"

# macOS/Linux
curl http://localhost:8080/produtos
```

**Resposta esperada:**
```json
{"message":"ok","data":[]}
// ou lista de produtos se tiver criado alguns
```

#### Teste Faturamento API
```bash
curl http://localhost:8081/notas
```

### Teste Completo

1. ✅ Cadastre 3 produtos
2. ✅ Crie 2 notas com múltiplos produtos
3. ✅ Imprima 1 nota
4. ✅ Verifica saldos foram atualizados
5. ✅ Tenta imprimir nota já fechada (deve ter apenas "Concluída")

Se tudo passou, **está funcionando 100%!**

---

## Limpeza / Reset de Dados

### Com Docker
```bash
# Sematar tudo
docker-compose down

# Remover dados (volumes)
docker-compose down -v

# Reiniciar limpo
docker-compose up --build
```

### Local
```bash
# No psql
DROP DATABASE estoque;
DROP DATABASE faturamento;

CREATE DATABASE estoque;
CREATE DATABASE faturamento;

# Serviços Go vão recriar as tabelas automaticamente (AutoMigrate)
```

---

## Performance & Logs

### Ver Logs do Frontend
Abra DevTools (F12) → Console
```
Erros aparecem em vermelho
Info em azul
Warnings em amarelo
```

### Ver Logs do Backend (Go)
```
[GIN-debug] GET /produtos
[GIN] 2026/04/12 - 10:30:45 | 200 | 123.456µs | ::1 | GET /produtos
```

- **200** = Sucesso
- **400** = Erro de input
- **500** = Erro interno
- **503** = Serviço indisponível

### Slow Queries

Se notas demora muito para imprimir:
1. Verifica conexão banco de dados
2. Se Estoque API demora muito (5+ seg), pode ter timeout
3. Logs vão informar o problema

---

## Build para Produção

### Frontend
```bash
cd frontend
npm run build
# Gera pasta dist/
# Deploy a statico hosting (Netlify, Vercel, S3, etc)
```

### Backend
```bash
# Já rodando em Docker automaticamente
docker build -t korp-estoque ./estoque
docker push seu-usuario/korp-estoque:latest
```

---

## Próximos Passos (Após Submissão)

- [ ] Deploy em cloud (Heroku, AWS, Google Cloud)
- [ ] Setup CI/CD (GitHub Actions)
- [ ] Mais testes automatizados
- [ ] Monitoramento (Sentry, DataDog)
- [ ] Autoscaling de containers
- [ ] Load balancing entre instâncias

---

## Suporte Técnico

Se encontrar problemas, verifique:

1. ✅ Logs (DevTools frontend, terminal backend)
2. ✅ Arquivo `.env` está correto
3. ✅ Banco de dados está acessível
4. ✅ Portas não estão em conflito
5. ✅ Versões compatíveis (Go 1.21+, Node 18+, Angular 17)

Arquivo: [DETALHAMENTO_TECNICO.md](DETALHAMENTO_TECNICO.md) tem mais detalhes.

