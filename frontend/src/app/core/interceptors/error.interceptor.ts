import { Injectable } from '@angular/core';
import {
  HttpInterceptor,
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpErrorResponse,
} from '@angular/common/http';
import { Observable, throwError, timer } from 'rxjs';
import { catchError, retry, retryWhen, mergeMap } from 'rxjs/operators';
import { MatSnackBar } from '@angular/material/snack-bar';

@Injectable()
export class ErrorInterceptor implements HttpInterceptor {
  constructor(private snackBar: MatSnackBar) {}

  intercept(
    request: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    return next.handle(request).pipe(
      // Retry uma vez automaticamente em case de erro transitório
      retryWhen((errors) =>
        errors.pipe(
          mergeMap((error, index) => {
            if (index < 1 && this.isTransientError(error)) {
              return timer(1000); // Wait 1 second before retrying
            }
            return throwError(() => error);
          })
        )
      ),
      catchError((error: HttpErrorResponse) => {
        let errorMsg = 'Erro desconhecido';

        if (error.error instanceof ErrorEvent) {
          // Erro do cliente
          errorMsg = error.error.message;
        } else {
          // Erro do servidor
          errorMsg =
            error.error?.message ||
            error.error?.erro ||
            `Erro ${error.status}: ${error.statusText}`;
        }

        // Mostrar erro para o usuário
        this.showError(errorMsg);

        return throwError(() => error);
      })
    );
  }

  private isTransientError(error: HttpErrorResponse): boolean {
    // Retry em timeouts e erros 5xx (exceto 501)
    return (
      error.status === 0 ||
      error.status === 408 ||
      (error.status >= 500 && error.status !== 501)
    );
  }

  private showError(message: string): void {
    this.snackBar.open(message, 'Fechar', {
      duration: 5000,
      horizontalPosition: 'end',
      verticalPosition: 'top',
      panelClass: ['error-snackbar'],
    });
  }
}
