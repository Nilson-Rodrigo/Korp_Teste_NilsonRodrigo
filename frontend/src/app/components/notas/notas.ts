import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatIconModule } from '@angular/material/icon';
import { catchError, finalize, Subject, takeUntil } from 'rxjs';
import { of } from 'rxjs';
import { NotaService, NotaFiscal, ItemNota } from '../../services/nota';
import { ProdutoService, Produto } from '../../services/produto';
import { StateService } from '../../services/state.service';

@Component({
  selector: 'app-notas',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatTableModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatCardModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatIconModule,
  ],
  templateUrl: './notas.html',
  styleUrl: './notas.css',
})
export class Notas implements OnInit, OnDestroy {
  // Dados
  notas: NotaFiscal[] = [];
  produtos: Produto[] = [];
  displayedColumns = ['numero', 'status', 'itens', 'acoes'];

  // Estado de formulário
  novoItem: ItemNota = { produto_id: 0, quantidade: 1 };
  itens: ItemNota[] = [];

  // Validação
  errosFormulario = { item: '', quantidade: '' };
  formularioValido = true;

  // Estado de UI
  imprimindo: number | null = null;
  salvando = false;
  carregando = true;

  // Cleanup subscriptions
  private destroy$ = new Subject<void>();

  constructor(
    private notaService: NotaService,
    private produtoService: ProdutoService,
    private snackBar: MatSnackBar,
    private stateService: StateService,
  ) {}

  // Ciclo de vida: OnInit
  ngOnInit() {
    console.log('Iniciando carregamento de notas e produtos...');
    this.carregarProdutos();
    this.carregarNotas();
  }

  carregarProdutos() {
    console.log('Carregando produtos...');
    this.produtoService.listar()
      .subscribe({
        next: (data) => {
          console.log('✓ Produtos carregados:', data);
          this.produtos = data;
        },
        error: (err) => {
          console.error('✗ Erro ao carregar produtos:', err);
          this.snackBar.open('Erro ao carregar produtos', 'Fechar', {
            duration: 5000,
            panelClass: 'snack-error',
          });
        }
      });
  }

  carregarNotas() {
    console.log('Carregando notas...');
    this.carregando = true;

    this.notaService.listar()
      .subscribe({
        next: (data) => {
          console.log('✓ Notas carregadas:', data);
          this.notas = data;
          this.carregando = false;
        },
        error: (err) => {
          console.error('✗ Erro ao carregar notas:', err);
          this.carregando = false;
          this.snackBar.open('Erro ao carregar notas', 'Fechar', {
            duration: 5000,
            panelClass: 'snack-error',
          });
        }
      });
  }

  // Validação em tempo real
  validarItem() {
    this.errosFormulario = { item: '', quantidade: '' };
    this.formularioValido = true;

    if (!this.novoItem.produto_id || this.novoItem.produto_id === 0) {
      this.errosFormulario.item = 'Selecione um produto';
      this.formularioValido = false;
    }

    if (this.itens.some(i => i.produto_id === this.novoItem.produto_id)) {
      this.errosFormulario.item = 'Produto já adicionado nesta nota';
      this.formularioValido = false;
    }

    if (!this.novoItem.quantidade || this.novoItem.quantidade <= 0) {
      this.errosFormulario.quantidade = 'Quantidade deve ser maior que zero';
      this.formularioValido = false;
    }

    // Validar saldo disponível
    if (this.novoItem.produto_id) {
      const produto = this.produtos.find(p => p.ID === this.novoItem.produto_id);
      if (produto && produto.saldo < this.novoItem.quantidade) {
        this.errosFormulario.quantidade = `Saldo insuficiente (disponível: ${produto.saldo})`;
        this.formularioValido = false;
      }
    }
  }

  nomeProduto(id: number): string {
    const p = this.produtos.find((x) => x.ID === id);
    return p ? `${p.codigo} — ${p.descricao}` : `Produto #${id}`;
  }

  obterSaldoProduto(id: number): number {
    return this.produtos.find(p => p.ID === id)?.saldo ?? 0;
  }

  adicionarItem() {
    this.validarItem();

    if (!this.formularioValido) {
      this.snackBar.open('Corrija os erros do formulário.', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }

    this.itens.push({ ...this.novoItem });
    this.snackBar.open(`✓ Produto adicionado (${this.novoItem.quantidade} un.)`, '', {
      duration: 2000,
      panelClass: 'snack-success',
    });
    this.limparFormularioItem();
  }

  removerItem(index: number) {
    this.itens.splice(index, 1);
    this.snackBar.open('Produto removido da nota.', '', {
      duration: 1500,
      panelClass: 'snack-info',
    });
  }

  limparFormularioItem() {
    this.novoItem = { produto_id: 0, quantidade: 1 };
    this.errosFormulario = { item: '', quantidade: '' };
    this.formularioValido = true;
  }

  salvar() {
    if (this.itens.length === 0) {
      this.snackBar.open('Adicione pelo menos um produto à nota fiscal.', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }

    this.salvando = true;
    this.notaService.criar({ itens: this.itens })
      .pipe(
        catchError((err) => {
          const msg = err.error?.erro || 'Erro ao criar nota fiscal. Tente novamente.';
          this.snackBar.open(msg, 'Fechar', {
            duration: 5000,
            panelClass: 'snack-error',
          });
          return of(null);
        }),
        finalize(() => {
          this.salvando = false;
        }),
        takeUntil(this.destroy$)
      )
      .subscribe((res) => {
        if (res) {
          this.stateService.adicionarNota(res);
          this.itens = [];
          this.limparFormularioItem();
          this.snackBar.open('✓ Nota fiscal criada com sucesso! NF #' + res.numero, 'Fechar', {
            duration: 4000,
            panelClass: 'snack-success',
          });
          // Recarregar para sincronizar
          setTimeout(() => this.carregarNotas(), 500);
        }
      });
  }

  imprimir(id: number) {
    // Confirmar antes de imprimir
    if (!confirm('Deseja imprimir esta nota fiscal? Esta ação não pode ser desfeita.')) {
      return;
    }

    this.imprimindo = id;
    this.notaService.imprimir(id)
      .pipe(
        catchError((err) => {
          const msg = err.error?.erro || 'Erro ao imprimir. Serviço indisponível.';
          this.snackBar.open(msg, 'Fechar', {
            duration: 6000,
            panelClass: 'snack-error',
          });
          return of(null);
        }),
        finalize(() => {
          this.imprimindo = null;
        }),
        takeUntil(this.destroy$)
      )
      .subscribe((res) => {
        if (res !== null) {
          // Atualizar estado
          const nota = this.notas.find(n => n.ID === id);
          if (nota) {
            this.stateService.atualizarNota({ ...nota, status: 'Fechada' });
          }

          this.snackBar.open('✓ Nota impressa! Status: FECHADA', 'Fechar', {
            duration: 4000,
            panelClass: 'snack-success',
          });

          // Recarregar ambas as telas
          setTimeout(() => {
            this.carregarNotas();
            this.carregarProdutos();
          }, 500);
        }
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
