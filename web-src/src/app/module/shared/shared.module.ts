import { NgModule } from '@angular/core';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatListModule } from '@angular/material/list';
import { MatMenuModule } from '@angular/material/menu';
import { MatIconModule } from '@angular/material/icon';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

import { NavigationComponent } from './component';
import { MatButtonModule } from '@angular/material/button';

import { AuthService, AuthGuardService } from './service';

@NgModule({
    declarations: [
        NavigationComponent,
    ],
    providers: [
        AuthService,
        AuthGuardService,
    ],
    exports: [
        NavigationComponent,
    ],
    imports: [
        CommonModule,
        MatToolbarModule,
        MatButtonModule,
        MatListModule,
        MatMenuModule,
        MatIconModule,
        RouterModule,
    ]
})
export class SharedModule { }
