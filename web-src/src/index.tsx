import 'reflect-metadata';
import React from 'react';
import ReactDOM from 'react-dom';
import {Provider} from 'react-redux';
import {Container} from 'inversify';

import App from './App';
import {DependencyProvider} from './container';
import {AuthService, RepositoryService, UserService} from './services';
import {store} from './store';
import {setAuthenticated} from './modules/layout';

const container = new Container();

const authService = new AuthService();
container.bind(AuthService).toConstantValue(authService);
container.bind(RepositoryService).toSelf();
container.bind(UserService).toSelf();
store.dispatch(setAuthenticated(authService.isAuthenticated()))

ReactDOM.render(
	<DependencyProvider value={container}>
		<Provider store={store}>
			<App />
		</Provider>
	</DependencyProvider>,
	document.getElementById('root'),
);
