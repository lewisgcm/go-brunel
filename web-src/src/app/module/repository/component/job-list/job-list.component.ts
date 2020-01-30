import { Component, Input, Output, EventEmitter, ViewChild, OnInit, OnDestroy } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';

import * as _ from 'lodash';
import * as moment from 'moment';
import { merge, Subject, BehaviorSubject, timer, Subscription } from 'rxjs';
import { debounceTime, switchMap, tap } from 'rxjs/operators';

import { Job, JobState, JobPageable } from '../../../shared';
import { RepositoryService } from '../../service';

const REFRESH_RATE_MS = 2000;

@Component({
    selector: 'app-repository-job-list',
    templateUrl: './job-list.component.html',
    styleUrls: ['./job-list.component.scss'],
})
export class JobListComponent implements OnInit, OnDestroy {
    id: string;
    data?: JobPageable;
    loading: boolean;
    filter: string;
    filterSubject: Subject<string>;
    idSubject: Subject<string>;
    subscription: Subscription;
    displayedColumns: string[] = [
        'state',
        'branch',
        'duration',
        'created_at',
        'started_by',
    ];

    get repositoryID(): string {
        return this.id;
    }

    @Input()
    set repositoryID(id: string) {
        this.idSubject.next(id);
    }

    @Output() jobSelect: EventEmitter<string> = new EventEmitter();

    @ViewChild(MatPaginator, { static: true }) paginator: MatPaginator;
    @ViewChild(MatSort, { static: true }) sort: MatSort;

    constructor(private repositoryService: RepositoryService) {
        this.filter = '';
        this.filterSubject = new BehaviorSubject<string>('');
        this.idSubject = new BehaviorSubject<string>(this.id);
        this.loading = true;
        this.data = { Jobs: [], Count: 0 };
    }

    public ngOnInit() {
        // Handle our search filtering, and reset our page index when sorting
        this.filterSubject.subscribe((f) => this.filter = f);
        this.sort.sortChange.subscribe(() => this.paginator.pageIndex = 0);
        this.idSubject.subscribe((d) => {
            this.id = d;
            this.paginator.pageIndex = 0;
        });

        // Handle any change to the table such as sorting, pagination or searching by reloading our table data
        this.subscription = merge(
            this.sort.sortChange,
            this.paginator.page,
            this.filterSubject,
            this.idSubject,
            timer(REFRESH_RATE_MS, REFRESH_RATE_MS)
        )
            .pipe(
                tap(() => {
                    this.loading = true;
                }),
                debounceTime(300),
                switchMap(
                    () => this.repositoryService.getJobs(
                        this.repositoryID,
                        this.filter,
                        this.paginator.pageIndex,
                        this.paginator.pageSize,
                        this.sort.active,
                        this.sort.direction,
                    )
                ),
                tap(() => this.loading = false),
            )
            .subscribe(
                (jobs) => {
                    this.data.Count = jobs.Count;
                    this.data.Jobs = _.map(jobs.Jobs, (j) => {
                        return {
                            ...j,
                            Duration: j.StartedAt && j.StoppedAt ?
                                moment.duration(moment(j.StoppedAt).diff(moment(j.StartedAt))).humanize() : 'n/a',
                            CreatedAt: moment.duration(moment(j.CreatedAt).diff(moment())).humanize() + ' ago',
                            Commit: {
                                ...j.Commit,
                                Branch: j.Commit.Branch.replace(/refs\/[^\/]+\//, ''),
                            },
                        };
                    });
                }
            );
    }

    public onClick(e: Job) {
        this.jobSelect.emit(e.ID);
    }

    public ngOnDestroy() {
        this.idSubject.complete();
        this.filterSubject.complete();
        this.subscription.unsubscribe();
    }

    public state(job: Job) {
        return JobState[job.State];
    }

    public stateIcon(job: Job): string {
        switch (job.State) {
            case JobState.Success:
                return 'check_circle';
            case JobState.Failed:
                return 'warning';
            case JobState.Cancelled:
                return 'remove_circle';
        }
        return 'sync';
    }

    public stateColor(job: Job): string {
        return `job-state ${JobState[job.State].toLowerCase()}`;
    }
}
