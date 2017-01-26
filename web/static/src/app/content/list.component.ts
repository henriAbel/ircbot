import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Data } from '@angular/router';
import {DomSanitizer} from '@angular/platform-browser';

import { ApiService } from '../shared/api.service';
import { AuthService } from '../auth/auth.service';

@Component({
  selector: 'irc-list',
  templateUrl: './list.component.html',
})
export class ListComponent implements OnInit {

  private links: ILink[];

  constructor(private route: ActivatedRoute, private service: ApiService, private sanitizer: DomSanitizer, private authService: AuthService) { }

  ngOnInit() {
    this.route.data.subscribe((data: Data) => {
        this.service.getLinks(10, 0, data['type']).subscribe(links => {
          this.links = links;
            console.log(links);
        });

        this.service.getLinkCount().subscribe(count => {
            console.log(count);
        });
    })
  }

  getRawLink(link: ILink): string {
    let type = link.Link_type == 'gif' ? 'webm' : link.Link_type
    let rawUrl = this.service.apiUrl + 'raw/' + link.Key + '/' + type;
    if (this.authService.hasToken()) {
      rawUrl += '?authorization=' + this.authService.getToken();
    }
    return rawUrl;
  }

  getYoutubeIdFromUrl(link: ILink) {
    return link.Link.substr(link.Link.lastIndexOf('=') + 1);
  }

  sanitize(url: string) {
    return this.sanitizer.bypassSecurityTrustResourceUrl(url);
  }

}
