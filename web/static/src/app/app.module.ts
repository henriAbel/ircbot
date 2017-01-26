import { NgModule, ApplicationRef } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';

import { AppComponent } from './app.component';
import { ListComponent } from './content/list.component';
import { LoginComponent } from './content/login.component';
import { StatComponent } from './content/stat.component';
import { HeaderComponent } from './header/header.component';
import { ApiService } from './shared';
import { AuthService } from './auth/auth.service';
import { routing } from './app.routing';

import { HttpService } from './auth/http.service';
import { Http } from '@angular/http';

import { ChartModule } from 'angular2-chartjs';

@NgModule({
  imports: [
    BrowserModule,
    HttpModule,
    FormsModule,
    routing,
    ChartModule
  ],
  declarations: [
    AppComponent,
    ListComponent,
    LoginComponent,
    StatComponent,
    HeaderComponent
  ],
  providers: [
    ApiService,
    AuthService,
    { provide: Http, useClass: HttpService }
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
  constructor(public appRef: ApplicationRef) {}
}
