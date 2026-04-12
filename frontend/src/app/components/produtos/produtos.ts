import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatIconModule } from '@angular/material/icon';
import { catchError, finalize, Subject, takeUntil } from 'rxjs';
import { of } from 'rxjs';
import { ProdutoService, Produto } from '../../services/produto';

@Component({
  selector: 'app-produtos',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatTableModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatCardModule,
    MatSnackBarModule,
    MatIconModule,
  ],
  templateUrl: './produtos.html',
  styleUrl: './produtos.css',
})
export class Produtos implements OnInit, OnDestroy {
  produtos: Produto[] = [];
  displayedColumns = ['codigo', 'descricao', 'saldo'];
  salvando = false;
  carregando = true;

  novoProduto = { codigo: '', descricao: '', saldo: 0 };

  private destroy$ = new Subject<void>();

  constructor(
    private produtoService: ProdutoService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit() {
    this.carregar();
  }

  carregar() {
    this.carregando = true;
    this.produtoService.listar()
      .pipe(
        catchError(() => {
          this.snackBar.open('Serviço de estoque indisponível. Tente novamente.', 'Fechar', {
            duration: 5000,
            panelClass: 'snack-error',
          });
          return of([]);
        }),
        finalize(() => {
          this.carregando = false;
        }),
        takeUntil(this.destroy$)
      )
      .subscribe((data) => {
        this.produtos = data;
      });
  }

  salvar() {
    if (!this.novoProduto.codigo || !this.novoProduto.descricao) {
      this.snackBar.open('Preencha todos os campos obrigatórios.', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }
    if (this.novoProduto.saldo < 0) {
      this.snackBar.open('Saldo não pode ser negativo.', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }

    this.salvando = true;
    this.produtoService.criar(this.novoProduto)
      .pipe(
        catchError((err) => {
          const msg = err.error?.erro || 'Erro ao cadastrar produto.';
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
          this.novoProduto = { codigo: '', descricao: '', saldo: 0 };
          this.carregar();
          this.snackBar.open('Produto cadastrado com sucesso!', 'OK', {
            duration: 3000,
            panelClass: 'snack-success',
          });
        }
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}