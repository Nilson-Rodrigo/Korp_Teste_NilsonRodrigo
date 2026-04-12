export interface ItemNota {
  produto_id: string;
  quantidade: number;
}

export interface NotaFiscal {
  id: string;
  numero: number;
  status: string;
  itens: ItemNota[];
  criado_em: string;
  atualizado_em: string;
}

export interface CriarNotaInput {
  itens: ItemNota[];
}
