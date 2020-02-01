import * as Actions from './actions';

export interface State {
	title: string;
	username: string;
	avatarUrl: string;
	authenticated: boolean;
}

const initialState: State = {
	title: '',
	username: '',
	avatarUrl: '',
	authenticated: false,
};

export function reducer(state = initialState, action: Actions.ActionTypes) {
	switch (action.type) {
	case Actions.Type.SET_USERNAME:
		return {
			...state,
			username: action.payload.username,
		};
	case Actions.Type.SET_AVATAR_URL:
		return {
			...state,
			avatarUrl: action.payload.avatarUrl,
		};
	case Actions.Type.SET_AUTHENTICATED:
		return {
			...state,
			authenticated: action.payload.authenticated,
		};
	default:
		return state;
	}
}
