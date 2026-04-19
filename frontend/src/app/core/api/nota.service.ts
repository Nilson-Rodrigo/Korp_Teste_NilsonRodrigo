import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { tap, shareReplay, finalize } from 'rxjs/operators';
import { NotaFiscal, CriarNotaRequest } from '../models';
import { NotaApiService } from '../api';

/**
 * Serviço de negócio para Notas Fiscais
 * Orquestra chamadas de API e gerencia estado local
 */
@Injectable({
  providedIn: 'root',
})
export class NotaService {
  private notasSubject = new BehaviorSubject<NotaFiscal[]>([]);
  public notas$ = this.notasSubject.asObservable();

  private loadingSubject = new BehaviorSubject<boolean>(false);
  public loading$ = this.loadingSubject.asObservable();

  private imprimindoSubject = new BehaviorSubject<number | null>(null);
  public imprimindo$ = this.imprimindoSubject.asObservable();

  constructor(private api: NotaApiService) {}

  /**
   * Carrega todas as notas
   */
  carregarNotas(): void {
    this.loadingSubject.next(true);
    this.api.listar().pipe(
      tap((notas) => {
        this.notasSubject.next(notas);
        this.loadingSubject.next(false);
      }),
      shareReplay(1)
    ).subscribe({
      error: () => this.loadingSubject.next(false),
    });
  }

  /**
   * Cria uma nova nota fiscal
   */
  criarNota(nota: CriarNotaRequest): Observable<NotaFiscal> {
    return this.api.criar(nota).pipe(
      tap((novaNota) => {
        const notas = this.notasSubject.value;
        this.notasSubject.next([...notas, novaNota]);
      })
    );
  }

  /**
   * Busca uma nota por ID
   */
  buscarNota(id: number): Observable<NotaFiscal> {
    return this.api.buscarPorId(id);
  }

  /**
   * Imprime uma nota fiscal
   */
  imprimirNota(id: number): Observable<any> {
    this.imprimindoSubject.next(id);
    return this.api.imprimir(id).pipe(
      tap(() => {
        // Atualizar status da nota
        const notas = this.notasSubject.value.map((n) =>
          n.id === id ? { ...n, status: 'Fechada' as const } : n
        );
        this.notasSubject.next(notas);
      }),
      finalize(() => {
        // Sempre remove o spinner, mesmo em erro
        this.imprimindoSubject.next(null);
      })
    );
  }

  /**
   * Retorna as notas atuais
   */
  getNotasAtuais(): NotaFiscal[] {
    return this.notasSubject.value;
  }
}
