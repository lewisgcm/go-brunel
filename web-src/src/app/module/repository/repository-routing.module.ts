import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { AuthGuardService as AuthGuard } from '../shared';
import { RepositoryPageComponent } from './pages/repository.page';

const routes: Routes = [
    { path: 'repository', component: RepositoryPageComponent, canActivate: [AuthGuard] },
    { path: 'repository/:id', component: RepositoryPageComponent, canActivate: [AuthGuard] },
    { path: '',   redirectTo: '/repository', pathMatch: 'full' },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class RepositoryRoutingModule { }
