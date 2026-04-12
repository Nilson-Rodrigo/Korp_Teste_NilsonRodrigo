import { Component, output, inject } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [RouterLink, RouterLinkActive],
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
