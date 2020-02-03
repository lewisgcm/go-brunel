import React from 'react';
import {Route, match} from 'react-router-dom';

import {JobComponent} from './components/JobComponent';

interface Props {
    match: match;
}

export function JobRoutes({match}: Props) {
	return <React.Fragment>
		<Route path={`${match.path}/:jobId`}
			exact
			component={JobComponent} />
	</React.Fragment>;
}
