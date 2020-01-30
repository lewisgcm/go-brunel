import { Component, OnInit, OnDestroy, AfterViewInit } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';

import * as _ from 'lodash';
import * as moment from 'moment';
import { parse } from 'ansicolor';
import { Subscription } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { map, takeWhile } from 'rxjs/operators';

import { JobService } from '../../service/job.service';
import {
    Job,
    JobStage,
    JobContainer,
    JobProgress,
    JobState,
    StageState
} from 'src/app/module/shared';

const LAST_STAGE_ID = 'clean';

@Component({
    selector: 'app-job-page',
    templateUrl: './job.page.html',
    styleUrls: [ './job.page.scss' ],
})
export class JobPageComponent implements OnInit, OnDestroy {

    stageSpacing = 100;
    job?: Job;
    progress: Subscription;

    stages: JobStage[] = [];
    stagesExpansion: {
        [key: string]: boolean
    } = {};

    userWantsToSeeBottom = true;

    constructor(
        private jobService: JobService,
        private activateRoute: ActivatedRoute,
        private sanitizer: DomSanitizer
    ) {
    }

    onScroll = (event: Event): void => {
        this.userWantsToSeeBottom = Math.round(window.scrollY + window.innerHeight) >= Math.round(document.body.scrollHeight);
    }

    public ngOnInit() {
        window.addEventListener('scroll', this.onScroll, true);

        this.activateRoute
            .params
            .pipe(
                map(p => p.id)
            )
            .subscribe(
                (id) => {
                    this.jobService
                        .get(id)
                        .subscribe(
                            (job) => {
                                this.job = job;
                                this.progress = this.jobService
                                    .progress(id)
                                    .pipe(
                                        takeWhile(
                                            (progress: JobProgress) => {
                                                const finalStage = this.stages
                                                    .find(s => s.ID === LAST_STAGE_ID);

                                                return !(finalStage && finalStage.State > StageState.Running);
                                            }
                                        )
                                    )
                                    .subscribe(
                                        (progress) => {
                                            this.stages = progress.Stages;

                                            if (this.stages.length > 0) {
                                                const lastStage = this.stages[this.stages.length - 1];

                                                if (lastStage.State === StageState.Running) {
                                                    this.stagesExpansion[lastStage.ID] = true;
                                                }
                                            }

                                            if (this.userWantsToSeeBottom) {
                                                window.scrollTo(0, document.body.scrollHeight);
                                            }
                                        }
                                    );
                            }
                        );
                },
            );
    }

    public getStyle(css: string) {
        return this.sanitizer.bypassSecurityTrustStyle(css);
    }

    public ngOnDestroy() {
        window.removeEventListener('scroll', this.onScroll, true);

        this.progress.unsubscribe();
    }

    public parse(log: string) {
        return parse(log);
    }

    public stagesCount() {
        return Object.keys(this.stages).length;
    }

    public getContainerTitle(container: JobContainer) {
        let title = `container '${container.Spec.Image}'`;

        if (container.Spec.EntryPoint) {
            title += `: ${container.Spec.EntryPoint}`;
        }

        if (container.Spec.Args) {
            title += ` ${container.Spec.Args.join(' ')}`;
        }

        return title;
    }

    public onExpand(stageID: string) {
        this.stagesExpansion[stageID] = !this.stagesExpansion[stageID];
    }

    public duration(stage: JobStage) {
        if (stage.StartedAt && stage.StoppedAt) {
            return moment.duration(moment(stage.StoppedAt).diff(moment(stage.StartedAt))).humanize();
        }
        return '';
    }
}
