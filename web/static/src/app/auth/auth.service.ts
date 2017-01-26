import { Injectable, EventEmitter } from '@angular/core';

@Injectable()
export class AuthService {

  private token: string;
  onTokenUpdated: EventEmitter<boolean> = new EventEmitter();

  constructor() {
    this.token = sessionStorage.getItem('irc-token');
    if (this.token === null) this.token = undefined;
  }

  hasToken(): boolean {
    return this.token !== undefined && this.token.length > 0;
  }

  getToken(): string {
    return this.token;
  }

  setToken(token: string): void {
    this.token = token;
    if (this.token === undefined) {
      sessionStorage.removeItem('irc-token');
    }
    else {
      sessionStorage.setItem('irc-token', this.token);
    }
    this.onTokenUpdated.emit(this.hasToken());
  }  
}