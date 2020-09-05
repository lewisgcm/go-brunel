import React, { PropsWithChildren } from "react";
import { connect } from "react-redux";

import { AuthenticatedLayout } from "../components/AuthenticatedLayout";
import { State } from "../reducer";
import { getAuthenticated } from "../selectors";
import { toggleSidebar } from "../actions";

interface Props {
	isAuthenticated: boolean;
	onSidebarToggle: () => void;
}

export const Layout = connect(
	(state: { layout: State }) => ({
		isAuthenticated: getAuthenticated(state.layout),
	}),
	(dispatch) => ({
		onSidebarToggle: () => dispatch(toggleSidebar()),
	})
)(
	({
		isAuthenticated,
		children,
		onSidebarToggle,
	}: PropsWithChildren<Props>) => {
		if (!isAuthenticated) {
			return <React.Fragment>{children}</React.Fragment>;
		}

		return (
			<AuthenticatedLayout onSidebarToggle={() => onSidebarToggle()}>
				{children}
			</AuthenticatedLayout>
		);
	}
);
