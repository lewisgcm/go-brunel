import * as Actions from './actions';

export interface State {
	authenticated: boolean;
	sideBarOpen: boolean;
}

const initialState: State = {
	authenticated: false,
	sideBarOpen: false,
};

export function reducer(state = initialState, action: Actions.ActionTypes) {
	switch (action.type) {
	case Actions.Type.SET_AUTHENTICATED:
		return {
			...state,
			authenticated: action.payload.authenticated,
		};
	case Actions.Type.TOGGLE_SIDEBAR:
		return {
			...state,
			sideBarOpen: action.payload.isOpen !== undefined ? action.payload.isOpen : !state.sideBarOpen,
		};
	default:
		return state;
	}
}
