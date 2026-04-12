import { Injectable, signal, computed } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private currentUser = signal<{ nome: string; email: string } | null>({
    nome: 'Administrador',
    email: 'admin@korp.com.br',
  });

  readonly user = this.currentUser.asReadonly();
  readonly isAuthenticated = computed(() => this.currentUser() !== null);

  login(nome: string, email: string): void {
    this.currentUser.set({ nome, email });
  }

  logout(): void {
    this.currentUser.set(null);
  }
}
