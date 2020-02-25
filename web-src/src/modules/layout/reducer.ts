import * as Actions from './actions';

export interface State {
	authenticated: boolean;
}

const initialState: State = {
	authenticated: false,
};

export function reducer(state = initialState, action: Actions.ActionTypes) {
	switch (action.type) {
	case Actions.Type.SET_AUTHENTICATED:
		return {
			...state,
			authenticated: action.payload.authenticated,
		};
	default:
		return state;
	}
}
