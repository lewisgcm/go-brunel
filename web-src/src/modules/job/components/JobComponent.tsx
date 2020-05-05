import React, {useEffect, useState} from 'react';
import {match, useHistory} from 'react-router';
import {AppBar, Button, Toolbar, Tooltip, Typography, withStyles} from '@material-ui/core';
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';
import {red} from '@material-ui/core/colors';
import {FaCodeBranch, FaUserPlus, FaUserTimes, FaRegClock} from 'react-icons/fa';

import {withDependency} from '../../../container';
import {Job, JobProgress, JobService, JobState} from '../../../services';
import {JobProgressGraph} from './JobProgressGraph';
import {JobContainerLogs} from './JobContainerLogs';
import {JobStageLogs} from './JobStageLogs';
import moment from 'moment';

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
			paddingLeft: theme.spacing(2),
			paddingRight: theme.spacing(2),
		},
		grow: {
			flexGrow: 1,
		},
		title: {
			fontWeight: 'bold',
			paddingLeft: theme.spacing(2),
		},
		titleJobInfo: {
			'paddingLeft': theme.spacing(2),
			'fontSize': theme.typography.body2.fontSize,
			'& svg': {
				verticalAlign: 'middle',
				height: '1.3em',
				width: '1.3em',
				marginRight: '8px',
			},
		},
	}),
);

const CancelButton = withStyles((theme: Theme) => ({
	root: {
		'color': theme.palette.getContrastText(red[500]),
		'backgroundColor': red[700],
		'marginLeft': theme.spacing(2),
		'&:hover': {
			backgroundColor: red[900],
		},
	},
}))(Button);

export const JobComponent = withDependency<Props, Dependencies>(
	(container) => ({
		jobService: container.get(JobService),
	}),
)(({jobService, match}) => {
	const history = useHistory();
	const classes = useStyles({});
	const {jobId} = match.params;
	const [job, setJob] = useState<Job | undefined>();
	const [jobProgress, setJobProgress] = useState<JobProgress>({State: JobState.Waiting, Stages: []});
	const [selectedStage, setSelectedStage] = useState<string | undefined>();

	const stageSelect = (newStageId: string) => {
		setSelectedStage(newStageId);
	};

	const onCancel = () => {
		jobService.cancel(jobId).subscribe(
			() => {},
		);
	};

	useEffect(() => {
		jobService
			.get(jobId)
			.subscribe(
				(job) => {
					setJob(job);
				},
			);

		const subscription = jobService
			.progress(jobId)
			.subscribe(
				(progress) => {
					setJobProgress(progress);
				},
			);

		return () => {
			return subscription.unsubscribe();
		};
	}, [jobService, jobId]);

	useEffect(() => {
		if (!selectedStage && jobProgress && jobProgress.Stages.length) {
			setSelectedStage(jobProgress.Stages[0].ID);
		}
	}, [jobProgress, selectedStage]);

	return <div>
		<AppBar className={classes.appBar} elevation={0}>
			<Toolbar disableGutters={true}>
				<Button color='inherit'
					onClick={() => history.goBack()}>
					Back
				</Button>
				{job && <React.Fragment>
					<Typography className={classes.title}>
						{job.Repository.Project}/{job.Repository.Name}
					</Typography>
				</React.Fragment>}
				<span className={classes.grow}/>
				{
					job && <React.Fragment>
						<Tooltip title={job.Commit.Revision}>
							<Typography className={classes.titleJobInfo}>
								<FaCodeBranch/>
								{job.Commit.Branch.replace('refs/heads/', '')}
							</Typography>
						</Tooltip>

						<Tooltip title={`Created by ${job.StartedBy}`} >
							<Typography className={classes.titleJobInfo}>
								<FaUserPlus />
								{job.StartedBy}
							</Typography>
						</Tooltip>

						{moment(job.StartedAt).isValid() && <Typography className={classes.titleJobInfo}>
							<FaRegClock />
							{moment(job.StartedAt).format('LLLL')}
						</Typography>}

						{job.StoppedBy && <Tooltip title={`Cancelled by ${job.StoppedBy}`} >
							<Typography className={classes.titleJobInfo}>
								<FaUserTimes />
								{job.StoppedBy}
							</Typography>
						</Tooltip>}
					</React.Fragment>
				}
				{jobProgress.State === JobState.Processing && <CancelButton onClick={() => onCancel()}>
					Cancel
				</CancelButton>}
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

