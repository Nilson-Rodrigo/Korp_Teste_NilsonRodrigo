import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { NotaFiscal, CriarNotaRequest } from '../models';

@Injectable({
  providedIn: 'root',
})
export class NotaApiService {
  private readonly apiUrl = 'http://localhost:8081/notas';

  constructor(private http: HttpClient) {}

  /**
   * Lista todas as notas fiscais
   */
  listar(): Observable<NotaFiscal[]> {
    return this.http.get<NotaFiscal[]>(this.apiUrl).pipe(
      tap(() => console.log('Notas fiscais listadas')),
      catchError(this.handleError)
    );
  }

  /**
   * Cria uma nova nota fiscal
   */
  criar(nota: CriarNotaRequest): Observable<NotaFiscal> {
    return this.http.post<NotaFiscal>(this.apiUrl, nota).pipe(
      tap((n) => console.log(`Nota fiscal criada: ${n.id}`)),
      catchError(this.handleError)
    );
  }

  /**
   * Busca uma nota por ID
   */
  buscarPorId(id: number): Observable<NotaFiscal> {
    return this.http.get<NotaFiscal>(`${this.apiUrl}/${id}`).pipe(
      tap(() => console.log(`Nota fiscal ${id} buscada`)),
      catchError(this.handleError)
    );
  }

  /**
   * Imprime uma nota fiscal
   */
  imprimir(id: number): Observable<any> {
    return this.http.post(`${this.apiUrl}/${id}/imprimir`, {}).pipe(
      tap(() => console.log(`Nota fiscal ${id} impressa`)),
      catchError(this.handleError)
    );
  }

  private handleError(error: any): Observable<never> {
    console.error('Erro na API de Notas:', error);
    return throwError(() => error);
  }
}
