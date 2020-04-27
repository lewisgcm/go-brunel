import React, {Dispatch} from 'react';
import {connect} from 'react-redux';
import {Hidden, Drawer as MaterialDrawer, IconButton} from '@material-ui/core';
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';
import CloseIcon from '@material-ui/icons/Close';

import {State} from '../reducer';
import {getSideBarOpen} from '../selectors';
import {ActionTypes, toggleSidebar} from '../actions';

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
		drawerHeader: {
			display: 'flex',
			alignItems: 'center',
			padding: theme.spacing(0, 1),
			...theme.mixins.toolbar,
			justifyContent: 'flex-end',
		},
		content: {
			marginLeft: drawerWidth,
			[theme.breakpoints.down('sm')]: {
				marginLeft: 0,
			},
		},
	}),
);

function mapStateToProps(state: { layout: State }, ownProps: Props) {
	return {
		...ownProps,
		isSideBarOpen: getSideBarOpen(state.layout),
	};
}

function mapDispatchToProps(dispatch: Dispatch<ActionTypes>) {
	return {
		onToggleSideBar: () => {
			dispatch(toggleSidebar());
		},
	};
}

export const Drawer = connect(
	mapStateToProps,
	mapDispatchToProps,
)(({sidebar, content, isSideBarOpen, onToggleSideBar}: ReturnType<typeof mapStateToProps> & ReturnType<typeof mapDispatchToProps>) => {
	const classes = useStyles({});

	return <React.Fragment>
		<nav className={classes.drawer}>
			<Hidden mdUp implementation="css">
				<MaterialDrawer
					variant="temporary"
					open={isSideBarOpen}
					onClose={() => onToggleSideBar()}
					classes={{
						paper: classes.drawerPaper,
					}}
					ModalProps={{
						keepMounted: true, // Better open performance on mobile.
					}}
				>
					<div>
						<div className={classes.drawerHeader}>
							<IconButton onClick={() => {
								onToggleSideBar();
							}}>
								<CloseIcon/>
							</IconButton>
						</div>
						{sidebar()}
					</div>
				</MaterialDrawer>
			</Hidden>

			<Hidden smDown implementation="css">
				<MaterialDrawer
					classes={{
						paper: classes.drawerPaper,
					}}
					variant="permanent"
					open={isSideBarOpen} >
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
});
