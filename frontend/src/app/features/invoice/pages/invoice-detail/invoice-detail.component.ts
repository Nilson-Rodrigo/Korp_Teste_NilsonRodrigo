import { Component, OnInit, inject, signal, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { catchError } from 'rxjs/operators';
import { of } from 'rxjs';
import { InvoiceService } from '../../services/invoice.service';
import { NotaFiscal } from '../../models/invoice.model';
import { InventoryService } from '../../../inventory/services/inventory.service';
import { Produto } from '../../../inventory/models/inventory.model';

@Component({
  selector: 'app-invoice-detail',
  standalone: true,
  imports: [
    CommonModule,
    RouterLink,
    MatCardModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
  ],
  templateUrl: './invoice-detail.component.html',
  styleUrl: './invoice-detail.component.css',
})
export class InvoiceDetailComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private invoiceService = inject(InvoiceService);
  private inventoryService = inject(InventoryService);
  private snackBar = inject(MatSnackBar);
  private cdr = inject(ChangeDetectorRef);

  nota = signal<NotaFiscal | null>(null);
  produtos = signal<Produto[]>([]);
  carregando = signal(true);
  imprimindo = signal(false);

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.carregar(id);
      this.carregarProdutos();
    }
  }

  carregar(id: string): void {
    this.carregando.set(true);
    this.invoiceService
      .buscarPorId(id)
      .pipe(
        catchError(() => {
          this.snackBar.open('Nota não encontrada.', 'Fechar', { duration: 4000, panelClass: 'snack-error' });
          this.router.navigate(['/faturamento']);
          return of(null);
        })
      )
      .subscribe((data) => {
        this.nota.set(data);
        this.carregando.set(false);
        this.cdr.markForCheck();
      });
  }

  carregarProdutos(): void {
    this.inventoryService
      .listar()
      .pipe(catchError(() => of([])))
      .subscribe((data) => {
        this.produtos.set(data);
        this.cdr.markForCheck();
      });
  }

  nomeProduto(id: string): string {
    const p = this.produtos().find((x) => x.id === id);
    return p ? `${p.codigo} — ${p.descricao}` : `Produto #${id}`;
  }

  imprimir(): void {
    const n = this.nota();
    if (!n) return;
    this.imprimindo.set(true);
    this.invoiceService
      .imprimir(n.id)
      .pipe(
        catchError((err) => {
          const msg = err.error?.erro || 'Erro ao imprimir nota.';
          this.snackBar.open(msg, 'Fechar', { duration: 5000, panelClass: 'snack-error' });
          this.imprimindo.set(false);
          return of(null);
        })
      )
      .subscribe((res) => {
        if (res !== null) {
          this.snackBar.open('Nota impressa com sucesso!', 'OK', { duration: 3000, panelClass: 'snack-success' });
          this.carregar(n.id);
        }
        this.imprimindo.set(false);
      });
  }
}
