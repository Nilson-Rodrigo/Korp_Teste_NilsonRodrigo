# ROTEIRO DE APRESENTAÇÃO - Sistema de Notas Fiscais

## Template para Vídeo (Duração recomendada: 8-10 minutos)

---

## 1. INTRODUÇÃO (30 segundos)

### Fala:
"Olá, meu nome é [SEU NOME]. Vou apresentar o Sistema de Emissão de Notas Fiscais que desenvolvi seguindo os requisitos do teste. O projeto é uma aplicação completa com arquitetura de microsserviços utilizando Angular no frontend e Go no backend."

### Mostrar:
- [Captura de tela inicial do projeto]
- Logo/Branding do Korp ERP

---

## 2. VISÃO GERAL DA ARQUITETURA (1 minuto)

### Fala:
"A aplicação foi dividida em 3 camadas principais: o frontend em Angular com Material Design, dois microsserviços em Go - um de Estoque e outro de Faturamento - e um banco de dados PostgreSQL compartilhado."

### Mostrar:
- Diagrama da arquitetura (pode desenhar na tela ou mostrar README_COMPLETO.md)
- Explicar fluxo de comunicação entre serviços

```
┌─────────────────┐
│  Angular App    │
│   (PORT 4200)   │
└────┬────────┬───┘
     │        │
  8080      8081
     │        │
┌────▼──┐  ┌──▼──┐
│Estoque│  │Fatura
│(GO)   │  │(GO) │
└────┬──┘  └──┬──┘
     └────┬───┘
        PostgreSQL
```

---

## 3. FUNCIONALIDADE 1: CADASTRO DE PRODUTOS (2 minutos)

### Fala:
"Começando com o cadastro de produtos. Você preenche o código, descrição e saldo inicial em estoque. O código é único, então não permite duplicatas."

### Demonstração:

1. **Acessar tela de Produtos**
   - Navegar para aba "Produtos"
   - Mostrar formulário vazio

2. **Preencher formulário**
   ```
   Código: PROD-001
   Descrição: Teclado USB
   Saldo: 50
   ```

3. **Clicar "Cadastrar Produto"**
   - Mostrar SnackBar de sucesso: "Produto cadastrado com sucesso!"
   - Produto aparece na tabela abaixo

4. **Adicionar mais 2-3 produtos**
   ```
   PROD-002 | Mouse Logitech | 30
   PROD-003 | Monitor 24" | 15
   PROD-004 | Mousepad | 100
   ```

5. **Tentar adicionar código duplicado**
   - Mostrar erro do backend
   - SnackBar vermelho: "Erro ao cadastrar produto"

### Comentários Técnicos:
- "Validação no frontend e backend"
- "SnackBar do Material Design para feedback"
- "RxJS catchError para tratamento de erro"

---

## 4. FUNCIONALIDADE 2: CADASTRO DE NOTAS FISCAIS (2 minutos)

### Fala:
"Agora vamos criar notas fiscais. Cada nota tem uma numeração única e sequencial gerada automaticamente pelo sistema."

### Demonstração:

1. **Navegar para "Notas Fiscais"**
   - Mostrar formulário de criação

2. **Criar primeira nota**
   - Click no dropdown "Produto"
   - Selecionar "PROD-001 — Teclado USB (50 un.)"
   - Input "Quantidade": 5
   - Click "Adicionar"
   - Item aparece na lista com chip

3. **Adicionar segundo produto**
   - Selecionar "PROD-002 — Mouse Logitech"
   - Quantidade: 3
   - Click "Adicionar"

4. **Click "Criar Nota Fiscal"**
   - Mostrar loading do botão: "Criando..."
   - Nota aparece na tabela abaixo com:
     - Número (ex: NF #1)
     - Status: "Aberta"
     - 2 Produtos
   - SnackBar: "Nota fiscal criada com sucesso!"

5. **Criar segunda nota de demonstração**
   - Similar ao anterior
   - Vai ter Número NF #2
   - Adicionar alguns produtos diferentes

### Comentários Técnicos:
- "Numeração sequencial com row-level locking (pessimistic locking)"
- "Múltiplos itens na mesma nota"
- "Validação: impede produtos duplicados"
- "Prevenção de race conditions no banco"

---

## 5. FUNCIONALIDADE 3: IMPRESSÃO DE NOTAS FISCAIS (2-3 minutos)

### Fala (IMPORTANTE):
"A impressão de nota é a funcionalidade mais crítica. Quando clica em imprimir, o sistema valida o saldo, atualiza no estoque, e só depois marca como fechada."

### Demonstração:

1. **Verificar saldos atuais de produtos**
   - Volta para aba Produtos
   - Mostrar saldos:
     ```
     PROD-001: 50 un.
     PROD-002: 30 un.
     ```

2. **Voltar para Notas**
   - Selecionar nota NF #1 (que tem 5x PROD-001 e 3x PROD-002)
   - Click botão "Imprimir"

3. **Mostrar processamento**
   - Botão muda para "Processando..."
   - Spinner/loading aparece

4. **Após sucesso (5-10 seg)**
   - Status muda de "Aberta" para "Fechada"
   - Botão desaparece (substitui por "Concluída")
   - SnackBar verde: "Nota impressa! Status atualizado para Fechada."

5. **Verificar atualização de saldos**
   - Ir para aba Produtos
   - Mostrar que saldos foram atualizados:
     ```
     PROD-001: 45 un. (eram 50, usaram 5)
     PROD-002: 27 un. (eram 30, usaram 3)
     ```

6. **Tentar imprimir nota já fechada**
   - Voltar para Notas
   - Tentar clicar em NF #1 (que está fechada)
   - Não há botão "Imprimir", apenas "Concluída"
   - Mostrar que não permite

### Comentários Técnicos:

**Mencionar cada passo:**
1. "Validação: verifica se status é Aberta"
2. "Requisição HTTP para o serviço de Estoque"
3. "Validação: verifica se saldo é suficiente"
4. "Atualização de saldo em transação"
5. "Em caso de erro, rollback automático"
6. "Fechamento da nota apenas se tudo suceder"
7. "Indicador de processamento (RxJS finalize)"

---

## 6. TRATAMENTO DE FALHAS - Cenário Fictício (1 minuto)

### Fala:
"Agora vou demonstrar como sistema se comporta quando há falha. Se por algum motivo o serviço de estoque desconectar..."

### Demonstração:

1. **Simular com manualmente parar backend**
   - (Se em Docker) `docker stop korp_estoque`
   - Ou simplesmente criar cenário onde não vai responder

2. **Tentar criar nova nota**
   - Ir para Notas
   - Criar nota com produtos normalmente
   - Click "Criar Nota Fiscal"

3. **Erro no processamento**
   - Requisição tira timeout (5 segundos)
   - SnackBar VERMELHO: "Serviço de faturamento indisponível"
   - Botão volta a estado normal (salvando = false)
   - Nota NÃO é criada (rollback de transação)

4. **Reiniciar serviço**
   - (Se Docker) `docker start korp_estoque`
   - Tentar novamente
   - Agora funciona

### Comentários Técnicos:
- "HTTP Client com timeout de 5 segundos"
- "RxJS catchError captura o erro"
- "Transações garantem não deixar dados inconsistentes"
- "Rollback automático em caso de falha"

---

## 7. DETALHAMENTO TÉCNICO - FRONTEND (1 minuto)

### Fala:
"Do ponto de vista técnico, o frontend utiliza Angular 17 com componentes standalone..."

### Pode mostrar (ou só falar):
- Abrir VS Code
- Mostrar arquivo `notas.ts`
- Explicar:

```typescript
// Ciclo de vida: OnInit carrega dados
ngOnInit() {
  this.carregarNotas();
}

// OnDestroy para cleanup
ngOnDestroy() {
  this.destroy$.next();
}

// RxJS para reatividade
.pipe(
  catchError(...),  // Trata erro
  finalize(...),    // Cleanup de state
  takeUntil(...)    // Cancela subscription
)
```

**Mencionar:**
- ✅ Ciclos de vida: OnInit, OnDestroy
- ✅ RxJS: catchError, finalize, takeUntil
- ✅ Material Design para UI
- ✅ TypeScript para type safety

---

## 8. DETALHAMENTO TÉCNICO - BACKEND (1 minuto)

### Fala:
"O backend em Go divide responsabilidades em dois microsserviços..."

### Pode mostrar (ou só falar):
- Abrir VS Code / arquivos Go
- Mostrar estrutura:

```
estoque/
├── model/produto.go          (Entity)
├── repository/               (Data access)
├── handler/produto_handler   (HTTP routes)
└── main.go

faturamento/
├── model/nota_fiscal.go      (Entity)
├── repository/               (Data access)
├── handler/nota_handler      (HTTP + inter-service)
└── main.go
```

**Mencionar:**
- ✅ Gin-gonic para HTTP routing
- ✅ GORM como ORM
- ✅ PostgreSQL persistência física
- ✅ Row-level locking para concorrência
- ✅ Transações para atomicidade

---

## 9. ARQUITETURA & DECISÕES (30 seg)

### Fala:
"A escolha de microsserviços permite que cada serviço escale independentemente. Se o estoque ficar muito lento, podemos aumentar recursos só pra ele, sem afetar faturamento."

### Explicar:
- ✅ Separação de responsabilidades
- ✅ Independência de deploy
- ✅ Comunicação via HTTP REST
- ✅ Banco compartilhado (para simplicidade)

---

## 10. COMO EXECUTAR (30 seg)

### Mostrar:
```bash
# Opção 1: Com Docker Compose (fácil)
docker-compose up --build
# http://localhost:4200

# Opção 2: Local (desenvolvimanto)
# Terminal 1
cd estoque && go run main.go

# Terminal 2
cd faturamento && go run main.go

# Terminal 3
cd frontend && npm start
```

---

## 11. CONCLUSÃO (30 seg)

### Fala:
"Resumindo, o sistema implementa todos os requisitos: cadastro de produtos, notas com numeração sequencial, impressão com atualização de saldos, microsserviços, tratamento de erros, e banco de dados real. Além disso, usa boas práticas como RxJS no frontend, transações no backend, e row-level locking para concorrência."

### Chamar atenção para:
- ✅ Funcionalidades completas
- ✅ Tecnologias solicitadas (Angular, Go)
- ✅ Qualidade de código
- ✅ Tratamento robusto de erros
- ✅ Pronto para produção (com Docker)

---

## DICAS PARA O VÍDEO

### Antes de Gravar:
- [ ] Testar tudo funcionando (backend + frontend)
- [ ] Ter dados já criados para não perder tempo
- [ ] Usar zoom/font grande para legibilidade
- [ ] Preparar fala / ter script próximo
- [ ] Fazer um teste de áudio e vídeo

### Durante:
- [ ] Fale claramente e pausadamente
- [ ] Aponte com mouse/seta o que está falando
- [ ] Não fale tão rápido
- [ ] Mostre a tela inteira (interface + código)
- [ ] Minimize distrações (notificações, abas extras)

### Edição:
- [ ] Cortar partes mortas (load time, erros)
- [ ] Adicionar legenda em pontos-chave
- [ ] Aumentar taxa de reprodução se tomar muito tempo
- [ ] Adicionar slide de title no início
- [ ] Adicionar slide com links no final

### Tempo:
- Objetivo: 8-10 minutos
- Máximo: 15 minutos
- Mínimo: 5 minutos (muito rápido)

---

## LINKS PARA INCLUIR NO FINAL DO VÍDEO

```
Repositório GitHub: https://github.com/seu-usuario/Korp_Teste_SEuNome
README Completo: [Link repo]/README_COMPLETO.md
Detalhamento: [Link repo]/DETALHAMENTO_TECNICO.md
```

---

## EXEMPLO DE FALA NATURAL

```
"Oi galera! Então, eu desenvolvi um sistema de notas fiscais pra Korp.

A arquitetura é dividida em três partes: frontend em Angular, backend em Go 
com dois microsserviços, e PostgreSQL. 

Vou começar mostrando como cadastrar produtos... [faz] ... depois criar 
notas... [faz] ... e agora a parte critical que é imprimir a nota que 
atualiza os saldos direto no estoque em tempo real.

Do ponto de vista technical, tá usando Angular Material pra interface, 
RxJS pra reatividade. No backend tá com Gin e GORM, implementei row-level 
locking pra garantir numeração única mesmo com requests simultâneos, e 
transações pra garantir atomicidade.

Em caso de erro, por exemplo se o serviço de estoque desconectar, o sistema 
faz um rollback automático pra não deixar dados inconsistentes.

Pra rodar é bem simples, só docker-compose up e tudo sobe, ou pode rodar 
local se tiver Go e PostgreSQL instalados.

É isso aí!"
```

