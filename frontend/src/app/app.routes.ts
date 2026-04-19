import { Routes } from '@angular/router';
import { Produtos } from './features/produtos/produtos';
import { Notas } from './features/notas/notas';

export const routes: Routes = [
  { path: '', redirectTo: 'produtos', pathMatch: 'full' },
  { path: 'produtos', component: Produtos },
  { path: 'notas', component: Notas },
];