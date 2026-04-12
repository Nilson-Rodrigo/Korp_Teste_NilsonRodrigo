import { Component, OnInit, inject } from '@angular/core';
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

  nota: NotaFiscal | null = null;
  produtos: Produto[] = [];
  carregando = true;
  imprimindo = false;

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.carregar(id);
      this.carregarProdutos();
    }
  }

  carregar(id: string): void {
    this.carregando = true;
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
        this.nota = data;
        this.carregando = false;
      });
  }

  carregarProdutos(): void {
    this.inventoryService
      .listar()
      .pipe(catchError(() => of([])))
      .subscribe((data) => (this.produtos = data));
  }

  nomeProduto(id: string): string {
    const p = this.produtos.find((x) => x.id === id);
    return p ? `${p.codigo} — ${p.descricao}` : `Produto #${id}`;
  }

  imprimir(): void {
    if (!this.nota) return;
    this.imprimindo = true;
    this.invoiceService
      .imprimir(this.nota.id)
      .pipe(
        catchError((err) => {
          const msg = err.error?.erro || 'Erro ao imprimir nota.';
          this.snackBar.open(msg, 'Fechar', { duration: 5000, panelClass: 'snack-error' });
          this.imprimindo = false;
          return of(null);
        })
      )
      .subscribe((res) => {
        if (res !== null) {
          this.snackBar.open('Nota impressa com sucesso!', 'OK', { duration: 3000, panelClass: 'snack-success' });
          this.carregar(this.nota!.id);
        }
        this.imprimindo = false;
      });
  }
}
