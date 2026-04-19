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
import { Observable } from 'rxjs';

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
  displayedColumns = ['id', 'codigo', 'descricao', 'saldo'];
  salvando = false;
  produtos$: Observable<Produto[]>;
  loading$: Observable<boolean>;

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

    this.produtos$ = this.produtoService.produtos$;
    this.loading$ = this.produtoService.loading$;
  }

  ngOnInit(): void {
    this.carregar();
  }

  ngOnDestroy(): void {}

  /**
   * Carrega os produtos
   */
  carregar(): void {
    this.produtoService.carregarProdutos();
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
