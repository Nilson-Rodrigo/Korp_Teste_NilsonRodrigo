import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  ReactiveFormsModule,
  FormBuilder,
  FormGroup,
  FormArray,
  Validators,
} from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar } from '@angular/material/snack-bar';
import { NotaService } from '../../core/api/nota.service';
import { ProdutoService } from '../../core/api/produto.service';
import { NotaFiscal, Produto } from '../../core/models';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-notas',
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
  templateUrl: './notas.html',
  styleUrl: './notas.css',
})
export class Notas implements OnInit, OnDestroy {
  form: FormGroup;
  notas: NotaFiscal[] = [];
  produtos: Produto[] = [];
  displayedColumns = ['numero', 'status', 'acoes'];
  displayedItensColumns = ['produto_id', 'quantidade', 'acoes'];
  loading = false;
  salvando = false;

  private destroy$ = new Subject<void>();

  constructor(
    private fb: FormBuilder,
    public notaService: NotaService,
    private produtoService: ProdutoService,
    private snackBar: MatSnackBar
  ) {
    this.form = this.fb.group({
      itens: this.fb.array([], Validators.minLength(1)),
    });
  }

  ngOnInit(): void {
    this.carregar();
    this.observarEstados();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  /**
   * Carrega notas e produtos
   */
  private carregar(): void {
    this.notaService.carregarNotas();
    this.produtoService.carregarProdutos();
  }

  /**
   * Observa estados do serviço
   */
  private observarEstados(): void {
    this.notaService.loading$
      .pipe(takeUntil(this.destroy$))
      .subscribe((loading) => {
        this.loading = loading;
      });

    this.notaService.notas$
      .pipe(takeUntil(this.destroy$))
      .subscribe((notas) => {
        this.notas = notas;
      });

    this.produtoService.produtos$
      .pipe(takeUntil(this.destroy$))
      .subscribe((produtos) => {
        this.produtos = produtos;
      });
  }

  /**
   * Retorna FormArray de itens
   */
  get itens(): FormArray {
    return this.form.get('itens') as FormArray;
  }

  /**
   * Cria novo item na FormArray
   */
  criarItemForm(): FormGroup {
    return this.fb.group({
      produto_id: ['', Validators.required],
      quantidade: [1, [Validators.required, Validators.min(1)]],
    });
  }

  /**
   * Adiciona novo item à nota
   */
  adicionarItem(): void {
    this.itens.push(this.criarItemForm());
  }

  /**
   * Remove item da nota
   */
  removerItem(index: number): void {
    this.itens.removeAt(index);
  }

  /**
   * Valida e cria nova nota fiscal
   */
  salvar(): void {
    if (!this.form.valid || this.itens.length === 0) {
      this.snackBar.open('Adicione ao menos um item à nota', 'Ok', {
        duration: 3000,
      });
      return;
    }

    this.salvando = true;

    this.notaService
      .criarNota({ itens: this.itens.value })
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (nota) => {
          this.form.reset();
          this.itens.clear();
          this.salvando = false;
          this.snackBar.open(`Nota fiscal #${nota.numero} criada!`, 'Ok', {
            duration: 3000,
          });
        },
        error: () => {
          this.salvando = false;
        },
      });
  }

  /**
   * Imprime uma nota fiscal
   */
  imprimir(id: number): void {
    this.notaService
      .imprimirNota(id)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: () => {
          this.snackBar.open('Nota impressa com sucesso!', 'Ok', {
            duration: 3000,
          });
        },
        error: () => {
          // Erro já é tratado pelo interceptor
        },
      });
  }

  /**
   * Verifica se nota pode ser impressa
   */
  podeImprimir(nota: NotaFiscal): boolean {
    return nota.status === 'Aberta';
  }

  /**
   * Retorna nome do produto
   */
  getNomeProduto(produtoId: number): string {
    const produto = this.produtos.find((p) => p.id === produtoId);
    return produto ? produto.descricao : `Produto #${produtoId}`;
  }

  /**
   * Valida campo de item
   */
  temErroItem(itemIndex: number, fieldName: string): boolean {
    const itemGroup = this.itens.at(itemIndex) as FormGroup;
    const field = itemGroup.get(fieldName);
    return !!(field && field.invalid && field.touched);
  }
}
