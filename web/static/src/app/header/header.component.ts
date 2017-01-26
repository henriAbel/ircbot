import { Component, OnInit } from '@angular/core';
import { AuthService } from '../auth/auth.service';

@Component({
  selector: 'irc-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit {

  private loggedIn: boolean = false;
  constructor(private authService: AuthService) {
    this.loggedIn = this.authService.hasToken();
  }

  ngOnInit() {
    this.authService.onTokenUpdated.subscribe((newStatus: boolean) => {
      this.loggedIn = newStatus;
    });    
  }
}
