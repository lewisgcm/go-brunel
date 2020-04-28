export class Type {
	static SET_AUTHENTICATED: '@@LAYOUT/SET_AUTHENTICATED' = '@@LAYOUT/SET_AUTHENTICATED';
	static TOGGLE_SIDEBAR: '@@LAYOUT/TOGGLE_SIDEBAR' = '@@LAYOUT/TOGGLE_SIDEBAR';
}

export function setAuthenticated(authenticated: boolean) {
	return {
		type: Type.SET_AUTHENTICATED,
		payload: {
			authenticated,
		},
	};
}

export function toggleSidebar(isOpen?: boolean) {
	return {
		type: Type.TOGGLE_SIDEBAR,
		payload: {
			isOpen,
		},
	};
}

export type ActionTypes = ReturnType<typeof setAuthenticated> | ReturnType<typeof toggleSidebar>;
