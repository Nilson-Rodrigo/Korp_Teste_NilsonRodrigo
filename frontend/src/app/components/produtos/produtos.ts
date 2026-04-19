import { Component, OnInit, OnDestroy, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ProdutoService, Produto } from '../../services/produto';

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
    MatIconModule,
    MatChipsModule,
  ],
  templateUrl: './produtos.html',
  styleUrl: './produtos.css',
})
export class Produtos implements OnInit, OnDestroy {
  private destroyRef = inject(takeUntilDestroyed);
  
  private fb = inject(FormBuilder);
  private produtoService = inject(ProdutoService);

  form!: FormGroup;
  produtos: Produto[] = [];
  displayedColumns = ['id', 'codigo', 'descricao', 'saldo'];
  loading = false;

  ngOnInit() {
    this.form = this.fb.group({
      codigo: ['', [Validators.required, Validators.minLength(2), Validators.maxLength(20)]],
      descricao: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(255)]],
      saldo: [0, [Validators.required, Validators.min(0)]],
    });
    this.carregar();
  }

  carregar() {
    this.loading = true;
    this.produtoService
      .listar()
      .pipe(takeUntilDestroyed())
      .subscribe({
        next: (data) => {
          this.produtos = data;
          this.loading = false;
        },
        error: (_) => {
          this.loading = false;
        },
      });
  }

  salvar() {
    if (this.form.invalid) return;
    
    this.loading = true;
    this.produtoService
      .criar(this.form.value)
      .pipe(takeUntilDestroyed())
      .subscribe({
        next: () => {
          this.form.reset({ saldo: 0 });
          this.carregar();
        },
        error: (_) => {
          this.loading = false;
        },
      });
  }

  temErro(campo: string): boolean {
    const control = this.form.get(campo);
    return !!(control && control.invalid && (control.dirty || control.touched));
  }

  getMensagemErro(campo: string): string {
    const control = this.form.get(campo);
    if (control?.hasError('required')) {
      return `${campo} é obrigatório`;
    }
    if (control?.hasError('minlength')) {
      return `${campo} mínimo de ${control.errors?.['minlength'].requiredLength} caracteres`;
    }
    if (control?.hasError('maxlength')) {
      return `${campo} máximo de ${control.errors?.['maxlength'].requiredLength} caracteres`;
    }
    if (control?.hasError('min')) {
      return `${campo} mínimo de 0`;
    }
    return '';
  }

  ngOnDestroy() {
    // takeUntilDestroyed já trata disso automaticamente
  }
}