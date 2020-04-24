import React from 'react';
import {Link} from 'react-router-dom';
import AppBar from '@material-ui/core/AppBar';
import {Button, IconButton, Toolbar, Drawer, Hidden} from '@material-ui/core';
import MenuIcon from '@material-ui/icons/Menu';
import CloseIcon from '@material-ui/icons/Close';
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
	container?: Element;
	children: React.ReactNode;
}

export function AuthenticatedLayout(props: ResponsiveDrawerProps) {
	const {children} = props;
	const classes = useStyles();
	const [mobileOpen, setMobileOpen] = React.useState(false);

	const handleDrawerToggle = () => {
		setMobileOpen(!mobileOpen);
	};

	return (
		<div className={classes.root}>
			<CssBaseline />
			<AppBar position="fixed" className={classes.appBar}>
				<Toolbar>
					<Hidden smUp>
						<IconButton edge="start" className={classes.menuButton} color="inherit" aria-label="menu" onClick={handleDrawerToggle}>
							<MenuIcon />
						</IconButton>
					</Hidden>

					<Typography className={classes.title} variant="h6" noWrap>
						Brunel CI
					</Typography>

					<Hidden xsDown>
						<Button component={Link} to={'/repository'} color="inherit">Repositories</Button>
						<Button component={Link} to={'/environment'} color="inherit">Environments</Button>
						<div className={classes.grow} />
						<CurrentUser/>
					</Hidden>
				</Toolbar>
			</AppBar>
			<div style={{width: '100%'}}>
				<div className={classes.toolbar} />
				<nav>
					<Hidden smUp implementation="css">
						<Drawer
							variant="temporary"
							open={mobileOpen}
							onClose={handleDrawerToggle}
							classes={{
								paper: classes.drawerPaper,
							}}
							ModalProps={{
								keepMounted: true,
							}}>
							{<React.Fragment>
								<div className={classes.drawerHeader}>
									<IconButton onClick={handleDrawerToggle}>
										<CloseIcon />
									</IconButton>
								</div>
							</React.Fragment>}
						</Drawer>
					</Hidden>
				</nav>
				<main className={classes.content}>
					{children}
				</main>
			</div>
		</div>
	);
}
