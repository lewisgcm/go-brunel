import {injectable} from 'inversify';
import {Observable, from, timer} from 'rxjs';
import {map, scan, switchMap, takeWhile} from 'rxjs/operators';
import jwtDecode from 'jwt-decode';
import moment from 'moment';

export interface Commit {
    Branch: string;
    Revision: string;
}

export interface Job {
    ID: string;
    StartedBy: string;
    CreatedAt: string;
    StartedAt: string;
    StoppedAt: string;
    Duration: string | number;
    State: JobState;
    Commit: Commit;
}

export interface Log {
    Message: string;
    Type: number;
    Time: string;
    StageID: string;
}

export interface ContainerLog extends Log {
    ContainerID: string;
}

export interface JobProgress {
    Stages: JobStage[];
}

export interface JobContainer extends Container {
    Logs: ContainerLog[];
}

export interface JobStage extends Stage {
    Containers: JobContainer[];
    Logs: Log[];
}

export interface Stage {
    ID: string;
    JobID: string;
    StartedAt: string;
    State: number;
    StoppedAt: string;
}

export interface Container {
    ID: string;
    JobID: string;
    ContainerID: string;
    State: number;
    Meta: {
        StageID: string;
        Service: boolean;
    };
    Spec: any;
    CreatedAt: string;
    StartedAt: string;
    StoppedAt: string;
}

export enum StageState {
    Running = 0,
    Success = 1,
    Error = 2
}

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

const POLL_INTERVAL_MS = 2 * 1000;
const LAST_STAGE_ID = 'clean';

@injectable()
export class JobService {
    constructor(private _authService: AuthService) { }

    public get(id: string): Observable<Job> {
        return from(fetch(
            `/api/job/${id}`,
            {headers: this._authService.getAuthHeaders()},
        )).pipe(
            switchMap((response) => response.json()),
        );
    }

    public containerLogs(containerId: string): Observable<string> {
        return from(fetch(
            `/api/container/${containerId}/logs`,
            {headers: this._authService.getAuthHeaders()},
        )).pipe(
            switchMap((response) => response.text()),
        );
    }

    public progress(id: string): Observable<JobProgress> {
        return timer(0, POLL_INTERVAL_MS).pipe(
            map(
                (i) => {
                    // Get everything from the start of time on first poll
                    if (i === 0) {
                        return 0;
                    } else {
                        // Otherwise only get things since our last poll
                        return (Date.now() - POLL_INTERVAL_MS);
                    }
                },
            ),
            switchMap(
                (since) => from(fetch(
                    `/api/job/${id}/progress?since=${since}`,
                    {
                        headers: this._authService.getAuthHeaders(),
                    }
                )).pipe(
                    switchMap((response) => response.json()),
                )
            ),
            scan(
                (accumulated: JobProgress, current: JobProgress) => {
                    if (accumulated === null) {
                        return current;
                    }

                    current.Stages.forEach(
                        (stage) => {
                            let oldStage = accumulated.Stages
                                .find(s => s.ID === stage.ID);

                            if (oldStage) {
                                stage.Logs = (oldStage.Logs || []).concat(stage.Logs);

                                (stage.Containers || []).forEach(
                                    (container) => {
                                        if (oldStage) {
                                            let oldContainer = (oldStage.Containers || [])
                                                .find(c => c.ID === container.ID);

                                            if (oldContainer) {
                                                container.Logs = oldContainer.Logs.concat(container.Logs);
                                            }
                                        }
                                    }
                                );
                            }
                        }
                    );

                    return current;
                },
                {
                    Stages: []
                }
            ),
            takeWhile(
                (progress: JobProgress, index) => {
                    const finalStage = progress
                        .Stages
                        .find(s => s.ID === LAST_STAGE_ID);

                    return index === 0 || !(finalStage && finalStage.State > StageState.Running);
                }
            ),
        );
    }
}
