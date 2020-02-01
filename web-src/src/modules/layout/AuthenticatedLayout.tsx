import React from 'react';
import AppBar from '@material-ui/core/AppBar';
import CssBaseline from '@material-ui/core/CssBaseline';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import {
	makeStyles,
	Theme,
	createStyles,
} from '@material-ui/core/styles';

import {CurrentUser} from './CurrentUser';

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
			[theme.breakpoints.up('sm')]: {
				display: 'none',
			},
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
					<IconButton
						color="inherit"
						edge="start"
						onClick={handleDrawerToggle}
						className={classes.menuButton} >
						<MenuIcon />
					</IconButton>
					<Typography variant="h6" noWrap>
						Brunel CI
					</Typography>
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
