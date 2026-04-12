# FAQ - Perguntas Frequentes & Respostas

Prepare-se com estas perguntas que podem ser feitas durante a apresentação ou entrevista.

---

## PERGUNTAS TÉCNICAS

### 1. Por que você escolheu Angular em vez de React/Vue?

**Resposta:**
"Angular foi escolhido pelos seguintes motivos:
- **Full framework** com tudo incluído (routing, HTTP client, forms, CLI)
- **Strongly typed** com TypeScript native, reduz bugs
- **Material Design** integrado e bem mantido
- **RxJS nativo** para reatividade (melhor que hooks do React)
- **Componentes standalone** (Angular 17+) são bem simples
- O requisito do teste mencionou Angular/C#, então Angular foi natural"

---

### 2. Explique a diferença entre OnInit e outras lifecycle hooks

**Resposta:**
"Angular tem 8 lifecycle hooks. Os principais:

- **OnInit**: Chamado UMA VEZ após o componente ser inicializado
  - Melhor lugar para carregar dados da API
  - Exemplo: `this.carregarProdutos()`

- **OnDestroy**: Chamado ANTES do componente ser destruído
  - Limpar subscriptions e liberar recursos
  - Previne memory leaks
  - Exemplo: `this.destroy$.next()`

- **OnChanges**: Detecta mudanças em @Input properties
  - Poderia usar se tivéssemos componentes pai-filho

- **Outros**: OnAfterViewInit, OnAfterContentInit, etc
  - Menos usados, mas importantes para casos específicos"

---

### 3. Como RxJS melhora a reatividade?

**Resposta:**
"RxJS oferece:

1. **Observables** - Streams de dados que emitem valores ao longo do tempo
   ```typescript
   this.productService.list() // Retorna Observable
   ```

2. **Operators** - Transformam dados no pipeline
   - `catchError()` - Trata erros
   - `map()` - Transforma dados
   - `filter()` - Filtra items
   - `finalize()` - Cleanup garantido
   - `takeUntil()` - Cancela subscription

3. **Lazy evaluation** - Código só executa quando subscribe
   ```typescript
   obs.pipe(op1, op2, op3).subscribe() // Agora executa
   ```

4. **Composition** - Fácil combinar lógica complexa
   ```typescript
   this.products$ = this.search$.pipe(
     debounceTime(300),
     distinctUntilChanged(),
     switchMap(term => this.search(term)),
     catchError(...)
   )
   ```

Isso torna o código mais declarativo e fácil de testar."

---

### 4. Por que usar row-level locking no banco?

**Resposta:**
"Problema: Dois usuários criam notas ao mesmo tempo
- Nota 1 vê últimaNumero = 0, incrementa para 1
- Nota 2 vê últimaNumero = 0, incrementa para 1
- Duas notas com número 1 (violação de constraint)

Solução: Row-level locking (pessimistic)
```go
tx.Clauses(clause.Locking{Strength: "UPDATE"}).
   Order("numero desc").
   First(&ultima)
```

Funciona assim:
1. Usuário A adquire LOCK na linha
2. Usuário B espera para não ler dado inconsistente
3. Usuário A incrementa e commit
4. Usuário B pega o lock, vê número maior
5. Cada nota fica com número único

Alternativa: Usar sequência no banco (SERIAL type)
Mas lock explícito deixa mais claro o intent."

---

### 5. Qual é o padrão de comunicação entre microsserviços?

**Resposta:**
"Estou usando **HTTP REST** simples:

Faturamento → Estoque
```
GET  /produtos/5         # Buscar para validar
PATCH /produtos/5/saldo  # Atualizar saldo
```

Vantagens:
- Simples de debug (curl, postman, console.log)
- Não precisa de broker message (Kafka, RabbitMQ)
- Bom para prototipagem

Desvantagens:
- Chamadas síncronas (se Estoque der erro, tudo falha)
- Sem retry automático
- Sem circuit breaker nativo

Melhorias futuras:
- Implementar retry com exponential backoff
- Circuit breaker pattern
- Message queue para operações async
- Service discovery (ao invés de hardcoded URLs)"

---

### 6. Como você trataria a idempotência?

**Resposta:**
"Idempotência = operações repetidas não causam efeitos duplicados

Meu sistema já tem idempotência por:
1. **Numeração única**: Se tenta criar mesma nota 2x, segunda tira erro de constraint
2. **Saldos com transação**: Atualizar saldo é atômico, não deixa meia-transação

Poderia melhorar com:
1. **Idempotency keys**: Cliente envia UUID único
   ```
   POST /notas
   Header: X-Idempotency-Key: abc-123
   ```
   Backend armazena key e resultado, retorna mesmo resultado se tenta 2x

2. **Versionamento**: Nota tem version field
   ```
   UPDATE nota SET ... WHERE id = ? AND version = ?
   ```
   Se outro thread mudou, versão não bate, tira erro

3. **Event sourcing**: Armazenar eventos, não estado
   - Tenta criar nota 2x
   - Primeira vez: evento criado, emitido
   - Segunda vez: detecta evento já existe, retorna resultado anterior"

---

### 7. Por quanto tempo você conhece Go/Angular?

**Resposta (adapte a seus dados):**
"Começei a aprender Angular há ~3 meses através do [source: curso Udemy/Github projects/etc].

Principais projetos:
- Projeto X com Angular + Material
- Projeto Y integrando APIs

Go tenho experiência de ~[X] meses/anos com:
- Building APIs com Gin
- Database access com GORM
- Microsserviços

Foco de aprendizado tem sido:
- RxJS operators e padrões reativos
- GORM relationships e transactions
- Docker para containerizar apps"

---

### 8. Como você debugaria um problema no banco?

**Resposta:**
"Passos que eu usaria:

1. **Verificar conexão**
   ```bash
   psql -h localhost -U postgres -d estoque
   ```

2. **Ver schema**
   ```sql
   \dt  # Listar tabelas
   \d produtos  # Ver estrutura
   ```

3. **Executar queries manualmente**
   ```sql
   SELECT * FROM produtos WHERE id = 5;
   SELECT * FROM notas ORDER BY numero DESC LIMIT 10;
   ```

4. **Ver logs do Go**
   - Adicionar logging em repository methods
   - Ver queries que estão sendo executadas

5. **Usar GORM query logging**
   ```go
   db := db.Session(&gorm.Session{DebugMode: true})
   ```

6. **Monitorar com pg_stat_statements**
   ```sql
   SELECT query, mean_exec_time FROM pg_stat_statements;
   ```

Ferramenta favorita: PgAdmin (GUI) ou psql (CLI"

---

## PERGUNTAS DE DESIGN

### 9. Como você estruturaria o projeto se tivesse que escalar para 100k usuários?

**Resposta:**
"Hoje: Monólito com banco compartilhado

Para escalar:

1. **Banco de dados**
   - Replicação Master-Slave
   - Sharding por cliente/região
   - Cache (Redis) para queries frequentes

2. **Microsserviços**
   - Estoque e Faturamento em servidores separados
   - Scale horizontalmente (múltiplas instâncias)
   - Load balancer na frente

3. **Frontend**
   - CDN global para assets
   - Progressive Web App
   - Cache inteligente

4. **Message Queue**
   - Kafka para notas de impressão (async)
   - Não esperar Estoque responder em tempo real

5. **Monitoring**
   - Prometheus para métricas
   - Grafana para dashboards
   - Sentry para erro tracking

6. **CI/CD**
   - GitHub Actions para deploy automático
   - Staging environment

Exemplo: Estoque em N servidores, Faturamento em M servidores, todos com DB replicada"

---

### 10. Qual foi o desafio técnico mais difícil?

**Resposta (adapte):**
"Maior desafio foi **garantir numeração sequencial nas notas fiscais sem duplicação**.

Inicialmente implementei sem lock:
- 2 requisições simultâneas → ambas veem número 0
- Ambas incrementam para 1
- CONSTRAINT viola

Descobri a solução através:
1. Pesquisa sobre race conditions
2. Leitura sobre GORM pessimistic locking
3. Aprendi sobre PostgreSQL row-level locks

Solução final: `clause.Locking{Strength: 'UPDATE'}`
- Trava a linha enquanto lê
- Segundo usuario espera
- Realmente funciona!

Lição aprendida: Concorrência é hard, sempre considere quando múltiplos users acessam dados"

---

## PERGUNTAS COMPORTAMENTAIS

### 11. Como você se mantém atualizado com tecnologias?

**Resposta:**
"Formas que uso:
1. **Documentação oficial** - Angular docs e Go docs
2. **YouTube** - Traversy Media, Fireship para resumos
3. **Comunidades** - Reddit r/golang, r/angular
4. **Projetos pessoais** - Aprendo fazendo
5. **Blogs técnicos** - Dev.to, Medium"

---

### 12. Como você abordaria um bug em produção?

**Resposta:**
"Processo:
1. **Reproduzir** - Consigo repetir o erro localmente?
2. **Logar** - Vejo o que tá acontecendo nos logs?
3. **Isolar** - É frontend, backend ou banco?
4. **Testar solução** - Faço em branch antes de prod
5. **Deploar** - Monitoro se não piora
6. **Documentar** - Registro pra evitar repetir"

---

### 13. Por que você quer trabalhar na Korp?

**Resposta (pesquise a empresa primeiro):**
"Interesses incluem:
- Stack moderno (Go, Angular, Postgres)
- Microsserviços - tema que estou aprofundando
- Cultura de [valor que Korp tem]
- Oportunidade de crescer em [área específica]"

---

## PERGUNTAS SOBRE O PROJETO

### 14. O que você adicionaria ao projeto se tivesse mais tempo?

**Resposta:**
"Features que gostaria de agregar:

**Curto prazo:**
- Testes automatizados (Jest + Angular Testing)
- Go unit tests com mocking
- Paginação na tabela de notas

**Médio prazo:**
- Autenticação JWT
- Roles/Permissões (admin, user, viewer)
- Relatórios em PDF
- Auditoria (quem fez o quê e quando)

**Longo prazo:**
- Dashboard com gráficos
- Notificações em tempo real (WebSocket)
- Mobile app (React Native)
- Analytics"

---

### 15. Como você testaria a aplicação?

**Resposta:**
"Cobertura de testes:

**Frontend (Angular)**
```
// Unit tests: componentes isolados
// Integration tests: component + service
// E2E tests: fluxo completo usuario
npm test  // Karma + Jasmine
```

**Backend (Go)**
```
// Unit tests: handlers sem DB
// Integration tests: handlers + DB mock
// End-to-end: real database
go test ./...
```

**Manual**
- QA testa scenarios:
  - Produto duplicado
  - Nota com items
  - Falha de rede
  - Concurrent operations"

---

## RESPOSTAS PRONTAS (Copy-Paste)

Use em apresentação se ficar branco:

**"Qual é o ponto forte do seu projeto?"**
→ "Tratamento robusto de erros com rollback automático + microsserviços bem separados"

**"Qual função do Go você mais gosto?"**
→ "Goroutines e channels - concorrência super simples comparado a outras linguagens"

**"Qual foi mais challengiendo - front ou back?"**
→ "Backend - garantir atomicidade e idempotência é complexo, row-level locking resolveu"

**"Quanto time você gastou?"**
→ "[X horas] - Maioria tempo em understanding requirements e arquitetura"

