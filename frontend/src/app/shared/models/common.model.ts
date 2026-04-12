export interface ApiResponse<T> {
  data: T;
  mensagem?: string;
}

export interface PageInfo {
  page: number;
  total: number;
  perPage: number;
}
