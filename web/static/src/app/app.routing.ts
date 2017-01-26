import { RouterModule, Routes } from '@angular/router';

import { ListComponent } from './content/list.component';
import { LoginComponent } from './content/login.component';
import { StatComponent } from './content/stat.component';


const routes: Routes = [
  { path: '', component: ListComponent, data: { type: 'image' } },
  { path: 'links', component: ListComponent, data: { type: 'link' } },
  { path: 'youtube', component: ListComponent, data: { type: 'youtube' } },
  { path: 'video', component: ListComponent, data: { type: 'webm,gif' } },
  { path: 'stats', component: StatComponent },
  { path: 'login', component: LoginComponent },
  { path: 'logout', component: LoginComponent, data: { action: 'logout' } },
];

export const routing = RouterModule.forRoot(routes);
