import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { AuthService, Repository, JobPageable } from '../../shared';
import { Observable } from 'rxjs';

@Injectable()
export class RepositoryService {
    constructor(private http: HttpClient, private authService: AuthService) { }

    public getRepositories(filter: string = ''): Observable<Repository[]> {
        return this.http.get<Repository[]>(
            encodeURI(`/api/repository/?filter=${filter}`),
            {
                headers: this.authService.authHeaders(),
            }
        );
    }

    public get(id: string): Observable<Repository> {
        return this.http.get<Repository>(
            `/api/repository/${id}`,
            {
                headers: this.authService.authHeaders(),
            }
        );
    }

    public getJobs(
        id: string,
        filter: string,
        pageIndex: number,
        pageSize: number,
        sortColumn: string,
        sortOrder: string,
    ): Observable<JobPageable> {
        return this.http.get<JobPageable>(
            `/api/repository/${id}/jobs`,
            {
                headers: this.authService.authHeaders(),
                params: {
                    filter,
                    pageIndex: pageIndex.toString(),
                    pageSize: pageSize.toString(),
                    sortColumn,
                    sortOrder,
                }
            }
        );
    }
}
