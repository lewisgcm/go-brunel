import React from 'react';
import {Redirect, Route, RouteProps} from 'react-router';

interface Props extends RouteProps {
	isAuthenticated: boolean;
}

export const ProtectedRoute = ({isAuthenticated, ...props}: Props) => {
	return isAuthenticated ?
		<Route {...props} /> :
		<Redirect to={'/user/login'} />;
};

