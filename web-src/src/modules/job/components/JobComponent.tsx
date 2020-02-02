import React, {useEffect, useState} from 'react';
import {withDependency} from "../../../container";
import {JobProgress, JobService} from "../../../services";
import {match} from "react-router";
import {JobProgressGraph} from "./JobProgressGraph";
import {JobContainerLogs} from "./JobContainerLogs";
import {Typography} from "@material-ui/core";

interface Dependencies {
    jobService: JobService;
}

interface Props {
    match: match<{jobId: string}>;
}

export const JobComponent = withDependency<Props, Dependencies>((container) => ({
    jobService: container.get(JobService),
}))(({jobService, match}) => {
    const { jobId } = match.params;
    const [jobProgress, setJobProgress] = useState<JobProgress>({Stages: []});

    useEffect(() => {
        const subscription = jobService
            .progress(jobId)
            .subscribe(
                (progress) => {
                    setJobProgress(progress);
                }
            );

        return () => {
            return subscription.unsubscribe();
        };
    }, [jobService, jobId]);

    return <div>
        <JobProgressGraph stages={jobProgress.Stages}
                          onStageSelect={() => {}}
                          selectedStageId={''} />
        {jobProgress.Stages.flatMap(s => (s.Containers || [])).map(c => {
            return <React.Fragment key={c.ContainerID}>
                <Typography>
                    {c.Spec.Image}
                </Typography>
                <JobContainerLogs containerId={c.ContainerID}/>
            </React.Fragment>
        })}
    </div>;
});

