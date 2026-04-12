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
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { catchError, finalize, Subject, takeUntil } from 'rxjs';
import { of } from 'rxjs';
import { ProdutoService, Produto } from '../../services/produto';
import { StateService } from '../../services/state.service';

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
    MatProgressSpinnerModule,
  ],
  templateUrl: './produtos.html',
  styleUrl: './produtos.css',
})
export class Produtos implements OnInit, OnDestroy {
  produtos: Produto[] = [];
  displayedColumns = ['codigo', 'descricao', 'saldo'];
  salvando = false;
  carregando = true;

  // Validação visual
  errosCampos = { codigo: '', descricao: '', saldo: '' };

  novoProduto = { codigo: '', descricao: '', saldo: 0 };
  formularioValido = false;

  private destroy$ = new Subject<void>();

  constructor(
    private produtoService: ProdutoService,
    private snackBar: MatSnackBar,
    private stateService: StateService
  ) {}

  ngOnInit() {
    this.carregar();
    // Subscrever a mudanças de estado
    this.stateService.produtos$
      .pipe(takeUntil(this.destroy$))
      .subscribe((produtos) => {
        this.produtos = produtos;
      });
  }

  carregar() {
    this.carregando = true;
    this.stateService.setCarregandoProdutos(true);
    
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
          this.stateService.setCarregandoProdutos(false);
        }),
        takeUntil(this.destroy$)
      )
      .subscribe((data) => {
        this.stateService.setProdutos(data);
      });
  }

  // Validação em tempo real
  validarCampos() {
    this.errosCampos = { codigo: '', descricao: '', saldo: '' };
    this.formularioValido = true;

    if (!this.novoProduto.codigo?.trim()) {
      this.errosCampos.codigo = 'Código do produto obrigatório';
      this.formularioValido = false;
    }

    if (this.novoProduto.codigo && this.produtos.some(p => p.codigo === this.novoProduto.codigo)) {
      this.errosCampos.codigo = 'Código já existe';
      this.formularioValido = false;
    }

    if (!this.novoProduto.descricao?.trim()) {
      this.errosCampos.descricao = 'Descrição obrigatória';
      this.formularioValido = false;
    }

    if (this.novoProduto.saldo == null || this.novoProduto.saldo < 0) {
      this.errosCampos.saldo = 'Saldo deve ser maior ou igual a zero';
      this.formularioValido = false;
    }

    return this.formularioValido;
  }

  salvar() {
    if (!this.validarCampos()) {
      this.snackBar.open('Corrija os erros do formulário', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }

    this.salvando = true;
    this.produtoService.criar(this.novoProduto)
      .pipe(
        catchError((err) => {
          const msg = err.error?.erro || 'Erro ao cadastrar produto. Tente novamente.';
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
          this.stateService.adicionarProduto(res);
          this.snackBar.open('✓ Produto cadastrado com sucesso!', 'Fechar', {
            duration: 3000,
            panelClass: 'snack-success',
          });
          this.limparFormulario();
        }
      });
  }

  limparFormulario() {
    this.novoProduto = { codigo: '', descricao: '', saldo: 0 };
    this.errosCampos = { codigo: '', descricao: '', saldo: '' };
    this.formularioValido = false;
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
