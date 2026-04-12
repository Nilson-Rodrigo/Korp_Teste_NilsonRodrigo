import { Component, OnInit, inject, signal, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { catchError } from 'rxjs/operators';
import { of } from 'rxjs';
import { InventoryService } from '../../services/inventory.service';
import { Produto } from '../../models/inventory.model';

@Component({
  selector: 'app-inventory-detail',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    RouterLink,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatSnackBarModule,
  ],
  templateUrl: './inventory-detail.component.html',
  styleUrl: './inventory-detail.component.css',
})
export class InventoryDetailComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private inventoryService = inject(InventoryService);
  private snackBar = inject(MatSnackBar);
  private cdr = inject(ChangeDetectorRef);

  produto = signal<Produto | null>(null);
  carregando = signal(true);
  novoSaldo = 0;
  salvando = signal(false);

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.carregar(id);
    }
  }

  carregar(id: string): void {
    this.carregando.set(true);
    this.inventoryService
      .buscarPorId(id)
      .pipe(
        catchError(() => {
          this.snackBar.open('Produto não encontrado.', 'Fechar', { duration: 4000, panelClass: 'snack-error' });
          this.router.navigate(['/estoque']);
          return of(null);
        })
      )
      .subscribe((data) => {
        this.produto.set(data);
        if (data) {
          this.novoSaldo = data.saldo;
        }
        this.carregando.set(false);
        this.cdr.markForCheck();
      });
  }

  atualizarSaldo(): void {
    const p = this.produto();
    if (!p || this.novoSaldo < 0) return;

    this.salvando.set(true);
    this.inventoryService
      .atualizarSaldo(p.id, this.novoSaldo)
      .pipe(
        catchError((err) => {
          const msg = err.error?.erro || 'Erro ao atualizar saldo.';
          this.snackBar.open(msg, 'Fechar', { duration: 5000, panelClass: 'snack-error' });
          this.salvando.set(false);
          return of(null);
        })
      )
      .subscribe((res) => {
        if (res !== null) {
          this.snackBar.open('Saldo atualizado com sucesso!', 'OK', { duration: 3000, panelClass: 'snack-success' });
          this.carregar(p.id);
        }
        this.salvando.set(false);
      });
  }
}
