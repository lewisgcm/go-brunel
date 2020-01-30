import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

import { AuthService, User } from '../../shared';

@Injectable()
export class UserService {
    constructor(private httpClient: HttpClient, private authService: AuthService) { }
}
