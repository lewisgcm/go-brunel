import React, {useEffect, useState} from 'react';
import {withDependency} from '../../../container';
import {JobProgress, JobService} from '../../../services';
import {match} from 'react-router';
import {JobProgressGraph} from './JobProgressGraph';
import {JobContainerLogs} from './JobContainerLogs';
import {Typography} from '@material-ui/core';

interface Dependencies {
    jobService: JobService;
}

interface Props {
    match: match<{jobId: string}>;
}

export const JobComponent = withDependency<Props, Dependencies>((container) => ({
	jobService: container.get(JobService),
}))(({jobService, match}) => {
	const {jobId} = match.params;
	const [jobProgress, setJobProgress] = useState<JobProgress>({Stages: []});
	const [selectedStage, setSelectedStage] = useState<string | undefined>();

	useEffect(() => {
		const subscription = jobService
			.progress(jobId)
			.subscribe(
				(progress) => {
					setJobProgress(progress);
					if (!selectedStage && progress.Stages.length) {
						setSelectedStage(progress.Stages[0].ID);
					}
				},
			);

		return () => {
			return subscription.unsubscribe();
		};
	}, [jobService, jobId, selectedStage]);

	return <div>
		<JobProgressGraph stages={jobProgress.Stages}
			onStageSelect={(s) => setSelectedStage(s.ID)}
			selectedStageId={selectedStage} />
		{jobProgress.Stages.filter((s) => s.ID === selectedStage).map(
			(s) => {
				return s.Logs.length > 0 ? <div key={s.ID} className={'term-container'} >
					{s.Logs.map((l, i)=> <React.Fragment key={`${l.StageID}-${i}`}>{l.Message} <br/></React.Fragment>)}
				</div> : <React.Fragment key={s.ID}></React.Fragment>;
			},
		)}
		{jobProgress.Stages
			.filter((s) => s.ID === selectedStage)
			.flatMap((s) => (s.Containers || []))
			.map((c) => {
				return <React.Fragment key={c.ContainerID}>
					<Typography>
						{c.Spec.Image}
					</Typography>
					<JobContainerLogs containerId={c.ContainerID} containerState={c.State}/>
				</React.Fragment>;
			})}
	</div>;
});

