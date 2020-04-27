import React from 'react';
import {NavLink} from 'react-router-dom';
import AppBar from '@material-ui/core/AppBar';
import {Button, IconButton, Toolbar, Hidden} from '@material-ui/core';
import MenuIcon from '@material-ui/icons/Menu';
import CssBaseline from '@material-ui/core/CssBaseline';
import Typography from '@material-ui/core/Typography';
import {
	makeStyles,
	Theme,
	createStyles,
} from '@material-ui/core/styles';

import {CurrentUser} from '../containers/CurrentUser';

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		root: {
			display: 'flex',
		},
		appBar: {
			width: `100%`,
			zIndex: theme.zIndex.drawer + 1,
		},
		buttonActive: {
			'&.is-active': {
				backgroundColor: theme.palette.primary.light,
			},
		},
		menuButton: {
			marginRight: theme.spacing(2),
		},
		drawerPaper: {
			width: '100%',
		},
		drawerHeader: {
			display: 'flex',
			alignItems: 'center',
			padding: theme.spacing(0, 1),
			...theme.mixins.toolbar,
			justifyContent: 'flex-end',
		},
		toolbar: theme.mixins.toolbar,
		content: {
			flexGrow: 1,
			padding: theme.spacing(3),
			paddingTop: 0,
		},
		grow: {
			flexGrow: 1,
		},
		title: {
			marginRight: theme.spacing(1),
		},
	}),
);

interface ResponsiveDrawerProps {
	children: React.ReactNode;
	onSidebarToggle: () => void;
}

export function AuthenticatedLayout(props: ResponsiveDrawerProps) {
	const {children} = props;
	const classes = useStyles({});

	const handleDrawerToggle = () => {
		props.onSidebarToggle();
	};

	return (
		<div className={classes.root}>
			<CssBaseline />
			<AppBar position="fixed" className={classes.appBar}>
				<Toolbar>
					<Hidden mdUp>
						<IconButton edge="start" className={classes.menuButton} color="inherit" aria-label="menu" onClick={handleDrawerToggle}>
							<MenuIcon />
						</IconButton>
					</Hidden>

					<Hidden smDown>
						<Typography className={classes.title} variant="h6" noWrap>
							Brunel CI
						</Typography>
					</Hidden>

					<Hidden xsDown>
						<Button className={classes.buttonActive} component={NavLink} activeClassName='is-active' to={'/repository'} color="inherit">Repositories</Button>
						<Button className={classes.buttonActive} component={NavLink} activeClassName='is-active' to={'/environment'} color="inherit">Environments</Button>
					</Hidden>
					<div className={classes.grow} />
					<CurrentUser/>
				</Toolbar>
			</AppBar>
			<div style={{width: '100%'}}>
				<div className={classes.toolbar} />
				<main className={classes.content}>
					{children}
				</main>
			</div>
		</div>
	);
}
