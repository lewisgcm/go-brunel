import React, {useEffect, useState, Dispatch} from 'react';
import {match, useHistory} from 'react-router';
import {BehaviorSubject, merge} from 'rxjs';
import {debounceTime, first, skip, switchMap, tap, distinctUntilChanged} from 'rxjs/operators';

import {Drawer, ActionTypes, toggleSidebar} from '../../layout';
import {RepositoryJobs} from './RepositoryJobs';
import {useDependency} from '../../../container';
import {Repository, RepositoryService} from '../../../services';
import {RepositoryListComponent} from '../components/RepositoryListComponent';
import {connect} from 'react-redux';

interface Props {
	match: match<{repositoryId: string}>;
}

function content(selectedRepository: Repository | undefined) {
	return () => selectedRepository ?
		<RepositoryJobs repository={selectedRepository}/> :
		<React.Fragment/>;
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
	const [isLoading, setLoading] = useState(false);
	const [repositoryItems, setRepositoryItems] = useState<Repository[]>([]);
	const [selectedRepositoryId, setSelectedRepositoryId] = useState<string | undefined>(match.params.repositoryId);
	const [selectedRepository, setSelectedRepository] = useState<Repository | undefined>();

	useEffect(() => {
		if (repositoryItems.length > 0 && (!selectedRepository || (selectedRepository.ID !== selectedRepositoryId))) {
			setSelectedRepository(
				repositoryItems.find((r) => r.ID === selectedRepositoryId),
			);
		}
	}, [repositoryItems, selectedRepositoryId, selectedRepository]);

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
				switchMap((term) => repositoryService.list(term)),
				tap(() => setLoading(false)),
			).subscribe(
				(items) => {
					setRepositoryItems(items);
					if (items.length && (match.params && !match.params.repositoryId)) {
						setSelectedRepositoryId(items[0].ID);
						history.push(`/repository/${items[0].ID}`);
					}
				},
			);

			return () => {
				subscription.unsubscribe();
			};
		},
		[repositoryService, subject, history, match.params],
	);

	return <Drawer
		sidebar={() => <RepositoryListComponent
			isLoading={isLoading}
			repositories={repositoryItems}
			selectedRepositoryId={selectedRepositoryId}
			onClick={(r) => {
				hideMobileSidebar();
				setSelectedRepositoryId(r.ID);
			}}
			onSearch={(term) =>subject.next(term)}/>}
		content={content(selectedRepository)}/>;
});

