import React, {useEffect, useState, Dispatch} from 'react';
import {match, useHistory} from 'react-router';
import {BehaviorSubject, merge} from 'rxjs';
import {debounceTime, first, skip, switchMap, tap, distinctUntilChanged} from 'rxjs/operators';

import {Drawer, ActionTypes, toggleSidebar} from '../../layout';
import {useDependency} from '../../../container';
import {UserList, UserService, User} from '../../../services';
import {UserListComponent} from '../components/UserListComponent';
import {connect} from 'react-redux';

interface Props {
	match: match<{username: string}>;
}

function mapDispatchToProps(dispatch: Dispatch<ActionTypes>) {
	return {
		hideMobileSidebar: () => {
			dispatch(toggleSidebar(false));
		},
	};
}

export const UserPage = connect(
	null,
	mapDispatchToProps,
)(({match, hideMobileSidebar}: Props & ReturnType<typeof mapDispatchToProps>) => {
	const userService = useDependency(UserService);
	const history = useHistory();
	const [subject] = useState(new BehaviorSubject(''));
	const [isLoading, setLoading] = useState(false);
	const [users, setUsers] = useState<UserList[]>([]);
	const [selectedUsername, setSelectedUsername] = useState<string | undefined>(undefined);
	const [selectedUser, setSelectedUser] = useState<User | undefined>(undefined);

	useEffect(() => {
		setSelectedUsername(match.params.username);
	}, [match.params]);

	useEffect(() => {
		if (selectedUsername) {
			userService
				.get(selectedUsername)
				.subscribe((user) => {
					setSelectedUser(user);
				});
		}
	}, [userService, selectedUsername]);

	useEffect(
		() => {
			const subscription = merge(
				subject.pipe(first()),
				subject.pipe(
					distinctUntilChanged(),
					skip(1),
					debounceTime(200),
				),
			).pipe(
				tap(() => setLoading(true)),
				switchMap((term) => userService.list(term)),
				tap(() => setLoading(false)),
			).subscribe(
				(items) => {
					setUsers(items);
					if (items.length && (match.params && !match.params.username)) {
						setSelectedUsername(items[0].Username);
						history.push(`/user/${items[0].Username}`);
					}
				},
			);

			return () => {
				subscription.unsubscribe();
			};
		},
		[userService, subject, history, match.params],
	);

	return <Drawer
		sidebar={() => <UserListComponent
			isLoading={isLoading}
			users={users}
			selectedUsername={selectedUsername}
			onClick={(user) => {
				hideMobileSidebar();
				setSelectedUsername(user.Username);
			}}
			onSearch={(term) =>subject.next(term)}/>}
		content={() => <div>
			<h1>{selectedUser && selectedUser.Name}</h1>
		</div>}/>;
});

