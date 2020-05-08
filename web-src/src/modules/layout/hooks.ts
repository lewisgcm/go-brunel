import {useSelector} from 'react-redux';
import {UserRole} from '../../services';
import {State} from './reducer';

export function useHasRole(role: UserRole) {
	return useSelector<{layout: State}, boolean>(
		(s) => role === s.layout.role,
	);
}
