import {State} from './reducer';

export function getUsername(state: State) {
	return state.username;
}

export function getAvatarUrl(state: State) {
	return state.avatarUrl;
}

export function getAuthenticated(state: State) {
	return state.authenticated;
}

