import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { Produto } from './produto';
import { NotaFiscal } from './nota';

@Injectable({
  providedIn: 'root'
})
export class StateService {
  
  // Produtos
  private produtosSubject = new BehaviorSubject<Produto[]>([]);
  public produtos$ = this.produtosSubject.asObservable();

  // Notas
  private notasSubject = new BehaviorSubject<NotaFiscal[]>([]);
  public notas$ = this.notasSubject.asObservable();

  // Estados de carregamento
  private carregandoProdutosSubject = new BehaviorSubject<boolean>(false);
  public carregandoProdutos$ = this.carregandoProdutosSubject.asObservable();

  private carregandoNotasSubject = new BehaviorSubject<boolean>(false);
  public carregandoNotas$ = this.carregandoNotasSubject.asObservable();

  // Notificações
  private notificacaoSubject = new BehaviorSubject<{ mensagem: string; tipo: 'sucesso' | 'erro' | 'aviso'; id: number } | null>(null);
  public notificacao$ = this.notificacaoSubject.asObservable();

  constructor() {}

  // ===== PRODUTOS =====
  setProdutos(produtos: Produto[]) {
    this.produtosSubject.next(produtos);
  }

  getProdutos(): Produto[] {
    return this.produtosSubject.value;
  }

  adicionarProduto(produto: Produto) {
    const produtos = this.produtosSubject.value;
    this.produtosSubject.next([...produtos, produto]);
  }

  atualizarProduto(produto: Produto) {
    const produtos = this.produtosSubject.value;
    const index = produtos.findIndex(p => p.ID === produto.ID);
    if (index !== -1) {
      produtos[index] = produto;
      this.produtosSubject.next([...produtos]);
    }
  }

  setCarregandoProdutos(isLoading: boolean) {
    this.carregandoProdutosSubject.next(isLoading);
  }

  // ===== NOTAS =====
  setNotas(notas: NotaFiscal[]) {
    this.notasSubject.next(notas);
  }

  getNotas(): NotaFiscal[] {
    return this.notasSubject.value;
  }

  adicionarNota(nota: NotaFiscal) {
    const notas = this.notasSubject.value;
    this.notasSubject.next([...notas, nota]);
  }

  atualizarNota(nota: NotaFiscal) {
    const notas = this.notasSubject.value;
    const index = notas.findIndex(n => n.ID === nota.ID);
    if (index !== -1) {
      notas[index] = nota;
      this.notasSubject.next([...notas]);
    }
  }

  setCarregandoNotas(isLoading: boolean) {
    this.carregandoNotasSubject.next(isLoading);
  }

  // ===== NOTIFICAÇÕES =====
  mostrarSucesso(mensagem: string) {
    this.mostrarNotificacao(mensagem, 'sucesso');
  }

  mostrarErro(mensagem: string) {
    this.mostrarNotificacao(mensagem, 'erro');
  }

  mostrarAviso(mensagem: string) {
    this.mostrarNotificacao(mensagem, 'aviso');
  }

  private mostrarNotificacao(mensagem: string, tipo: 'sucesso' | 'erro' | 'aviso') {
    const id = Date.now();
    this.notificacaoSubject.next({ mensagem, tipo, id });
  }

  limparNotificacao() {
    this.notificacaoSubject.next(null);
  }
}
