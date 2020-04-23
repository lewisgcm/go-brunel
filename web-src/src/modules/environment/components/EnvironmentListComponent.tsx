import React from 'react';
import {LinearProgress, List, ListItem, ListItemText, TextField, Typography} from '@material-ui/core';
import {Link} from 'react-router-dom';
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';
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
	}),
);

interface Props {
	isLoading: boolean;
	environments: EnvironmentList[];
	onClick: (repository: string) => void;
	onSearch: (term: string) => void;
	selectedEnvironmentId?: string;
}

export function EnvironmentListComponent({
	isLoading,
	environments,
	selectedEnvironmentId,
	onClick,
	onSearch,
}: Props) {
	const classes = useStyles();

	return <List className={classes.list}>
		<TextField className={classes.input}
			label="Search for an environment"
			onChange={(e) => onSearch(e.target.value)} />
		<LinearProgress className={isLoading ? '' : classes.hidden} />
		{environments.map(
			(r) => {
				return <ListItem
					className={`${classes.listItem}`}
					button
					component={Link}
					key={r.ID}
					to={`/repository/${r.ID}`}
					onClick={() => onClick(r.ID)} >
					<ListItemText>{r.Name}/{r.Name}</ListItemText>
				</ListItem>;
			},
		)}
		{environments.length === 0 && <Typography className={classes.empty}>
			No environments found.
		</Typography>}
	</List>;
}
