import React, {useEffect, useState} from 'react';
import {LinearProgress} from '@material-ui/core';
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';
import {delay} from 'rxjs/operators';

import {useDependency} from '../../../container';
import {ContainerState, JobService} from '../../../services';

interface Props {
	containerId: string;
	containerState: ContainerState;
}

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		'hidden': {
			visibility: 'hidden',
		},
		'failed': {
			borderLeft: '10px solid red',
		},
	}),
);

require('./JobContainerLogs.css');

export const JobContainerLogs = React.memo<Props>((({containerId, containerState}) => {
	const classes = useStyles();
	const jobService = useDependency(JobService);
	const [content, setContent] = useState<null | HTMLDivElement>();
	const [isLoading, setIsLoading] = useState(false);

	useEffect(
		() => {
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
						() => {
							setIsLoading(false);
						},
					);

				return () => {
					return subscription.unsubscribe();
				};
			}
		},
		[jobService, containerId, content],
	);

	return <React.Fragment>
		<LinearProgress className={isLoading ? '' : classes.hidden}/>
		<div className={'term-container ' + (containerState === ContainerState.Error ? classes.failed : '')}
			ref={(ref) => setContent(ref)} />
	</React.Fragment>;
}),
(prevProps, nextProps) =>
	(prevProps.containerId === nextProps.containerId) && prevProps.containerState > ContainerState.Running,
);

