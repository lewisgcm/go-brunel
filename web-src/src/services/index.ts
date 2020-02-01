import {injectable} from 'inversify';
import {Observable, from} from 'rxjs';
import {switchMap} from 'rxjs/operators';
import jwtDecode from 'jwt-decode';
import moment from 'moment';

export interface User {
	Username: string;
	Email: string;
	Name: string;
	AvatarURL: string;
}

export enum JobState {
	Waiting = 0,
	Processing = 1,
	Failed = 2,
	Success = 3,
	Cancelled = 4
}

export interface RepositoryJobPage {
	Count: number;
	Jobs: RepositoryJob[];
}

export interface RepositoryJob {
	ID: string;
	RepositoryID: string;
	Commit: {
		Branch: string;
		Revision: string;
	};
	State: JobState;
	StartedBy: string;
	CreatedAt: string;
	StartedAt: string;
	StoppedAt: string;
}

export interface Repository {
	ID: string;
	Project: string;
	Name: string;
	URI: string;
	CreatedAt: string;
}

@injectable()
export class RepositoryService {
	list(filter = ''): Observable<Repository[]> {
		const query = new URLSearchParams({
			filter: filter.toString(),
		});

		return from(fetch(`/api/repository?${query}`)).pipe(
			switchMap((response) => response.json()),
		);
	}

	jobs(
		id: string,
		filter: string,
		pageIndex: number,
		pageSize: number,
		sortOrder: 'asc' | 'desc',
		sortColumn: string,
	): Observable<RepositoryJobPage> {
		const query = new URLSearchParams({
			sortOrder: sortOrder.toString(),
			sortColumn: sortColumn.toString(),
			pageIndex: pageIndex.toString(),
			pageSize: pageSize.toString(),
			filter: filter.toString(),
		});

		return from(fetch(`/api/repository/${id}/jobs?${query}`)).pipe(
			switchMap((response) => response.json()),
		);
	}
}

@injectable()
export class AuthService {
	isAuthenticated(): boolean {
		const token = window.localStorage.getItem('jwt');
		if (token) {
			try {
				return moment
					.utc((jwtDecode(token) as any).exp)
					.isBefore(moment());
			} catch (e) {
				return false;
			}
		}
		return false;
	}

	setAuthentication(token: string) {
		window.localStorage.setItem('jwt', token);
	}

	getAuthHeaders(): HeadersInit {
		return {
			Authorization: `Bearer ${window.localStorage.getItem('jwt')}`,
		} as HeadersInit;
	}
}

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
