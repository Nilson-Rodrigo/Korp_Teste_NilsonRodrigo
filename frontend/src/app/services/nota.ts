import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

export interface ItemNota {
  produto_id: number;
  quantidade: number;
}

export interface NotaFiscal {
  ID: number;
  numero: number;
  status: string;
  itens: ItemNota[];
}

@Injectable({
  providedIn: 'root'
})
export class NotaService {
  private apiUrl = environment.faturamentoApiUrl;

  constructor(private http: HttpClient) {}

  listar(): Observable<NotaFiscal[]> {
    return this.http.get<NotaFiscal[]>(`${this.apiUrl}/notas`);
  }

  criar(nota: { itens: ItemNota[] }): Observable<NotaFiscal> {
    return this.http.post<NotaFiscal>(`${this.apiUrl}/notas`, nota);
  }

  imprimir(id: number): Observable<any> {
    return this.http.post(`${this.apiUrl}/notas/${id}/imprimir`, {});
  }
}