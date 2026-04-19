// Models para Notas Fiscais
export interface ItemNota {
  id?: number;
  produto_id: number;
  quantidade: number;
}

export interface NotaFiscal {
  id: number;
  numero: number;
  status: 'Aberta' | 'Fechada';
  itens: ItemNota[];
}

export interface CriarNotaRequest {
  itens: ItemNota[];
}
