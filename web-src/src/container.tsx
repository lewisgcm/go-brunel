import React, { useContext, useMemo } from "react";
import { Container, interfaces } from "inversify";

const context = React.createContext<Container>({} as Container);

export const DependencyProvider = context.Provider;

export function withDependency<P, T>(resolver: (container: Container) => T) {
	return (WrappedComponent: React.ComponentType<P & T>) => {
		return (props: P) => {
			const container = useContext(context);
			const dependencies = resolver(container);
			const totalProps = { ...props, ...dependencies };

			return <WrappedComponent {...totalProps} />;
		};
	};
}

export function useDependency<T>(
	serviceIdentifier: interfaces.ServiceIdentifier<T>
) {
	const container: Container = useContext(context);
	return useMemo(() => container.get<T>(serviceIdentifier), [
		container,
		serviceIdentifier,
	]);
}
