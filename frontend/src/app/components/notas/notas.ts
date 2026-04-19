import { Component, OnInit, OnDestroy, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, FormArray, Validators } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatSelectModule } from '@angular/material/select';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { NotaService, NotaFiscal } from '../../services/nota';
import { ProdutoService, Produto } from '../../services/produto';

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
    MatIconModule,
    MatChipsModule,
    MatSelectModule,
  ],
  templateUrl: './notas.html',
  styleUrl: './notas.css',
})
export class Notas implements OnInit, OnDestroy {
  private destroyRef = inject(takeUntilDestroyed);
  
  private fb = inject(FormBuilder);
  notaService = inject(NotaService);
  produtoService = inject(ProdutoService);

  form!: FormGroup;
  notas: NotaFiscal[] = [];
  produtos: Produto[] = [];
  displayedColumns = ['numero', 'status', 'acoes'];
  imprimindo: number | null = null;
  loading = false;

  ngOnInit() {
    this.form = this.fb.group({
      itens: this.fb.array([], Validators.minLength(1)),
    });
    this.carregarProdutos();
    this.carregar();
  }

  get itens(): FormArray {
    return this.form.get('itens') as FormArray;
  }

  criarItemForm(): FormGroup {
    return this.fb.group({
      produto_id: [null, [Validators.required]],
      quantidade: [1, [Validators.required, Validators.min(1)]],
    });
  }

  adicionarItem() {
    this.itens.push(this.criarItemForm());
  }

  removerItem(index: number) {
    this.itens.removeAt(index);
  }

  carregarProdutos() {
    this.produtoService
      .listar()
      .pipe(takeUntilDestroyed())
      .subscribe({
        next: (data) => {
          this.produtos = data;
        },
      });
  }

  carregar() {
    this.loading = true;
    this.notaService
      .listar()
      .pipe(takeUntilDestroyed())
      .subscribe({
        next: (data) => {
          this.notas = data;
          this.loading = false;
        },
        error: (_) => {
          this.loading = false;
        },
      });
  }

  salvar() {
    console.log('DEBUG: Clicou em salvar');
    console.log('DEBUG: form.invalid =', this.form.invalid);
    console.log('DEBUG: itens.length =', this.itens.length);
    
    if (this.form.invalid || this.itens.length === 0) {
      console.warn('Formulário inválido ou sem itens');
      return;
    }

    this.loading = true;
    
    const itensTransformados = this.itens.value.map((item: any) => ({
      produto_id: Number(item.produto_id),
      quantidade: Number(item.quantidade),
    }));

    const payload = { itens: itensTransformados };

    this.notaService
      .criar(payload)
      .pipe(takeUntilDestroyed())
      .subscribe({
        next: (response) => {
          console.log('DEBUG: Nota criada! Response:', response);
          
          // Limpar todos os itens
          while (this.itens.length > 0) {
            this.itens.removeAt(0);
          }
          
          // Adicionar um novo item vazio
          this.adicionarItem();
          console.log('DEBUG: Item novo adicionado. Itens agora:', this.itens.length);
          
          // Atualizar notas SEM usar carregar() que pode resetar o form
          this.notas.push(response);
          
          this.loading = false;
          console.log('DEBUG: Tudo pronto para criar a próxima nota');
        },
        error: (err) => {
          console.error('DEBUG: Erro:', err);
          this.loading = false;
        },
      });
  }

  compareFn(c1: any, c2: any): boolean {
    return c1 && c2 ? c1 === c2 : c1 === c2;
  }

  imprimir(id: number) {
    this.imprimindo = id;
    this.notaService
      .imprimir(id)
      .pipe(takeUntilDestroyed())
      .subscribe({
        next: () => {
          this.imprimindo = null;
          this.carregar();
        },
        error: (_) => {
          this.imprimindo = null;
        },
      });
  }

  getNomeProduto(produtoId: number): string {
    return this.produtos.find((p) => p.id === produtoId)?.descricao || '';
  }

  temErroItem(index: number, campo: string): boolean {
    const control = this.itens.at(index)?.get(campo);
    return !!(control && control.invalid && (control.dirty || control.touched));
  }

  ngOnDestroy() {
    // takeUntilDestroyed já trata disso automaticamente
  }
}