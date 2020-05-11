import React from 'react';
import {Link} from 'react-router-dom';
import {makeStyles, Theme, createStyles} from '@material-ui/core/styles';
import {
	List,
	ListItem,
	ListItemText,
	Typography,
	TextField,
	LinearProgress,
} from '@material-ui/core';

import {UserList} from '../../../services';

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
	users: UserList[];
	onClick: (user: UserList) => void;
	onSearch: (term: string) => void;
	selectedUsername?: string;
}

export function UserListComponent({
	isLoading,
	users,
	selectedUsername,
	onClick,
	onSearch,
}: Props) {
	const classes = useStyles({});

	return <List className={classes.list}>
		<TextField className={classes.input}
			label="Search for a user"
			onChange={(e) => onSearch(e.target.value)} />
		<LinearProgress className={isLoading ? '' : classes.hidden} />
		{users.map(
			(user) => {
				return <ListItem
					className={`${classes.listItem} ${selectedUsername === user.Username ? classes.selectedItem : ''}`}
					button
					component={Link}
					key={user.Username}
					to={`/user/${user.Username}`}
					onClick={() => onClick(user)} >
					<ListItemText>{user.Username}</ListItemText>
				</ListItem>;
			},
		)}
		{users.length === 0 && <Typography className={classes.empty}>
			No users found.
		</Typography>}
	</List>;
}
