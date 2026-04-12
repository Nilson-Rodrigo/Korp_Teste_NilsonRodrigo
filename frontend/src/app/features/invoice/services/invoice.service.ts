import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../../../core/services/api.service';
import { NotaFiscal, CriarNotaInput } from '../models/invoice.model';

@Injectable({ providedIn: 'root' })
export class InvoiceService {
  private api = inject(ApiService);

  listar(): Observable<NotaFiscal[]> {
    return this.api.get<NotaFiscal[]>('/notas');
  }

  buscarPorId(id: string): Observable<NotaFiscal> {
    return this.api.get<NotaFiscal>(`/notas/${id}`);
  }

  criar(input: CriarNotaInput): Observable<NotaFiscal> {
    return this.api.post<NotaFiscal>('/notas', input);
  }

  imprimir(id: string): Observable<unknown> {
    return this.api.post(`/notas/${id}/imprimir`, {});
  }
}
