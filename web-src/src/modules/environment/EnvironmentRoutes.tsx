import React from 'react';
import {Route, match} from 'react-router-dom';

import {EnvironmentPage} from './containers/EnvironmentPage';

interface Props {
	match: match;
}

export function EnvironmentRoutes({match}: Props) {
	return <React.Fragment>
		<Route path={`${match.path}`}
			   exact
			   component={EnvironmentPage} />
	</React.Fragment>;
}
