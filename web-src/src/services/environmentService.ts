import {injectable} from 'inversify';
import {Observable, from} from 'rxjs';
import {switchMap} from 'rxjs/operators';

import {AuthService} from './authService';
import {EnvironmentList, Environment} from './models';
import {handleResponse} from './util';

@injectable()
export class EnvironmentService {
	constructor(private _authService: AuthService) {
	}

	list(filter: string): Observable<EnvironmentList[]> {
		const query = new URLSearchParams({
			filter: filter.toString(),
		});

		return from(fetch(
			`/api/environment?${query}`,
			{headers: this._authService.getAuthHeaders()})).pipe(
			switchMap(handleResponse),
		);
	}

	get(id: string): Observable<Environment> {
		return from(fetch(
			`/api/environment/${id}`,
			{headers: this._authService.getAuthHeaders()})).pipe(
			switchMap(handleResponse),
		);
	}

	save(environment: Environment): Observable<Environment> {
		return from(fetch(
			`/api/environment`,
			{
				method: 'POST',
				headers: this._authService.getAuthHeaders(),
				body: JSON.stringify(environment),
			})).pipe(
			switchMap(handleResponse),
		);
	}
}
