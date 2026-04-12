import { Component, OnInit, inject, signal, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatIconModule } from '@angular/material/icon';
import { catchError } from 'rxjs/operators';
import { of } from 'rxjs';
import { InventoryService } from '../../services/inventory.service';
import { Produto, CriarProdutoInput } from '../../models/inventory.model';

@Component({
  selector: 'app-inventory-list',
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
  ],
  templateUrl: './inventory-list.component.html',
  styleUrl: './inventory-list.component.css',
})
export class InventoryListComponent implements OnInit {
  private inventoryService = inject(InventoryService);
  private snackBar = inject(MatSnackBar);
  private cdr = inject(ChangeDetectorRef);

  produtos = signal<Produto[]>([]);
  displayedColumns = ['codigo', 'descricao', 'saldo'];
  salvando = signal(false);
  carregando = signal(true);

  novoProduto: CriarProdutoInput = { codigo: '', descricao: '', saldo: 0 };

  ngOnInit(): void {
    this.carregar();
  }

  carregar(): void {
    this.carregando.set(true);
    this.inventoryService
      .listar()
      .pipe(
        catchError(() => {
          this.snackBar.open('Serviço de estoque indisponível. Tente novamente.', 'Fechar', {
            duration: 5000,
            panelClass: 'snack-error',
          });
          return of([]);
        })
      )
      .subscribe((data) => {
        this.produtos.set(data);
        this.carregando.set(false);
        this.cdr.markForCheck();
      });
  }

  salvar(): void {
    if (!this.novoProduto.codigo || !this.novoProduto.descricao) {
      this.snackBar.open('Preencha todos os campos obrigatórios.', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }
    if (this.novoProduto.saldo < 0) {
      this.snackBar.open('Saldo não pode ser negativo.', 'Fechar', {
        duration: 3000,
        panelClass: 'snack-error',
      });
      return;
    }

    this.salvando.set(true);
    this.inventoryService
      .criar(this.novoProduto)
      .pipe(
        catchError((err) => {
          const msg = err.error?.erro || 'Erro ao cadastrar produto.';
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
          this.novoProduto = { codigo: '', descricao: '', saldo: 0 };
          this.carregar();
          this.snackBar.open('Produto cadastrado com sucesso!', 'OK', {
            duration: 3000,
            panelClass: 'snack-success',
          });
        }
        this.salvando.set(false);
      });
  }
}
