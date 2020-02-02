import React, {useEffect, useState} from 'react';
import {match, useHistory} from 'react-router';
import {BehaviorSubject} from 'rxjs';
import {debounceTime, switchMap, tap} from 'rxjs/operators';

import {Drawer} from '../../layout';
import {RepositoryJobs} from './RepositoryJobs';
import {withDependency} from '../../../container';
import {Repository, RepositoryService} from '../../../services';
import {RepositoryListComponent} from '../components/RepositoryListComponent';

interface Props {
	match: match<{repositoryId: string}>;
}

interface Dependencies {
	repositoryService: RepositoryService;
}

function content(match: match<{repositoryId: string}>) {
	return () => {
		if (match.params && match.params.repositoryId) {
			return <RepositoryJobs match={match}/>;
		}
		return <React.Fragment />;
	};
}

export const RepositoryPage = withDependency<Props, Dependencies>(
	(container) => ({
		repositoryService: container.get(RepositoryService),
	}),
)(({match, repositoryService}) => {
	const history = useHistory();
	const [subject] = useState(new BehaviorSubject(''));
	const [search, setSearch] = useState('');
	const [isLoading, setLoading] = useState(false);
	const [repositoryItems, setRepositoryItems] = useState<Repository[]>([]);

	useEffect(
		() => {
			const subscription = subject
				.pipe(
					tap(() => setLoading(true)),
					debounceTime(300),
					switchMap((term) => repositoryService.list(term)),
					tap(() => setLoading(false)),
				).subscribe(
					(items) => {
						setRepositoryItems(items);
						if (items.length && (match.params && !match.params.repositoryId)) {
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

	useEffect(() => {
		subject.next(search);
	}, [search, subject]);

	return <Drawer
		sidebar={() => <RepositoryListComponent
			isLoading={isLoading}
			repositories={repositoryItems}
			onClick={() => {}}
			onSearch={setSearch}/>}
		content={content(match)}/>;
});

