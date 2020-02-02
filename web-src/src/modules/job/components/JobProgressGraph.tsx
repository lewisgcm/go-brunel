import React from 'react';
import {JobStage, JobState, StageState} from '../../../services';
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';
import {FaCheck, FaTimes, FaSync} from 'react-icons/fa';

interface StageGraphProps {
    stages: JobStage[];
    selectedStageId?: string;
    onStageSelect: (stage: JobStage) => void;
}

const stageSpacing = 100;

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		'svg': {
			marginBottom: '10px',
		},
		'line': {
			stroke: 'lightslategrey',
			strokeWidth: 5,
		},
		'stageText': {
			fill: 'lightslategrey',
		},
		'default': {
			fill: 'lightslategrey',
		},
		'stage': {
			'color': 'white',
			'fill': '#2e7d2e',
			'&:hover': {
				strokeWidth: '5px',
				stroke: '#4bd24b',
				cursor: 'pointer',
			},
			'&.error': {
				fill: '#e00000',
			},
			'&.in-progress': {
				fill: 'grey',
			},
			'&.error:hover': {
				stroke: '#ff5858',
			},
			'&.in-progress:hover': {
				stroke: 'darkgrey',
			},
		},
		'@keyframes spinRound': {
			from: {
				transform: 'rotate(0deg)',
			},
			to: {
				transform: 'rotate(360deg)',
			},
		},
		'inProgress': {
			animation: '$spinRound 2s linear infinite',
		},
	}),
);

function jobStateClass(state: StageState): string {
	switch (state) {
	case StageState.Running:
		return 'in-progress';
	case StageState.Error:
		return 'error';
	default:
		return '';
	}
}

export const JobProgressGraph = (props: StageGraphProps) => {
	const classes = useStyles();
	return <svg width="100%" height="100" viewBox={`0 0 300 100`} className={classes.svg} >

		<g>
			{/* Render the starting point in our graph. */}
			<g transform={`translate(${-stageSpacing}, 50) rotate(0)`} >
				<line x1={0} y1="0" x2={stageSpacing} y2="0" className={classes.line} />
				<text x="0" y="35" textAnchor="middle" className={classes.stageText}>start</text>
				<circle cx="0" cy="0" r="10" className={classes.default} />
			</g>
			{
				props.stages.map(
					(stage, index) => <g key={stage.ID}
						transform={`translate(${index * stageSpacing}, 50) rotate(0)`}
						onClick={() => props.onStageSelect(stage)} >
						<line x1={0} y1="0" x2={stageSpacing} y2="0" className={classes.line} />
						<text x="0" y="35" textAnchor="middle" className={classes.stageText}>{stage.ID}</text>
						<g className={`${classes.stage} ${jobStateClass(stage.State)}`}>
							<circle cx="0" cy="0" r="20" />
							<g>
								{stage.State === StageState.Running && <g transform={'translate(0, 0)'}>
									<animateTransform
										attributeName="transform"
										type="rotate"
										from="0"
										to="360"
										dur="4s"
										repeatCount="indefinite" />
									<g transform={'translate(-9, -9) scale(1.3, 1.3)'}>
										<FaSync/>
									</g>
								</g>}
								{stage.State === StageState.Error && <g transform={'translate(-9, -9) scale(1.3, 1.3)'}><FaTimes/></g>}
								{stage.State === StageState.Success && <g transform={'translate(-9, -9) scale(1.3, 1.3)'}><FaCheck/></g>}
							</g>
						</g>
					</g>,
				)
			}
			{/* Render the ending point in our graph. */}
			<g transform={`translate(${props.stages.length * stageSpacing},50) rotate(0)`} >
				<text x="0" y="35" textAnchor="middle" className={classes.stageText}>end</text>
				<circle cx="0" cy="0" r="10" className={classes.default} />
			</g>
		</g>
	</svg>;
};
