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

import {Trigger} from './Trigger';
import {
	RepositoryTrigger,
	RepositoryTriggerType,
	RepositoryService,
} from '../../../services';
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

interface RepositoryTriggerItem extends RepositoryTrigger {
	isValid?: boolean;
}

export function RepositoryTriggers(props: Props) {
	const classes = useStyles({});
	const repositoryService = useDependency(RepositoryService);
	const [triggers, setTriggers] = useState<RepositoryTriggerItem[]>(props.triggers || []);

	useEffect(() => {
		setTriggers(props.triggers || []);
	}, [props]);

	const isTriggerValid = (trigger: RepositoryTrigger) => {
		return trigger && trigger.Pattern.trim().length > 0 && trigger.Pattern.trim().length < 100;
	};

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
									isValid: false,
								}]),
							);
						}}>
						<AddIcon/>
					</IconButton>
				</Grid>
				<Grid container item xs={11} spacing={3}>
					{
						triggers.map((trigger, index) => {
							return <Trigger
								key={index}
								trigger={trigger}
								onRemove={() => {
									const copy = triggers.slice();
									copy.splice(index, 1);
									setTriggers(copy);
								}}
								isValid={trigger.isValid === undefined ? true : trigger.isValid}
								onChange={(newTrigger: RepositoryTriggerItem) => {
									newTrigger.isValid = isTriggerValid(newTrigger);
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
			<Button disabled={triggers.some((t) => t.isValid === false)} size="small" color="primary" onClick={() => onSave(triggers)}>
				Save
			</Button>
		</ExpansionPanelActions>
	</ExpansionPanel>;
}
