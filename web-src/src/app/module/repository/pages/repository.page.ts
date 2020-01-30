import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Location } from '@angular/common';
import { Subject } from 'rxjs';
import { map, filter, debounceTime, tap } from 'rxjs/operators';

import { Repository } from '../../shared';
import { RepositoryService } from '../service';

@Component({
    selector: 'app-repository-page',
    templateUrl: './repository.page.html',
    styleUrls: ['./repository.page.scss'],
})
export class RepositoryPageComponent implements OnInit, OnDestroy {
    repositories: Repository[];
    selectedRepository?: Repository;
    searchLoading: boolean;
    searchSubject: Subject<string>;

    constructor(
        private router: Router,
        private location: Location,
        private activatedRoute: ActivatedRoute,
        private repositoryService: RepositoryService,
    ) {
        this.searchSubject = new Subject<string>();
        this.searchLoading = false;
        activatedRoute.
            params.
            pipe(
                map(p => p.id),
                filter(i => i),
            ).
            subscribe(
                (id) => {
                    this.onSelect(id);
                }
            );
    }

    public onJobSelect(id: string) {
        this.router.navigate([`/job/${id}`]);
    }

    public onSelect(id: string) {
        this.repositoryService
            .get(id).subscribe(
                (r) => {
                    this.location.go(`/repository/${r.ID}`);
                    this.selectedRepository = r;
                }
            );
    }

    public onSearch(term: string) {
        this.searchSubject.next(term);
    }

    public ngOnDestroy() {
        this.searchSubject.complete();
    }

    public ngOnInit() {
        this.searchSubject
            .pipe(
                tap(() => this.searchLoading = true),
                debounceTime(400)
            )
            .subscribe(
                (term) => {
                    this.repositoryService
                        .getRepositories(term)
                        .subscribe(
                            (repositories) => {
                                this.searchLoading = false;
                                this.repositories = repositories;
                            }
                        );
                }
            );
        this.onSearch('');
    }
}
