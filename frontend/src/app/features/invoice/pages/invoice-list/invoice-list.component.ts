import { Component, OnInit, inject, signal, ChangeDetectorRef } from '@angular/core';
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
import { InvoiceService } from '../../services/invoice.service';
import { NotaFiscal, ItemNota } from '../../models/invoice.model';
import { InventoryService } from '../../../inventory/services/inventory.service';
import { Produto } from '../../../inventory/models/inventory.model';

@Component({
  selector: 'app-invoice-list',
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
  templateUrl: './invoice-list.component.html',
  styleUrl: './invoice-list.component.css',
})
export class InvoiceListComponent implements OnInit {
  private invoiceService = inject(InvoiceService);
  private inventoryService = inject(InventoryService);
  private snackBar = inject(MatSnackBar);
  private cdr = inject(ChangeDetectorRef);

  notas = signal<NotaFiscal[]>([]);
  produtos = signal<Produto[]>([]);
  displayedColumns = ['numero', 'status', 'itens', 'acoes'];

  novoItem: ItemNota = { produto_id: '', quantidade: 1 };
  itens = signal<ItemNota[]>([]);

  imprimindo = signal<string | null>(null);
  salvando = signal(false);
  carregando = signal(true);

  ngOnInit(): void {
    this.carregarNotas();
    this.carregarProdutos();
  }

  carregarProdutos(): void {
    this.inventoryService
      .listar()
      .pipe(
        catchError(() => {
          this.snackBar.open('Não foi possível carregar a lista de produtos.', 'Fechar', {
            duration: 5000,
            panelClass: 'snack-error',
          });
          return of([]);
        })
      )
      .subscribe((data) => {
        this.produtos.set(data);
        this.cdr.markForCheck();
      });
  }

  carregarNotas(): void {
    this.carregando.set(true);
    this.invoiceService
      .listar()
      .pipe(
        catchError(() => {
          this.snackBar.open('Serviço de faturamento indisponível.', 'Fechar', {
            duration: 5000,
            panelClass: 'snack-error',
          });
          this.carregando.set(false);
          return of([]);
        })
      )
      .subscribe((data) => {
        this.notas.set(data);
        this.carregando.set(false);
        this.cdr.markForCheck();
      });
  }

  nomeProduto(id: string): string {
    const p = this.produtos().find((x) => x.id === id);
    return p ? `${p.codigo} — ${p.descricao}` : `Produto #${id}`;
  }

  adicionarItem(): void {
    if (!this.novoItem.produto_id || this.novoItem.quantidade <= 0) {
      this.snackBar.open('Selecione um produto e informe a quantidade.', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }
    const duplicado = this.itens().find((i) => i.produto_id === this.novoItem.produto_id);
    if (duplicado) {
      this.snackBar.open('Este produto já foi adicionado à nota.', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }
    this.itens.update(current => [...current, { ...this.novoItem }]);
    this.novoItem = { produto_id: '', quantidade: 1 };
  }

  removerItem(index: number): void {
    this.itens.update(current => current.filter((_, i) => i !== index));
  }

  salvar(): void {
    if (this.itens().length === 0) {
      this.snackBar.open('Adicione pelo menos um produto à nota.', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }
    this.salvando.set(true);
    this.invoiceService
      .criar({ itens: this.itens() })
      .pipe(
        catchError((err) => {
          const msg = err.error?.erro || 'Erro ao criar nota fiscal.';
          this.snackBar.open(msg, 'Fechar', {
            duration: 5000,
            panelClass: 'snack-error',
          });
          this.salvando.set(false);
          return of(null);
        })
      )
      .subscribe((res) => {
        if (res) {
          this.itens.set([]);
          this.carregarNotas();
          this.snackBar.open('Nota fiscal criada com sucesso!', 'OK', {
            duration: 3000,
            panelClass: 'snack-success',
          });
        }
        this.salvando.set(false);
      });
  }

  imprimir(id: string): void {
    this.imprimindo.set(id);
    this.invoiceService
      .imprimir(id)
      .pipe(
        catchError((err) => {
          const msg = err.error?.erro || 'Serviço de estoque indisponível. Tente novamente.';
          this.snackBar.open(msg, 'Fechar', {
            duration: 6000,
            panelClass: 'snack-error',
          });
          this.imprimindo.set(null);
          return of(null);
        })
      )
      .subscribe((res) => {
        if (res !== null) {
          this.carregarNotas();
          this.carregarProdutos();
          this.snackBar.open('Nota impressa! Status atualizado para Fechada.', 'OK', {
            duration: 4000,
            panelClass: 'snack-success',
          });
        }
        this.imprimindo.set(null);
      });
  }
}
