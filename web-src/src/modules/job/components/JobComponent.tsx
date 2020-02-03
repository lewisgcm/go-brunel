import React, {useEffect, useState} from 'react';
import {match, useHistory} from 'react-router';
import {AppBar, Button, Toolbar, Typography} from '@material-ui/core';
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';

import {withDependency} from '../../../container';
import {JobProgress, JobService} from '../../../services';
import {JobProgressGraph} from './JobProgressGraph';
import {JobContainerLogs} from './JobContainerLogs';
import {JobStageLogs} from './JobStageLogs';

interface Dependencies {
    jobService: JobService;
}

interface Props {
    match: match<{jobId: string}>;
}

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		appBar: {
			zIndex: theme.zIndex.drawer + 1,
		},
	}),
);

export const JobComponent = withDependency<Props, Dependencies>(
	(container) => ({
		jobService: container.get(JobService),
	}),
)(({jobService, match}) => {
	const history = useHistory();
	const classes = useStyles();
	const {jobId} = match.params;
	const [jobProgress, setJobProgress] = useState<JobProgress>({Stages: []});
	const [selectedStage, setSelectedStage] = useState();
	const stageSelect = (newStageId: string) => {
		setSelectedStage(newStageId);
	};

	useEffect(() => {
		const subscription = jobService
			.progress(jobId)
			.subscribe(
				(progress) => {
					setJobProgress(progress);
					if (!selectedStage && progress.Stages.length) {
						stageSelect(progress.Stages[0].ID);
					}
				},
			);

		return () => {
			return subscription.unsubscribe();
		};
	}, [jobService, jobId, selectedStage]);

	return <div>
		<AppBar className={classes.appBar} elevation={0}>
			<Toolbar>
				<Button color={'inherit'}
					onClick={() => history.goBack()}>
					{'Back'}
				</Button>
			</Toolbar>
		</AppBar>
		<JobProgressGraph stages={jobProgress.Stages}
			onStageSelect={(s) => stageSelect(s.ID)}
			selectedStageId={selectedStage} />
		{jobProgress
			.Stages
			.filter((s) => s.ID === selectedStage)
			.map(
				(s) => <JobStageLogs key={s.ID} stage={s} />,
			)
		}
		{jobProgress.Stages
			.filter((s) => s.ID === selectedStage)
			.flatMap((s) => (s.Containers || []))
			.map((c) => {
				return <React.Fragment key={c.ContainerID}>
					<Typography>
						{c.Spec.Image}
					</Typography>
					<JobContainerLogs containerId={c.ContainerID}
						containerState={c.State} />
				</React.Fragment>;
			})}
	</div>;
});

