// Models para Produtos
export interface Produto {
  id: number;
  codigo: string;
  descricao: string;
  saldo: number;
}

export interface CriarProdutoRequest {
  codigo: string;
  descricao: string;
  saldo: number;
}

export interface AtualizarSaldoRequest {
  saldo: number;
}
