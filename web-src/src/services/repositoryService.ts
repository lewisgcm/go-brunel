import { injectable } from "inversify";
import { Observable, from, throwError } from "rxjs";
import { switchMap } from "rxjs/operators";

import { AuthService } from "./authService";
import { RepositoryTrigger, Repository, RepositoryJobPage } from "./models";
import { handleResponse } from "./util";

@injectable()
export class RepositoryService {
	constructor(private _authService: AuthService) {}

	setTriggers(id: string, triggers: RepositoryTrigger[]): Observable<{}> {
		return from(
			fetch(`/api/repository/${id}/triggers`, {
				method: "PUT",
				headers: this._authService.getAuthHeaders(),
				body: JSON.stringify(triggers),
			})
		).pipe(
			switchMap((response) =>
				response.ok
					? response.text()
					: throwError(new Error(response.statusText))
			)
		);
	}

	list(filter = ""): Observable<Repository[]> {
		const query = new URLSearchParams({
			filter: filter.toString(),
		});

		return from(
			fetch(`/api/repository?${query}`, {
				method: "GET",
				headers: this._authService.getAuthHeaders(),
			})
		).pipe(switchMap(handleResponse));
	}

	get(id: string): Observable<Repository> {
		return from(
			fetch(`/api/repository/${id}`, {
				method: "GET",
				headers: this._authService.getAuthHeaders(),
			})
		).pipe(switchMap(handleResponse));
	}

	jobs(
		id: string,
		filter: string,
		pageIndex: number,
		pageSize: number,
		sortOrder: "asc" | "desc",
		sortColumn: string
	): Observable<RepositoryJobPage> {
		const query = new URLSearchParams({
			sortOrder: sortOrder.toString(),
			sortColumn: sortColumn.toString(),
			pageIndex: pageIndex.toString(),
			pageSize: pageSize.toString(),
			filter: filter.toString(),
		});

		return from(
			fetch(`/api/repository/${id}/jobs?${query}`, {
				method: "GET",
				headers: this._authService.getAuthHeaders(),
			})
		).pipe(switchMap(handleResponse));
	}
}
