import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, timer } from 'rxjs';

import { AuthService, Job, JobProgress } from '../../shared';
import { switchMap, map, scan } from 'rxjs/operators';
import * as _ from "lodash";

const POLL_INTERVAL_MS = 2 * 1000;

@Injectable()
export class JobService {
    constructor(private httpClient: HttpClient, private authService: AuthService) { }

    public get(id: string): Observable<Job> {
        return this.httpClient.get<Job>(
            `/api/job/${id}`,
            { headers: this.authService.authHeaders() }
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
                (since) => this.httpClient.get<JobProgress>(
                    `/api/job/${id}/progress`,
                    {
                        headers: this.authService.authHeaders(),
                        params: {
                            since: since.toString()
                        },
                    }
                )
            ),
            scan(
                (accumulated, current) => {
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
                                                .find(c => c.ID == container.ID);

                                            if (oldContainer) {
                                                container.Logs = oldContainer.Logs.concat(container.Logs);
                                            }
                                        }
                                    }
                                );
                            }
                        }
                    )

                    return current;
                },
                {
                    Stages: []
                }
            )
        );
    }
}
