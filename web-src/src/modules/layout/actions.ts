export enum Type {
	SET_AUTHENTICATED = '@@LAYOUT/SET_AUTHENTICATED'
}

export function setAuthenticated(authenticated: boolean) {
	return {
		type: Type.SET_AUTHENTICATED,
		payload: {
			authenticated,
		},
	};
}

export type ActionTypes = ReturnType<typeof setAuthenticated>;
