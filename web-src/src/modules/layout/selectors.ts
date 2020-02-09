import {State} from './reducer';

export function getAuthenticated(state: State) {
	return state.authenticated;
}

