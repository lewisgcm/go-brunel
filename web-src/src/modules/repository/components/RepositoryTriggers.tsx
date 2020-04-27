import React, {useState, useEffect} from 'react';
import {
	Theme,
	makeStyles,
	createStyles,
	Grid,
	IconButton,
	Button,
	Divider,
	ExpansionPanel,
	ExpansionPanelActions,
	ExpansionPanelDetails,
	ExpansionPanelSummary,
} from '@material-ui/core';
import AddIcon from '@material-ui/icons/Add';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';

import {TriggerEntry} from './TriggerEntry';
import {RepositoryTrigger, RepositoryTriggerType, RepositoryService} from '../../../services';
import {useDependency} from '../../../container';

const useStyles = makeStyles((theme: Theme) => createStyles({
	'triggers': {
		marginBottom: theme.spacing(2),
	},
}));

interface Props {
	id: string;
	triggers?: RepositoryTrigger[];
}

export function RepositoryTriggers(props: Props) {
	const classes = useStyles({});
	const repositoryService = useDependency(RepositoryService);
	const [triggers, setTriggers] = useState<RepositoryTrigger[]>(props.triggers || []);

	useEffect(() => {
		setTriggers(props.triggers || []);
	}, [props]);

	const onSave = (triggers: RepositoryTrigger[]) => {
		repositoryService
			.setTriggers(props.id, triggers)
			.subscribe(() => { });
	};

	return <ExpansionPanel className={classes.triggers}>
		<ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
			<p style={{margin: 0}}>Build Triggers</p>
		</ExpansionPanelSummary>
		<ExpansionPanelDetails>
			<Grid container justify="space-between" spacing={3}>
				<Grid item xs={1}>
					<IconButton
						onClick={() => {
							setTriggers(
								triggers.concat([{
									Type: RepositoryTriggerType.Branch,
									Pattern: '',
								}]),
							);
						}}>
						<AddIcon/>
					</IconButton>
				</Grid>
				<Grid container item xs={11} spacing={3}>
					{
						triggers.map((trigger, index) => {
							return <TriggerEntry
								key={index}
								trigger={trigger}
								onRemove={() => {
									const copy = triggers.slice();
									copy.splice(index, 1);
									setTriggers(copy);
								}}
								onChange={(newTrigger) => {
									const copy = triggers.slice();
									copy[index] = newTrigger;
									setTriggers(copy);
								}}/>;
						})
					}
				</Grid>
			</Grid>
		</ExpansionPanelDetails>
		<Divider />
		<ExpansionPanelActions>
			<Button size="small" color="primary" onClick={() => onSave(triggers)}>
				Save
			</Button>
		</ExpansionPanelActions>
	</ExpansionPanel>;
}
