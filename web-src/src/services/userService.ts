import {injectable} from 'inversify';
import {Observable, from} from 'rxjs';
import {switchMap} from 'rxjs/operators';

import {User, UserList} from './models';
import {AuthService} from './authService';
import {handleResponse} from './util';

@injectable()
export class UserService {
	constructor(private _authService: AuthService) {
	}

	getProfile(): Observable<User> {
		return from(fetch(
			'/api/user/profile',
			{headers: this._authService.getAuthHeaders()},
		)).pipe(
			switchMap(handleResponse),
		);
	}

	list(filter: string): Observable<UserList[]> {
		return from(fetch(
			`/api/user?filter=${filter}`,
			{headers: this._authService.getAuthHeaders()},
		)).pipe(
			switchMap(handleResponse),
		);
	}

	get(username: string): Observable<User> {
		return from(fetch(
			`/api/user/profile/${username}`,
			{headers: this._authService.getAuthHeaders()},
		)).pipe(
			switchMap(handleResponse),
		);
	}
}
