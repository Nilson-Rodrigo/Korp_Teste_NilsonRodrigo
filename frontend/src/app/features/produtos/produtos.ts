import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  ReactiveFormsModule,
  FormBuilder,
  FormGroup,
  Validators,
} from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ProdutoService } from '../../core/api/produto.service';
import { Produto } from '../../core/models';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-produtos',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatTableModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatCardModule,
    MatProgressSpinnerModule,
  ],
  templateUrl: './produtos.html',
  styleUrl: './produtos.css',
})
export class Produtos implements OnInit, OnDestroy {
  form: FormGroup;
  produtos: Produto[] = [];
  displayedColumns = ['id', 'codigo', 'descricao', 'saldo'];
  loading = false;
  salvando = false;

  private destroy$ = new Subject<void>();

  constructor(
    private fb: FormBuilder,
    private produtoService: ProdutoService,
    private snackBar: MatSnackBar
  ) {
    this.form = this.fb.group({
      codigo: [
        '',
        [
          Validators.required,
          Validators.minLength(1),
          Validators.maxLength(50),
        ],
      ],
      descricao: [
        '',
        [
          Validators.required,
          Validators.minLength(3),
          Validators.maxLength(200),
        ],
      ],
      saldo: [0, [Validators.required, Validators.min(0), Validators.max(999999)]],
    });
  }

  ngOnInit(): void {
    this.carregar();
    this.observarLoading();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  /**
   * Carrega os produtos
   */
  carregar(): void {
    this.produtoService.carregarProdutos();
  }

  /**
   * Observa estado de loading
   */
  private observarLoading(): void {
    this.produtoService.loading$
      .pipe(takeUntil(this.destroy$))
      .subscribe((loading) => {
        this.loading = loading;
      });

    this.produtoService.produtos$
      .pipe(takeUntil(this.destroy$))
      .subscribe((produtos) => {
        this.produtos = produtos;
      });
  }

  /**
   * Valida e cria um novo produto
   */
  salvar(): void {
    if (!this.form.valid) {
      this.snackBar.open('Preencha todos os campos corretamente', 'Ok', {
        duration: 3000,
      });
      return;
    }

    this.salvando = true;

    this.produtoService
      .criarProduto(this.form.value)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (produto) => {
          this.form.reset();
          this.salvando = false;
          this.snackBar.open(
            `Produto ${produto.codigo} criado com sucesso!`,
            'Ok',
            { duration: 3000 }
          );
        },
        error: () => {
          this.salvando = false;
        },
      });
  }

  /**
   * Retorna se um campo tem erro
   */
  temErro(fieldName: string, tipoErro: string): boolean {
    const field = this.form.get(fieldName);
    return !!(field && field.hasError(tipoErro) && field.touched);
  }

  /**
   * Retorna mensagem de erro
   */
  getMensagemErro(fieldName: string): string {
    const field = this.form.get(fieldName);
    if (field?.hasError('required')) {
      return `${fieldName.charAt(0).toUpperCase() + fieldName.slice(1)} é obrigatório`;
    }
    if (field?.hasError('minlength')) {
      return `Mínimo de ${field.getError('minlength').requiredLength} caracteres`;
    }
    if (field?.hasError('maxlength')) {
      return `Máximo de ${field.getError('maxlength').requiredLength} caracteres`;
    }
    if (field?.hasError('min')) {
      return `Valor mínimo é ${field.getError('min').min}`;
    }
    return 'Campo inválido';
  }
}
