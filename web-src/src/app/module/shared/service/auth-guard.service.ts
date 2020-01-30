import { Injectable } from '@angular/core';
import { Router, CanActivate } from '@angular/router';

import { AuthService } from './auth.service';

@Injectable()
export class AuthGuardService implements CanActivate {
    constructor(public auth: AuthService, public router: Router) { }

    public canActivate(): boolean {
        if (!this.auth.isAuthenticated()) {
            this.router.navigate(['user/login']);
            return false;
        }
        return true;
    }
}
