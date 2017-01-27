import { Injectable } from '@angular/core';
import { Http, Response, RequestOptions, URLSearchParams, Headers } from '@angular/http';
import { Observable } from 'rxjs/Observable';


@Injectable()
export class ApiService {

  public apiUrl = '/api/';

  constructor(private http: Http) { }

  getLinks(limit: number, offset: number, filter?: string): Observable<ILink[]> {
    let opts = this.doOptions({
      limit: limit,
      offset: offset,
      filter: filter
    });
    return this.http.get(this.apiUrl + 'links', opts)
      .map(this.extractData)
      .catch(this.handleError);
  }

  getLinkCount(filter?: string): Observable<any> {
    let opts = this.doOptions({
      filter: filter
    });
    return this.http.get(this.apiUrl + 'links/count', opts)
      .map(this.extractData)
      .catch(this.handleError);
  }

  doLogin(password: string): Observable<ILoginResponse> {
    return this.http.post(this.apiUrl + 'login', { 'password': password })
      .map(this.extractData)
      .catch(this.handleError);
  }

  getStats(): Observable<IStats> {
    return this.http.get(this.apiUrl + 'stat/all', this.doOptions())
      .map(this.extractData)
      .catch(this.handleError);
  }

  private doOptions(param?: any): RequestOptions {
    let params: URLSearchParams = new URLSearchParams();
    if (param) {
      for (let key of Object.keys(param)) {
        params.set(key, param[key]);
      }
    }
    let headers = new Headers({ 'Content-Type': 'application/json' });
    return new RequestOptions({
      headers: headers,
      search: params
    });
  }

  private extractData(res: Response) {
    let body = res.json();
    return body || {};
  }

  private handleError(error: Response | any) {
    let errMsg: string;
    if (error instanceof Response) {
      const body = error.json() || '';
      const err = body.error || JSON.stringify(body);
      errMsg = `${error.status} - ${error.statusText || ''} ${err}`;
    } else {
      errMsg = error.message ? error.message : error.toString();
    }
    console.error(errMsg);
    return Observable.throw(errMsg);
  }

}
