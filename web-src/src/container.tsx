import React, {useContext} from 'react';
import {Container} from 'inversify';

const context = React.createContext<Container>({} as Container);

export const DependencyProvider = context.Provider;

export function withDependency<P, T>(
	resolver: (container: Container) => T,
) {
	return (WrappedComponent: React.ComponentType<P & T>) => {
		return (props: P) => {
			const container = useContext(context);
			const dependencies = resolver(container);
			const totalProps = {...props, ...dependencies};

			return <WrappedComponent {...totalProps} />;
		};
	};
}
