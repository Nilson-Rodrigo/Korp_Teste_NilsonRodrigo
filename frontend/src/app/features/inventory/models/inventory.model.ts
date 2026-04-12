export interface Produto {
  id: string;
  codigo: string;
  descricao: string;
  saldo: number;
  criado_em: string;
  atualizado_em: string;
}

export interface CriarProdutoInput {
  codigo: string;
  descricao: string;
  saldo: number;
}
