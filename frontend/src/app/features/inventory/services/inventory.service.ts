import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../../../core/services/api.service';
import { Produto, CriarProdutoInput } from '../models/inventory.model';

@Injectable({ providedIn: 'root' })
export class InventoryService {
  private api = inject(ApiService);

  listar(): Observable<Produto[]> {
    return this.api.get<Produto[]>('/produtos');
  }

  buscarPorId(id: string): Observable<Produto> {
    return this.api.get<Produto>(`/produtos/${id}`);
  }

  criar(input: CriarProdutoInput): Observable<Produto> {
    return this.api.post<Produto>('/produtos', input);
  }

  atualizarSaldo(id: string, saldo: number): Observable<unknown> {
    return this.api.patch(`/produtos/${id}/saldo`, { saldo });
  }
}
