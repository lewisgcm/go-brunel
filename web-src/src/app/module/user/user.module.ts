import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';

import { UserRoutingModule } from './user-routing.module';
import { LoginPageComponent } from './pages/login/login.page';

@NgModule({
    declarations: [
        LoginPageComponent,
    ],
    exports: [
        LoginPageComponent,
    ],
    imports: [
        MatButtonModule,
        UserRoutingModule,
    ]
})
export class UserModule { }
