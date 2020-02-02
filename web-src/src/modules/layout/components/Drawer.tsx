import React from 'react';
import {Hidden, Drawer as MaterialDrawer} from '@material-ui/core';
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';

interface Props {
	sidebar: () => React.ReactNode;
	content: () => React.ReactNode;
}

const drawerWidth = 320;

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		drawer: {
			width: drawerWidth,
			flexShrink: 0,
		},
		toolbar: theme.mixins.toolbar,
		drawerPaper: {
			width: drawerWidth,
		},
		content: {
			marginLeft: drawerWidth,
		},
	}),
);

export function Drawer({sidebar, content}: Props) {
	const classes = useStyles();

	return <React.Fragment>
		<nav className={classes.drawer}>
			<Hidden xsDown implementation="css">
				<MaterialDrawer
					classes={{
						paper: classes.drawerPaper,
					}}
					variant="permanent"
					open >
					<div>
						<div className={classes.toolbar} />
						{sidebar()}
					</div>
				</MaterialDrawer>
			</Hidden>
		</nav>
		<div className={classes.content}>
			{content()}
		</div>
	</React.Fragment>;
}
