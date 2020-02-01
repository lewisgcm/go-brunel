import React from 'react';
import {connect} from 'react-redux';
import {Route} from 'react-router';
import {Layout, State} from './modules/layout';
import {BrowserRouter, Switch} from 'react-router-dom';

import {RepositoryRoutes} from './modules/repository';
import {ProtectedRoute} from './modules/layout/ProtectedRoute';
import {getAuthenticated} from './modules/layout/selectors';
import {Login} from './modules/user/Login';

require('./App.css');

export default connect(
	(state: { layout: State }) => ({
		isAuthenticated: getAuthenticated(state.layout),
	}))(({isAuthenticated}) => {
	return (
		<BrowserRouter>
			<Layout>
				<Switch>
					<ProtectedRoute isAuthenticated={isAuthenticated}
						path='/repository'
						component={RepositoryRoutes} />
					<Route path={'/user/login'} component={Login} exact/>
				</Switch>
			</Layout>
		</BrowserRouter>
	);
});
