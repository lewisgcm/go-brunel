import React from 'react';
import {Route, match} from 'react-router-dom';

import {UserPage} from './containers/UserPage';

interface Props {
	match: match;
}

export function UserRoutes({match}: Props) {
	return <React.Fragment>
		<Route path={`${match.path}/:username?`}
			exact
			component={UserPage} />
	</React.Fragment>;
}
