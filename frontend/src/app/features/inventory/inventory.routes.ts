import { Routes } from '@angular/router';
import { InventoryListComponent } from './pages/inventory-list/inventory-list.component';
import { InventoryDetailComponent } from './pages/inventory-detail/inventory-detail.component';

export const INVENTORY_ROUTES: Routes = [
  { path: '', component: InventoryListComponent },
  { path: ':id', component: InventoryDetailComponent },
];
