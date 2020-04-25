import React from 'react';
import {LinearProgress, List, ListItem, ListItemText, TextField, Typography, Button} from '@material-ui/core';
import {Link} from 'react-router-dom';
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';
import AddIcon from '@material-ui/icons/Add';

import {EnvironmentList} from '../../../services';

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		list: {
			padding: theme.spacing(2),
		},
		listItem: {
			borderBottom: `1px solid ${theme.palette.grey[300]}`,
		},
		input: {
			width: '100%',
		},
		hidden: {
			visibility: 'hidden',
		},
		empty: {
			textAlign: 'center',
			paddingTop: theme.spacing(1),
		},
		selectedItem: {
			backgroundColor: theme.palette.grey[300],
		},
		addButton: {
			marginBottom: theme.spacing(1),
		},
	}),
);

interface Props {
	isLoading: boolean;
	environments: EnvironmentList[];
	onClick: (repository: string) => void;
	onSearch: (term: string) => void;
	onAdd: () => void;
	selectedEnvironmentId?: string;
}

export function EnvironmentListComponent({
	isLoading,
	environments,
	selectedEnvironmentId,
	onClick,
	onSearch,
	onAdd,
}: Props) {
	const classes = useStyles();

	return <List className={classes.list}>
		<Button
			onClick={() => onAdd()}
			className={classes.addButton}
			variant="contained"
			color="primary"
			fullWidth
			startIcon={<AddIcon />}>
				Add Environment
		</Button>

		<TextField className={classes.input}
			label="Search for an environment"
			onChange={(e) => onSearch(e.target.value)} />

		<LinearProgress className={isLoading ? '' : classes.hidden} />

		{environments.map(
			(r) => {
				return <ListItem
					className={`${classes.listItem} ${selectedEnvironmentId === r.ID ? classes.selectedItem : ''}`}
					button
					component={Link}
					key={r.ID}
					to={`/environment/${r.ID}`}
					onClick={() => onClick(r.ID)} >
					<ListItemText>{r.Name}</ListItemText>
				</ListItem>;
			},
		)}
		{environments.length === 0 && <Typography className={classes.empty}>
			No environments found.
		</Typography>}
	</List>;
}
