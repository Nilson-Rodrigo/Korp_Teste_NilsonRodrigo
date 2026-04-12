import { HttpInterceptorFn, HttpErrorResponse } from '@angular/common/http';
import { catchError, throwError } from 'rxjs';

export const httpErrorInterceptor: HttpInterceptorFn = (req, next) => {
  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      let mensagem = 'Erro desconhecido. Tente novamente.';

      if (error.status === 0) {
        mensagem = 'Servidor indisponível. Verifique a conexão.';
      } else if (error.status === 404) {
        mensagem = 'Recurso não encontrado.';
      } else if (error.status === 400) {
        mensagem = error.error?.erro || 'Dados inválidos.';
      } else if (error.status >= 500) {
        mensagem = 'Erro interno do servidor.';
      }

      console.error(`[HTTP ${error.status}] ${mensagem}`, error);
      return throwError(() => ({ ...error, mensagem }));
    })
  );
};
