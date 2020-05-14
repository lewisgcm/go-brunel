import React, {useEffect, useState, Dispatch} from 'react';
import {match, useHistory} from 'react-router';
import {connect} from 'react-redux';
import {BehaviorSubject, merge} from 'rxjs';
import {first, distinctUntilChanged, skip, debounceTime, tap, switchMap} from 'rxjs/operators';

import {Drawer, ActionTypes, toggleSidebar, SearchableList, SearchListState} from '../../layout';
import {useDependency} from '../../../container';
import {UserList, UserService, User} from '../../../services';
import { UserDetail } from '../components/UserDetail';

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
	const [listState, setListState] = useState(SearchListState.Loaded);
	const [users, setUsers] = useState<UserList[]>([]);
	const [selectedUser, setSelectedUser] = useState<User | undefined>(undefined);

	useEffect(
		() => {
			if (match.params.username) {
				userService
					.get(match.params.username)
					.subscribe(
						(user) => {
							setSelectedUser(user);
						},
						() => {
							setSelectedUser(undefined);
						},
					);
			}
		},
		[userService, match.params.username],
	);

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
				tap(() => setListState(SearchListState.Loading)),
				switchMap((term) => userService.list(term)),
				tap(() => setListState(SearchListState.Loaded)),
			).subscribe(
				(users) => {
					setUsers(users);
				},
				() => {
					setUsers([]);
					setListState(SearchListState.Error);
				},
			);

			return () => {
				subscription.unsubscribe();
			};
		},
		[userService, subject],
	);

	useEffect(
		() => {
			if (users && users.length && !match.params.username) {
				history.replace(`/user/${users[0].Username}`);
			}
		},
		[users, history, match.params.username],
	);

	const sidebar = () => <SearchableList
		state={listState}
		emptyPlaceholder='No users found.'
		errorPlaceholder='Error fetching users.'
		searchPlaceholder='Search for a user'
		items={users}
		render={(item: UserList) => ({
			selected: match.params.username === item.Username,
			text: item.Username,
			key: item.Username,
		})}
		onClick={(user) => {
			history.replace(`/user/${user.Username}`);
			hideMobileSidebar();
		}}
		onSearch={(term) => subject.next(term)} />;

	const content = () => selectedUser ?
		<UserDetail user={selectedUser} /> :
		<React.Fragment></React.Fragment>;

	return <Drawer
		sidebar={sidebar}
		content={content} />;
});

