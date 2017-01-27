import { Component, OnInit, ViewChildren, QueryList, ElementRef, HostListener } from '@angular/core';
import { ActivatedRoute, Data } from '@angular/router';
import { DomSanitizer } from '@angular/platform-browser';

import { ApiService } from '../shared/api.service';
import { AuthService } from '../auth/auth.service';

@Component({
  selector: 'irc-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss'],
})
export class ListComponent implements OnInit {

  @ViewChildren('thumbImage') private thumbnails: QueryList<ElementRef>;

  links: ILink[];
  page = 0;
  linkCount: number = -1;
  itemsInPage = 10;
  private filter: string;

  constructor(private route: ActivatedRoute, private service: ApiService,
    private sanitizer: DomSanitizer, private authService: AuthService) { }

  @HostListener('window:keydown', ['$event'])
  doSomething($event) {
    let index;
    switch ($event.keyCode) {
      case 39:
        index = this.getActivatedElementIndex() + 1;
        if (this.links[index].Link_type === 'fullView') {
          index++;
        };
        if (index >= this.links.length) {
          this.next(() => {
            this.activate(this.links[0]);
          });
          return;
        }
        this.activate(this.links[index]);
        break;
      case 37:
        index = this.getActivatedElementIndex() - 1;
        if (index < 0) {
          this.prev(() => {
            this.activate(this.links[this.links.length - 1]);
          });
        }
        this.activate(this.links[index]);
        break;
      case 27:
        this.links.splice(this.getFullViewIndex(), 1);
        break;
    }
  }

  ngOnInit() {
    this.route.data.subscribe((data: Data) => {
      this.filter = data['type'];
      this.service.getLinkCount(this.filter).subscribe(res => {
        this.linkCount = res.Count;
        this.refreshLinks();
      });
    });

  }

  refreshLinks(callback?: () => any) {
    let offset = this.page * this.itemsInPage;
    this.service.getLinks(this.itemsInPage, offset, this.filter).subscribe(links => {
      this.links = links;
      if (callback) {
        callback.call(this);
      }
    });
  }

  getRawLink(link: ILink, forceType?: string): string {
    let type = forceType ? forceType : (link.Link_type === 'gif' ? 'webm' : link.Link_type);
    let rawUrl = this.service.apiUrl + 'raw/' + link.Key + '/' + type;
    if (this.authService.hasToken()) {
      rawUrl += '?authorization=' + this.authService.getToken();
    }
    return rawUrl;
  }

  // TODO might not always work
  getYoutubeIdFromUrl(link: ILink): string {
    return link.Link.substr(link.Link.lastIndexOf('=') + 1);
  }

  // For Youtube iframe src attribute, url must be sanitized
  sanitize(url: string) {
    return this.sanitizer.bypassSecurityTrustResourceUrl(url);
  }

  next(callback?: () => any): void {
    if (this.itemsInPage * this.page < this.linkCount - this.itemsInPage) {
      this.page++;
      this.refreshLinks(callback);
    }
  }

  prev(callback?: () => any): void {
    if (this.page > 0) {
      this.page--;
      this.refreshLinks(callback);
    }
  }

  activate(link: ILink) {
    if (link === undefined) {
      return;
    }
    // Remove any open images
    let fvi = this.getFullViewIndex();
    if (fvi >= 0) {
      this.links.splice(fvi, 1);
    }
    // Make copy of current image
    let fullView: any = {};
    for (let property in link) {
      if (link.hasOwnProperty(property)) {
        fullView[property] = link[property];
      }
    }
    fullView.Link_type = 'fullView';
    // Insert it at end of row
    this.links.splice(this.getRowLastElement(link), 0, fullView);
  }

  getFullViewIndex(): number {
    for (let i = 0; i < this.links.length; i++) {
      if (this.links[i].Link_type === 'fullView') {
        return i;
      }
    }
    return -1;
  }

  getActivatedElementIndex() {
    let fullView = this.getFullViewIndex();
    if (-1 === fullView) {
      return -1;
    }
    let dbId = this.links[fullView].Key;
    for (let i = 0; i < this.links.length; i++) {
      if (this.links[i].Key === dbId) {
        return i;
      }
    }
    return -1;
  }

  getRowLastElement(link: ILink) {
    let lastOffsetTop = -1;
    let elements = this.thumbnails.toArray();
    // Find element where offsetTop differs from previous element
    for (let i = this.links.indexOf(link); i < elements.length; i++) {
      let ofst = elements[i].nativeElement.offsetTop;
      if (lastOffsetTop < 0) {
        lastOffsetTop = ofst;
        continue;
      }
      if (ofst !== lastOffsetTop) {
        return i;
      }
    }
    // Last element is clicked
    return elements.length;
  }

  closeFullView() {
    this.links.splice(this.getFullViewIndex(), 1);
  }
}
