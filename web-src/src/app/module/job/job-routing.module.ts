import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { JobPageComponent } from './pages/job/job.page';

const routes: Routes = [
    { path: 'job/:id', component: JobPageComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class JobRoutingModule { }
