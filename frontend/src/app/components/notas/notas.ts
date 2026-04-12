import { Component, OnInit } from '@angular/core';
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
import { catchError } from 'rxjs/operators';
import { of } from 'rxjs';
import { NotaService, NotaFiscal, ItemNota } from '../../services/nota';
import { ProdutoService, Produto } from '../../services/produto';

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
export class Notas implements OnInit {
  // Dados
  notas: NotaFiscal[] = [];
  produtos: Produto[] = [];
  displayedColumns = ['numero', 'status', 'itens', 'acoes'];

  // Estado de formulário
  novoItem: ItemNota = { produto_id: 0, quantidade: 1 };
  itens: ItemNota[] = [];

  // Estado de UI
  imprimindo: number | null = null;
  salvando = false;
  carregando = true;

  constructor(
    private notaService: NotaService,
    private produtoService: ProdutoService,
    private snackBar: MatSnackBar
  ) {}

  // Ciclo de vida: OnInit — carrega dados ao montar o componente
  ngOnInit() {
    this.carregarNotas();
    this.carregarProdutos();
  }

  carregarProdutos() {
    this.produtoService.listar()
      .pipe(
        catchError(() => {
          this.snackBar.open('Não foi possível carregar a lista de produtos.', 'Fechar', {
            duration: 5000,
            panelClass: 'snack-error',
          });
          return of([]);
        })
      )
      .subscribe((data) => (this.produtos = data));
  }

  carregarNotas() {
    this.carregando = true;
    this.notaService.listar()
      .pipe(
        catchError(() => {
          this.snackBar.open('Serviço de faturamento indisponível.', 'Fechar', {
            duration: 5000,
            panelClass: 'snack-error',
          });
          this.carregando = false;
          return of([]);
        })
      )
      .subscribe((data) => {
        this.notas = data;
        this.carregando = false;
      });
  }

  nomeProduto(id: number): string {
    const p = this.produtos.find((x) => x.ID === id);
    return p ? `${p.codigo} — ${p.descricao}` : `Produto #${id}`;
  }

  adicionarItem() {
    if (!this.novoItem.produto_id || this.novoItem.quantidade <= 0) {
      this.snackBar.open('Selecione um produto e informe a quantidade.', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }
    const duplicado = this.itens.find((i) => i.produto_id === this.novoItem.produto_id);
    if (duplicado) {
      this.snackBar.open('Este produto já foi adicionado à nota.', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }
    this.itens.push({ ...this.novoItem });
    this.novoItem = { produto_id: 0, quantidade: 1 };
  }

  removerItem(index: number) {
    this.itens.splice(index, 1);
  }

  salvar() {
    if (this.itens.length === 0) {
      this.snackBar.open('Adicione pelo menos um produto à nota.', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }
    this.salvando = true;
    this.notaService.criar({ itens: this.itens })
      .pipe(
        catchError((err) => {
          const msg = err.error?.erro || 'Erro ao criar nota fiscal.';
          this.snackBar.open(msg, 'Fechar', {
            duration: 5000,
            panelClass: 'snack-error',
          });
          this.salvando = false;
          return of(null);
        })
      )
      .subscribe((res) => {
        if (res) {
          this.itens = [];
          this.carregarNotas();
          this.snackBar.open('Nota fiscal criada com sucesso!', 'OK', {
            duration: 3000,
            panelClass: 'snack-success',
          });
        }
        this.salvando = false;
      });
  }

  imprimir(id: number) {
    this.imprimindo = id;
    this.notaService.imprimir(id)
      .pipe(
        catchError((err) => {
          const msg = err.error?.erro || 'Serviço de estoque indisponível. Tente novamente.';
          this.snackBar.open(msg, 'Fechar', {
            duration: 6000,
            panelClass: 'snack-error',
          });
          this.imprimindo = null;
          return of(null);
        })
      )
      .subscribe((res) => {
        if (res !== null) {
          this.carregarNotas();
          this.carregarProdutos(); // Atualiza saldos exibidos
          this.snackBar.open('Nota impressa! Status atualizado para Fechada.', 'OK', {
            duration: 4000,
            panelClass: 'snack-success',
          });
        }
        this.imprimindo = null;
      });
  }
}