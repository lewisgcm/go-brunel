import React, {useEffect, useState} from 'react';
import {LinearProgress} from '@material-ui/core';
import {delay} from 'rxjs/operators';

import {withDependency} from '../../../container';
import {ContainerState, JobService} from '../../../services';
import {createStyles, makeStyles, Theme} from "@material-ui/core/styles";

interface Dependencies {
    jobService: JobService;
}

interface Props {
    containerId: string;
    containerState: ContainerState;
}

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		'hidden': {
			visibility: 'hidden',
		},
	}),
);

require('./JobContainerLogs.css');

export const JobContainerLogs = React.memo<Props>(
	withDependency<Props, Dependencies>((container) => ({
		jobService: container.get(JobService),
	}))(({jobService, containerId}) => {
		const classes = useStyles();
		const [content, setContent] = useState<null | HTMLDivElement>();
		const [isLoading, setIsLoading] = useState(false);

		useEffect(() => {
			if (content) {
				setIsLoading(true);

				const subscription = jobService
					.containerLogs(containerId)
					.pipe(
						delay(200),
					)
					.subscribe(
						(progress) => {
							content.innerHTML = progress;
							setIsLoading(false);
						},
					);

				return () => {
					return subscription.unsubscribe();
				};
			}
		}, [jobService, containerId, content]);

		return <React.Fragment>
			<LinearProgress className={isLoading ? '' : classes.hidden}/>
			<div className={'term-container'} ref={(r) => setContent(r)} />
		</React.Fragment>;
	}),
	(prevProps, nextProps) =>
		(prevProps.containerId === nextProps.containerId) && prevProps.containerState > ContainerState.Running,
);

