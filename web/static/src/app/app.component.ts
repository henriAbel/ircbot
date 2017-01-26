import { Component } from '@angular/core';

import '../style/app.scss';

@Component({
  selector: 'irc-app',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  url = 'https://github.com/preboot/angular2-webpack';
  title: string;

  constructor() {
    
  }
}
