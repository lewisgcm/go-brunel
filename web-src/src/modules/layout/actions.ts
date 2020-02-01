export enum Type {
	SET_USERNAME = '@@LAYOUT/SET_USERNAME',
	SET_AVATAR_URL = '@@LAYOUT/SET_AVATAR_URL',
	SET_AUTHENTICATED = '@@LAYOUT/SET_AUTHENTICATED'
}

export function setUsername(username: string) {
	return {
		type: Type.SET_USERNAME,
		payload: {
			username,
		},
	};
}

export function setAvatarUrl(avatarUrl: string) {
	return {
		type: Type.SET_AVATAR_URL,
		payload: {
			avatarUrl,
		},
	};
}

export function setAuthenticated(authenticated: boolean) {
	return {
		type: Type.SET_AUTHENTICATED,
		payload: {
			authenticated,
		},
	};
}

export type ActionTypes = ReturnType<typeof setUsername>
	& ReturnType<typeof setAvatarUrl> & ReturnType<typeof setAuthenticated>;
