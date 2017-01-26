import { Component, OnInit } from '@angular/core';
import { ApiService } from '../shared';
import { AuthService } from '../auth/auth.service';
import { Router, ActivatedRoute } from '@angular/router';


@Component({
  selector: 'irc-login',
  templateUrl: './login.component.html',
})
export class LoginComponent implements OnInit {

  model: string;
  loading: boolean;  
  errorMessage: string = '';

  constructor(private apiService: ApiService, private authService: AuthService, private router: Router, private route: ActivatedRoute) {
    this.loading = false;
  }

  ngOnInit() {
    this.route.data.subscribe(data => {
      if (undefined !== data) {
        if (data['action'] === 'logout') {
          this.authService.setToken(undefined);
          this.router.navigate(['/']);
        }
      }
    })
  }

  onSubmit() {
    this.loading = true;
    this.apiService.doLogin(this.model).subscribe(res => {
      this.authService.setToken(res.token);
      this.router.navigate(['/']);
    }, 
    err => {
      this.loading = !1;
      this.errorMessage = err;
    });
    console.log(this.model);
  }

}
