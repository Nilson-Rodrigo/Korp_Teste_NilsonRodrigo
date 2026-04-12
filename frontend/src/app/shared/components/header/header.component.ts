import { Component, output, inject } from '@angular/core';
import { AuthService } from '../../../core/services/auth.service';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [],
  templateUrl: './header.component.html',
  styleUrl: './header.component.css',
})
export class HeaderComponent {
  toggleSidebar = output();

  private authService = inject(AuthService);
  user = this.authService.user;

  onToggleSidebar(): void {
    this.toggleSidebar.emit();
  }
}
