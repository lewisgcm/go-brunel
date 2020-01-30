import { Component, OnInit } from '@angular/core';
import { AuthService } from '../../service';
import { User } from '../../models';
import { Router } from '@angular/router';

@Component({
    selector: 'app-shared-navigation',
    templateUrl: './navigation.component.html',
    styleUrls: ['./navigation.component.scss'],
})
export class NavigationComponent implements OnInit {
    visible: boolean;
    user?: User;

    constructor(private authService: AuthService, private router: Router) {
        this.visible = false;
    }

    ngOnInit() {
        this.visible = this.authService.isAuthenticated();
        if (this.visible) {
            this.authService.getProfile().subscribe(
                (user) => {
                    this.user = user;
                }
            );
        }
    }

    logout() {
        this.authService.logOut().subscribe(
            () => {
                this.visible = false;
                this.router.navigate(['/user/login']);
            }
        );
    }
}
