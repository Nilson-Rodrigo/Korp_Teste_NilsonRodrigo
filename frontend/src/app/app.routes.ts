import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', redirectTo: 'estoque', pathMatch: 'full' },
  {
    path: 'estoque',
    loadChildren: () =>
      import('./features/inventory/inventory.routes').then((m) => m.INVENTORY_ROUTES),
  },
  {
    path: 'faturamento',
    loadChildren: () =>
      import('./features/invoice/invoice.routes').then((m) => m.INVOICE_ROUTES),
  },
];