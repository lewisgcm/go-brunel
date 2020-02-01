import React from 'react';
import {Route, match} from 'react-router-dom';

import {RepositoryPage} from './RepositoryPage';

interface Props {
	match: match;
}

export function RepositoryRoutes({match}: Props) {
	return <React.Fragment>
		<Route path={`${match.path}/:repositoryId?`}
			exact
			component={RepositoryPage} />
	</React.Fragment>;
}
