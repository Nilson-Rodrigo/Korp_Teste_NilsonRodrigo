import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { tap, shareReplay } from 'rxjs/operators';
import { Produto, CriarProdutoRequest } from '../models';
import { ProdutoApiService } from '../api';

/**
 * Serviço de negócio para Produtos
 * Orquestra chamadas de API e gerencia estado local
 */
@Injectable({
  providedIn: 'root',
})
export class ProdutoService {
  private produtosSubject = new BehaviorSubject<Produto[]>([]);
  public produtos$ = this.produtosSubject.asObservable();

  private loadingSubject = new BehaviorSubject<boolean>(false);
  public loading$ = this.loadingSubject.asObservable();

  constructor(private api: ProdutoApiService) {}

  /**
   * Carrega todos os produtos
   */
  carregarProdutos(): void {
    this.loadingSubject.next(true);
    this.api.listar().pipe(
      tap((produtos) => {
        this.produtosSubject.next(produtos);
        this.loadingSubject.next(false);
      }),
      shareReplay(1)
    ).subscribe({
      error: () => this.loadingSubject.next(false),
    });
  }

  /**
   * Cria um novo produto
   */
  criarProduto(produto: CriarProdutoRequest): Observable<Produto> {
    return this.api.criar(produto).pipe(
      tap((novoProduto) => {
        const produtos = this.produtosSubject.value;
        this.produtosSubject.next([...produtos, novoProduto]);
      })
    );
  }

  /**
   * Busca um produto por ID
   */
  buscarProduto(id: number): Observable<Produto> {
    return this.api.buscarPorId(id);
  }

  /**
   * Atualiza o saldo de um produto
   */
  atualizarSaldo(id: number, novoSaldo: number): Observable<any> {
    return this.api.atualizarSaldo(id, novoSaldo).pipe(
      tap(() => {
        const produtos = this.produtosSubject.value.map((p) =>
          p.id === id ? { ...p, saldo: novoSaldo } : p
        );
        this.produtosSubject.next(produtos);
      })
    );
  }

  /**
   * Retorna os produtos atuais
   */
  getProdutosAtuais(): Produto[] {
    return this.produtosSubject.value;
  }
}
