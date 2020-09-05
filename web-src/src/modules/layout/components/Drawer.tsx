import React, { Dispatch } from "react";
import { NavLink } from "react-router-dom";
import { connect } from "react-redux";
import {
	Hidden,
	Drawer as MaterialDrawer,
	IconButton,
	List,
	ListItem,
	Typography,
	ListItemText,
	Divider,
} from "@material-ui/core";
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";
import CloseIcon from "@material-ui/icons/Close";

import { State } from "../reducer";
import { getSideBarOpen, getRole } from "../selectors";
import { ActionTypes, toggleSidebar } from "../actions";
import { UserRole } from "../../../services";

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
		grow: {
			flexGrow: 1,
		},
		headerTitle: {
			marginLeft: theme.spacing(1),
		},
		drawerHeader: {
			display: "flex",
			alignItems: "center",
			padding: theme.spacing(0, 1),
			...theme.mixins.toolbar,
			justifyContent: "flex-end",
		},
		content: {
			marginLeft: drawerWidth,
			[theme.breakpoints.down("sm")]: {
				marginLeft: 0,
			},
		},
	})
);

function mapStateToProps(state: { layout: State }, ownProps: Props) {
	return {
		...ownProps,
		isSideBarOpen: getSideBarOpen(state.layout),
		isAdmin: getRole(state.layout) === UserRole.Admin,
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
	mapDispatchToProps
)(
	({
		sidebar,
		content,
		isSideBarOpen,
		onToggleSideBar,
		isAdmin,
	}: ReturnType<typeof mapStateToProps> &
		ReturnType<typeof mapDispatchToProps>) => {
		const classes = useStyles({});

		return (
			<React.Fragment>
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
									<Typography
										variant="h6"
										className={classes.headerTitle}
										noWrap
									>
										Brunel CI
									</Typography>
									<div className={classes.grow} />
									<IconButton
										onClick={() => {
											onToggleSideBar();
										}}
									>
										<CloseIcon />
									</IconButton>
								</div>
								<Divider />
								<List>
									<ListItem
										button
										component={NavLink}
										to={"/repository"}
										activeClassName="Mui-selected"
										onClick={() => onToggleSideBar()}
									>
										<ListItemText
											primary={"Repositories"}
										/>
									</ListItem>
									{isAdmin && (
										<ListItem
											button
											component={NavLink}
											to={"/environment"}
											activeClassName="Mui-selected"
											onClick={() => onToggleSideBar()}
										>
											<ListItemText
												primary={"Environments"}
											/>
										</ListItem>
									)}
									{isAdmin && (
										<ListItem
											button
											component={NavLink}
											to={"/user"}
											activeClassName="Mui-selected"
											onClick={() => onToggleSideBar()}
										>
											<ListItemText primary={"Users"} />
										</ListItem>
									)}
								</List>
								<Divider />
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
							open={isSideBarOpen}
						>
							<div>
								<div className={classes.toolbar} />
								{sidebar()}
							</div>
						</MaterialDrawer>
					</Hidden>
				</nav>
				<div className={classes.content}>{content()}</div>
			</React.Fragment>
		);
	}
);
