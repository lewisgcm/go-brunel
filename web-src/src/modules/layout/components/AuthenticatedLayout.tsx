import React from 'react';
import {Link} from 'react-router-dom';
import AppBar from '@material-ui/core/AppBar';
import {Button} from '@material-ui/core';
import CssBaseline from '@material-ui/core/CssBaseline';
import Toolbar from '@material-ui/core/Toolbar';
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
	container?: Element;
	children: React.ReactNode;
}

export function AuthenticatedLayout(props: ResponsiveDrawerProps) {
	const {children} = props;
	const classes = useStyles();

	return (
		<div className={classes.root}>
			<CssBaseline />
			<AppBar position="fixed" className={classes.appBar}>
				<Toolbar>
					<Typography className={classes.title} variant="h6" noWrap>
						Brunel CI
					</Typography>
					<Button component={Link} to={'/repository'} color="inherit">Repositories</Button>
					<Button component={Link} to={'/environment'} color="inherit">Environments</Button>
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
