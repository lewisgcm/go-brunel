import { NgModule } from '@angular/core';
import { MatInputModule } from '@angular/material/input';
import { MatListModule } from '@angular/material/list';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSortModule } from '@angular/material/sort';
import { MatCardModule } from '@angular/material/card';
import { CommonModule } from '@angular/common';

import { SharedModule } from '../shared';

import { ListComponent, DetailComponent, JobListComponent } from './component';

import { RepositoryService } from './service';

import { RepositoryRoutingModule } from './repository-routing.module';
import { RepositoryPageComponent } from './pages/repository.page';

@NgModule({
    declarations: [
        ListComponent,
        DetailComponent,
        JobListComponent,
        RepositoryPageComponent,
    ],
    providers: [
        RepositoryService,
    ],
    exports: [
        ListComponent,
        RepositoryPageComponent,
    ],
    imports: [
        MatInputModule,
        MatListModule,
        MatTableModule,
        MatIconModule,
        MatTooltipModule,
        MatButtonModule,
        MatProgressBarModule,
        MatProgressSpinnerModule,
        MatPaginatorModule,
        MatSortModule,
        MatCardModule,
        CommonModule,
        SharedModule,
        RepositoryRoutingModule,
    ]
})
export class RepositoryModule { }
