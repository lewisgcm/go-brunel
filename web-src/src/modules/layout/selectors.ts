import {State} from './reducer';

export function getAuthenticated(state: State) {
	return state.authenticated;
}

export function getSideBarOpen(state: State) {
	return state.sideBarOpen;
}


