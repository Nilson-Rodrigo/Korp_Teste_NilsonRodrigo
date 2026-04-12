import { Component } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  template: `
    <header class="navbar">
      <div class="navbar-inner">
        <a routerLink="/" class="navbar-brand">
          <div class="brand-icon">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M9 12h6M9 16h4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </div>
          <div class="brand-text">
            <span class="brand-name">Korp ERP</span>
            <span class="brand-sub">Notas Fiscais</span>
          </div>
        </a>

        <nav class="navbar-nav">
          <a routerLink="/produtos" routerLinkActive="active" class="nav-link">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            Produtos
          </a>
          <a routerLink="/notas" routerLinkActive="active" class="nav-link">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            Notas Fiscais
          </a>
        </nav>
      </div>
    </header>

    <main class="app-main">
      <router-outlet />
    </main>
  `,
  styles: [`
    .navbar {
      position: sticky;
      top: 0;
      z-index: 100;
      background: #fff;
      border-bottom: 1px solid #e2e8f0;
      box-shadow: 0 1px 3px rgb(0 0 0 / 0.06);
    }

    .navbar-inner {
      max-width: 1100px;
      margin: 0 auto;
      padding: 0 24px;
      height: 60px;
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .navbar-brand {
      display: flex;
      align-items: center;
      gap: 10px;
      text-decoration: none;
      color: inherit;
    }

    .brand-icon {
      width: 36px;
      height: 36px;
      border-radius: 8px;
      background: linear-gradient(135deg, #2563eb, #0ea5e9);
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      flex-shrink: 0;
    }

    .brand-text {
      display: flex;
      flex-direction: column;
      line-height: 1.2;
    }

    .brand-name {
      font-size: 0.9rem;
      font-weight: 700;
      color: #0f172a;
      letter-spacing: -0.01em;
    }

    .brand-sub {
      font-size: 0.7rem;
      font-weight: 400;
      color: #94a3b8;
      letter-spacing: 0.01em;
    }

    .navbar-nav {
      display: flex;
      align-items: center;
      gap: 4px;
    }

    .nav-link {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 7px 14px;
      border-radius: 8px;
      font-size: 0.875rem;
      font-weight: 500;
      color: #475569;
      text-decoration: none;
      transition: all 200ms cubic-bezier(0.4, 0, 0.2, 1);
    }

    .nav-link:hover {
      background: #f1f5f9;
      color: #0f172a;
    }

    .nav-link.active {
      background: #eff6ff;
      color: #2563eb;
    }

    .nav-link svg {
      opacity: 0.8;
    }

    .nav-link.active svg {
      opacity: 1;
    }

    .app-main {
      background-color: #f8fafc;
      min-height: calc(100vh - 60px);
    }

    @media (max-width: 600px) {
      .navbar-inner { padding: 0 16px; }
      .brand-sub { display: none; }
      .nav-link span { display: none; }
    }
  `]
})
export class AppComponent {}