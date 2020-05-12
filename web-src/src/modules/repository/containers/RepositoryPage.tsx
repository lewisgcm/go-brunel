import React, {useEffect, useState, Dispatch} from 'react';
import {match, useHistory} from 'react-router';
import {BehaviorSubject, merge} from 'rxjs';
import {debounceTime, first, skip, switchMap, tap, distinctUntilChanged} from 'rxjs/operators';
import {connect} from 'react-redux';

import {Drawer, ActionTypes, toggleSidebar, SearchableList, SearchListState} from '../../layout';
import {RepositoryJobs} from './RepositoryJobs';
import {useDependency} from '../../../container';
import {Repository, RepositoryService} from '../../../services';

interface Props {
	match: match<{repositoryId: string}>;
}

function mapDispatchToProps(dispatch: Dispatch<ActionTypes>) {
	return {
		hideMobileSidebar: () => {
			dispatch(toggleSidebar(false));
		},
	};
}

export const RepositoryPage = connect(
	null,
	mapDispatchToProps,
)(({match, hideMobileSidebar}: Props & ReturnType<typeof mapDispatchToProps>) => {
	const repositoryService = useDependency(RepositoryService);
	const history = useHistory();
	const [subject] = useState(new BehaviorSubject(''));
	const [listState, setListState] = useState(SearchListState.Loaded);
	const [repositories, setRepositories] = useState<Repository[]>([]);
	const [selectedRepository, setSelectedRepository] = useState<Repository | undefined>();

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
				switchMap((term) => repositoryService.list(term)),
				tap(() => setListState(SearchListState.Loaded)),
			).subscribe(
				(items) => {
					setRepositories(items);
				},
				() => {
					setListState(SearchListState.Error);
				},
			);

			return () => {
				subscription.unsubscribe();
			};
		},
		[repositoryService, subject],
	);

	useEffect(
		() => {
			if (repositories.length && !match.params.repositoryId) {
				history.push(`/repository/${repositories[0].ID}`);
			} else if (repositories.length && match.params.repositoryId) {
				setSelectedRepository(
					repositories.find((r) => r.ID === match.params.repositoryId),
				);
			}
		},
		[repositories, match.params.repositoryId, history],
	);

	const sidebar = () => <SearchableList
		emptyPlaceholder='No repositories found.'
		errorPlaceholder='Error fetching repositories.'
		searchPlaceholder='Search for a repository'
		state={listState}
		items={repositories}
		render={(item) => ({
			selected: item.ID === match.params.repositoryId,
			text: `${item.Project}/${item.Name}`,
			key: item.ID,
		})}
		onClick={(r) => {
			hideMobileSidebar();
			history.push(`/repository/${r.ID}`);
		}}
		onSearch={(term) =>subject.next(term)}/>;

	const content = () => selectedRepository ?
		<RepositoryJobs repository={selectedRepository}/> :
		<React.Fragment/>;

	return <Drawer
		sidebar={sidebar}
		content={content}/>;
});

