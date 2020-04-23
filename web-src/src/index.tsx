import 'reflect-metadata';
import React from 'react';
import ReactDOM from 'react-dom';
import {Provider} from 'react-redux';
import {Container} from 'inversify';

import App from './App';
import {store} from './store';
import {DependencyProvider} from './container';
import {setAuthenticated} from './modules/layout';
import {
	AuthService,
	EnvironmentService,
	JobService,
	RepositoryService,
	UserService,
} from './services';

const container = new Container();
const authService = new AuthService();

container.bind(AuthService).toConstantValue(authService);
container.bind(RepositoryService).toSelf();
container.bind(UserService).toSelf();
container.bind(JobService).toSelf();
container.bind(EnvironmentService).toSelf();

store.dispatch(setAuthenticated(authService.isAuthenticated()));

ReactDOM.render(
	<DependencyProvider value={container}>
		<Provider store={store}>
			<App />
		</Provider>
	</DependencyProvider>,
	document.getElementById('root'),
);
