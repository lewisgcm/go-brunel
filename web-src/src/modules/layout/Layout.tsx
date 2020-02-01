import React, {PropsWithChildren} from 'react';
import {connect} from 'react-redux';

import {AuthenticatedLayout} from './AuthenticatedLayout';
import {State} from './reducer';
import {getAuthenticated} from './selectors';

interface Props {
	isAuthenticated: boolean;
}

export const Layout = connect(
	(state: { layout: State }) => ({
		isAuthenticated: getAuthenticated(state.layout),
	}),
)(
	({isAuthenticated, children}: PropsWithChildren<Props>) => {
		if (!isAuthenticated) {
			return <React.Fragment>
				{children}
			</React.Fragment>;
		}

		return <AuthenticatedLayout>
			{children}
		</AuthenticatedLayout>;
	},
);
