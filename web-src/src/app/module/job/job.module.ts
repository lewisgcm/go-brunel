import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { CommonModule } from '@angular/common';

import { JobRoutingModule } from './job-routing.module';
import { JobPageComponent } from './pages/job/job.page';
import { JobService } from './service/job.service';

@NgModule({
    declarations: [
        JobPageComponent,
    ],
    providers: [
        JobService,
    ],
    exports: [
        JobPageComponent,
    ],
    imports: [
        MatIconModule,
        MatButtonModule,
        CommonModule,
        JobRoutingModule,
    ]
})
export class JobModule { }
