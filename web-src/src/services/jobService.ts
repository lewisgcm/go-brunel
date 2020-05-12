import {injectable} from 'inversify';
import {Observable, from, timer} from 'rxjs';
import {map, scan, switchMap, takeWhile} from 'rxjs/operators';

import {AuthService} from './authService';
import {Job, JobProgress, JobState} from './models';
import {handleResponse} from './util';

const POLL_INTERVAL_MS = 2 * 1000;

@injectable()
export class JobService {
	constructor(private _authService: AuthService) {
	}

	public get(id: string): Observable<Job> {
		return from(fetch(
			`/api/job/${id}`,
			{headers: this._authService.getAuthHeaders()},
		)).pipe(
			switchMap(handleResponse),
		);
	}

	public cancel(id: string): Observable<{}> {
		return from(fetch(
			`/api/job/${id}`,
			{headers: this._authService.getAuthHeaders(), method: 'DELETE'},
		)).pipe(
			switchMap((response) => response.text() as Promise<{}>),
		);
	}

	public reSchedule(id: string): Observable<Job> {
		return from(fetch(
			`/api/job/${id}/reschedule`,
			{headers: this._authService.getAuthHeaders(), method: 'POST'},
		)).pipe(
			switchMap(handleResponse),
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
					},
				)).pipe(
					switchMap(handleResponse),
				),
			),
			scan(
				(accumulated: JobProgress, current: JobProgress) => {
					if (accumulated === null) {
						return current;
					}

					if (!current.Stages) {
						current.Stages = [];
					}

					current.Stages.forEach(
						(stage) => {
							const oldStage = accumulated.Stages
								.find((s) => s.ID === stage.ID);

							stage.Logs = stage.Logs || [];

							if (oldStage) {
								oldStage.Logs = oldStage.Logs || [];

								stage.Logs = oldStage.Logs.concat(stage.Logs);
							}
						},
					);

					return current;
				},
				{
					State: JobState.Waiting,
					Stages: [],
				},
			),
			takeWhile(
				(progress: JobProgress) => progress.State === JobState.Waiting || progress.State === JobState.Processing,
				true,
			),
		);
	}
}
