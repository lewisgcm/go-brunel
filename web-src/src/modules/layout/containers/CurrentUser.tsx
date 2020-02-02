import React, {Fragment, useEffect, useState} from 'react';
import {connect} from 'react-redux';
import {
	Avatar,
	Menu,
	MenuItem,
	Typography,
	IconButton,
} from '@material-ui/core';
import MoreIcon from '@material-ui/icons/MoreVert';
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';

import {withDependency} from '../../../container';
import {AuthService, UserService} from '../../../services';
import {setAuthenticated} from '../actions';

interface Dependencies {
	authService: AuthService;
	userService: UserService;
}

interface Props {
	setLoggedOut: () => void;
}

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		username: {
			paddingLeft: theme.spacing(2),
			fontWeight: 'bold',
			fontSize: '1.1em',
		},
	}),
);

export const CurrentUser = connect(
	null,
	((dispatch) => ({
		setLoggedOut: () => dispatch(setAuthenticated(false)),
	})),
)(
	withDependency<Props, Dependencies>((container) => ({
		authService: container.get(AuthService),
		userService: container.get(UserService),
	}))(
		({authService, userService, setLoggedOut}) => {
			const classes = useStyles();
			const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
			const [username, setUsername] = useState('');
			const [avatarUrl, setAvatarUrl] = useState('');

			const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
				setAnchorEl(event.currentTarget);
			};

			const handleClose = () => {
				setAnchorEl(null);
			};

			const onLogout = () =>{
				authService.setAuthentication('');
				handleClose();
				setLoggedOut();
			};

			useEffect(() => {
				userService
					.getProfile()
					.subscribe(
						(user) => {
							setUsername(user.Username);
							setAvatarUrl(user.AvatarURL);
						},
					);
			}, [userService]);

			return <Fragment>
				{avatarUrl && username && <React.Fragment>
					<Avatar src={avatarUrl} />
					<Typography className={classes.username}>
						{username}
					</Typography>
				</React.Fragment>}
				<IconButton edge="end" color="inherit" onClick={handleClick}>
					<MoreIcon />
				</IconButton>
				<Menu
					id="simple-menu"
					anchorEl={anchorEl}
					keepMounted
					open={Boolean(anchorEl)}
					onClose={handleClose}
				>
					<MenuItem onClick={onLogout}>Logout</MenuItem>
				</Menu>
			</Fragment>;
		}),
);
