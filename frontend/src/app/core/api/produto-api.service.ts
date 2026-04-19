import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { Produto, CriarProdutoRequest } from '../models';

@Injectable({
  providedIn: 'root',
})
export class ProdutoApiService {
  private readonly apiUrl = 'http://localhost:8080/produtos';

  constructor(private http: HttpClient) {}

  /**
   * Lista todos os produtos
   */
  listar(): Observable<Produto[]> {
    return this.http.get<Produto[]>(this.apiUrl).pipe(
      tap(() => console.log('Produtos listados com sucesso')),
      catchError(this.handleError)
    );
  }

  /**
   * Cria um novo produto
   */
  criar(produto: CriarProdutoRequest): Observable<Produto> {
    return this.http.post<Produto>(this.apiUrl, produto).pipe(
      tap((p) => console.log(`Produto criado: ${p.id}`)),
      catchError(this.handleError)
    );
  }

  /**
   * Busca um produto por ID
   */
  buscarPorId(id: number): Observable<Produto> {
    return this.http.get<Produto>(`${this.apiUrl}/${id}`).pipe(
      tap(() => console.log(`Produto ${id} busca realizada`)),
      catchError(this.handleError)
    );
  }

  /**
   * Atualiza o saldo de um produto
   */
  atualizarSaldo(id: number, novoSaldo: number): Observable<any> {
    return this.http
      .patch(`${this.apiUrl}/${id}/saldo`, {
        saldo: novoSaldo,
      })
      .pipe(
        tap(() => console.log(`Saldo do produto ${id} atualizado`)),
        catchError(this.handleError)
      );
  }

  private handleError(error: any): Observable<never> {
    console.error('Erro na API de Produtos:', error);
    return throwError(() => error);
  }
}
