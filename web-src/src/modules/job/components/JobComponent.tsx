import React, { useEffect, useState } from "react";
import { match, useHistory } from "react-router";
import {
	AppBar,
	Button,
	Toolbar,
	Typography,
	withStyles,
	Hidden,
} from "@material-ui/core";
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";
import { red, blue } from "@material-ui/core/colors";
import {
	FaCodeBranch,
	FaUserPlus,
	FaUserTimes,
	FaRegClock,
} from "react-icons/fa";
import moment from "moment";

import { useDependency } from "../../../container";
import {
	Job,
	JobProgress,
	JobService,
	JobState,
	UserRole,
} from "../../../services";
import { JobProgressGraph } from "./JobProgressGraph";
import { JobContainerLogs } from "./JobContainerLogs";
import { JobStageLogs } from "./JobStageLogs";
import { useHasRole } from "../../layout";
import { JobInformationPopover } from "./JobInformationPopover";

interface Props {
	match: match<{ jobId: string }>;
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
			fontWeight: "bold",
			paddingLeft: theme.spacing(2),
		},
	})
);

const CancelButton = withStyles((theme: Theme) => ({
	root: {
		color: theme.palette.getContrastText(red[500]),
		backgroundColor: red[700],
		marginLeft: theme.spacing(2),
		"&:hover": {
			backgroundColor: red[900],
		},
	},
}))(Button);

const TriggerButton = withStyles((theme: Theme) => ({
	root: {
		color: theme.palette.getContrastText(blue[500]),
		backgroundColor: blue[700],
		marginLeft: theme.spacing(2),
		"&:hover": {
			backgroundColor: blue[900],
		},
	},
}))(Button);

export const JobComponent = ({ match }: Props) => {
	const history = useHistory();
	const classes = useStyles({});
	const jobService = useDependency(JobService);
	const { jobId } = match.params;
	const [job, setJob] = useState<Job | undefined>();
	const [jobProgress, setJobProgress] = useState<JobProgress>({
		State: JobState.Waiting,
		Stages: [],
	});
	const [selectedStage, setSelectedStage] = useState<string | undefined>();
	const isAdmin = useHasRole(UserRole.Admin);

	const stageSelect = (newStageId: string) => {
		setSelectedStage(newStageId);
	};

	const onCancel = () => {
		jobService.cancel(jobId).subscribe(() => {});
	};

	const onReSchedule = (id: string) => {
		jobService.reSchedule(id).subscribe((newJob) => {
			history.push(`/job/${newJob.ID}`);
		});
	};

	useEffect(() => {
		jobService.get(jobId).subscribe((job) => {
			setJob(job);
		});

		const subscription = jobService
			.progress(jobId)
			.subscribe((progress) => {
				setJobProgress(progress);
			});

		return () => {
			return subscription.unsubscribe();
		};
	}, [jobService, jobId]);

	useEffect(() => {
		if (!selectedStage && jobProgress && jobProgress.Stages.length) {
			setSelectedStage(jobProgress.Stages[0].ID);
		}
	}, [jobProgress, selectedStage]);

	return (
		<div>
			<AppBar className={classes.appBar} elevation={0}>
				<Toolbar disableGutters={true}>
					<Button color="inherit" onClick={() => history.goBack()}>
						Back
					</Button>
					{job && (
						<React.Fragment>
							<Hidden xsDown>
								<Typography className={classes.title}>
									{job.Repository.Project}/
									{job.Repository.Name}
								</Typography>
							</Hidden>
						</React.Fragment>
					)}
					<span className={classes.grow} />
					{job && (
						<React.Fragment>
							<JobInformationPopover
								icon={FaCodeBranch}
								information={job.Commit.Branch.replace(
									"refs/heads/",
									""
								)}
								tooltipText={job.Commit.Revision}
								popover={
									<React.Fragment>
										Branch{" "}
										<b>
											{job.Commit.Branch.replace(
												"refs/heads/",
												""
											)}
										</b>{" "}
										at revision <b>{job.Commit.Revision}</b>
									</React.Fragment>
								}
							/>

							<JobInformationPopover
								icon={FaUserPlus}
								information={job.StartedBy}
								tooltipText={`Created by ${job.StartedBy}`}
								popover={
									<React.Fragment>
										Created by <b>{job.StartedBy}</b>
									</React.Fragment>
								}
							/>

							{moment(job.StartedAt).isValid() && (
								<JobInformationPopover
									icon={FaRegClock}
									information={moment(job.StartedAt).format(
										"LLLL"
									)}
									tooltipText={`Started at ${moment(
										job.StartedAt
									).format("LLLL")}`}
									popover={
										<React.Fragment>
											Started at{" "}
											<b>
												{moment(job.StartedAt).format(
													"LLLL"
												)}
											</b>
										</React.Fragment>
									}
								/>
							)}

							{job.StoppedBy && (
								<JobInformationPopover
									icon={FaUserTimes}
									information={job.StoppedBy}
									tooltipText={`Cancelled by ${job.StoppedBy}`}
									popover={
										<React.Fragment>
											Cancelled by <b>{job.StoppedBy}</b>
										</React.Fragment>
									}
								/>
							)}
						</React.Fragment>
					)}
					{jobProgress.State === JobState.Processing && isAdmin && (
						<CancelButton onClick={() => onCancel()}>
							Cancel
						</CancelButton>
					)}
					{jobProgress.State > JobState.Processing && isAdmin && (
						<TriggerButton onClick={() => onReSchedule(jobId)}>
							Retry
						</TriggerButton>
					)}
				</Toolbar>
			</AppBar>
			<JobProgressGraph
				stages={jobProgress.Stages}
				onStageSelect={(s) => stageSelect(s.ID)}
				selectedStageId={selectedStage}
			/>
			{jobProgress.Stages.filter((s) => s.ID === selectedStage).map(
				(s) => (
					<JobStageLogs key={s.ID} stage={s} />
				)
			)}
			{jobProgress.Stages.filter((s) => s.ID === selectedStage)
				.flatMap((s) => s.Containers || [])
				.map((c) => {
					return (
						<React.Fragment key={c.ContainerID}>
							<Typography>{c.Spec.Image}</Typography>
							<JobContainerLogs
								containerId={c.ContainerID}
								containerState={c.State}
							/>
						</React.Fragment>
					);
				})}
		</div>
	);
};
