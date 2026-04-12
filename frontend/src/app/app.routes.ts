import { Routes } from '@angular/router';
import { Produtos } from './components/produtos/produtos';
import { Notas } from './components/notas/notas';

export const routes: Routes = [
  { path: '', redirectTo: 'produtos', pathMatch: 'full' },
  { path: 'produtos', component: Produtos },
  { path: 'notas', component: Notas },
];