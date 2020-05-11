import {injectable} from 'inversify';
import {Observable, from} from 'rxjs';
import {switchMap} from 'rxjs/operators';

import {User} from './models';
import {AuthService} from './authService';

@injectable()
export class UserService {
	constructor(private _authService: AuthService) {
	}

	getProfile(): Observable<User> {
		return from(fetch(
			'/api/user/profile',
			{headers: this._authService.getAuthHeaders()},
		)).pipe(
			switchMap((response) => response.json()),
		);
	}
}
