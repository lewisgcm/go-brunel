import React from 'react';
import {connect} from 'react-redux';
import {Redirect, Route} from 'react-router';
import {BrowserRouter, Switch} from 'react-router-dom';

import {JobRoutes} from './modules/job';
import {RepositoryRoutes} from './modules/repository';
import {EnvironmentRoutes} from './modules/environment';
import {getAuthenticated, ProtectedRoute, Layout, State} from './modules/layout';
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

					<ProtectedRoute isAuthenticated={isAuthenticated}
						path='/job'
						component={JobRoutes} />

					<ProtectedRoute isAuthenticated={isAuthenticated}
						path='/environment'
						component={EnvironmentRoutes} />

					<Route path={'/user/login'} component={Login} exact />
					{!isAuthenticated && <Redirect to={'/user/login'} />}
					{isAuthenticated && <Redirect to={'/repository'} />}
				</Switch>
			</Layout>
		</BrowserRouter>
	);
});
