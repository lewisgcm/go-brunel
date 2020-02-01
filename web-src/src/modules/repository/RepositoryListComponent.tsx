import React from 'react';
import {
	makeStyles,
	Theme,
	createStyles,
} from '@material-ui/core/styles';

import {
	List,
	ListItem,
	ListItemText,
	Typography,
	TextField,
	LinearProgress,
} from '@material-ui/core';

import {Link} from 'react-router-dom';

import {Repository} from '../../services';

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
	}),
);

interface Props {
	isLoading: boolean;
	repositories: Repository[];
	onClick: (repository: Repository) => void;
	onSearch: (term: string) => void;
}

export function RepositoryListComponent({
	isLoading,
	repositories,
	onClick,
	onSearch,
}: Props) {
	const classes = useStyles();

	return <List className={classes.list}>
		<TextField className={classes.input}
			label="Search for a repository"
			onChange={(e) => onSearch(e.target.value)} />
		<LinearProgress className={isLoading ? '' : classes.hidden} />
		{repositories.map(
			(r) => {
				return <ListItem
					className={classes.listItem}
					button
					component={Link}
					key={r.ID}
					to={`/repository/${r.ID}`}
					onClick={() => onClick(r)} >
					<ListItemText>{r.Project}/{r.Name}</ListItemText>
				</ListItem>;
			},
		)}
		{repositories.length === 0 && <Typography className={classes.empty}>
			No repositories found.
		</Typography>}
	</List>;
}