import { Injectable } from '@angular/core';
import { HttpHeaders, HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';


import { User } from '..';
import { map } from 'rxjs/operators';

@Injectable()
export class AuthService {
    constructor(private httpClient: HttpClient) { }

    private getToken(): string {
        return localStorage.getItem('jwt');
    }

    public isAuthenticated(): boolean {
        return this.getToken() !== null;
    }

    public authHeaders(): HttpHeaders {
        return new HttpHeaders({
            Authorization: `Bearer ${this.getToken()}`
        });
    }

    public logOut(): Observable<void> {
        return of('').pipe(
            map(
                () => {
                    localStorage.removeItem('jwt');
                },
            )
        );
    }

    public getProfile(): Observable<User> {
        return this.httpClient.get<User>(
            `/api/user/profile`,
            { headers: this.authHeaders() }
        );
    }
}
