import React, {useState, useEffect} from 'react';
import {match} from 'react-router';
import {BehaviorSubject} from 'rxjs';
import {tap, delay, debounceTime, switchMap} from 'rxjs/operators';

import {withDependency} from '../../../container';
import {RepositoryJobsComponent} from '../components/RepositoryJobsComponent';
import {RepositoryService, RepositoryJobPage} from '../../../services';

interface Props {
	match: match<{ repositoryId: string }>;
}

interface Dependencies {
	repositoryService: RepositoryService;
}

interface QueryParams {
	repositoryId: string;
	rowsPerPage: number;
	currentPage: number;
	sortOrder: 'asc' | 'desc';
	sortColumn: string;
	search: string;
}

export const RepositoryJobs = withDependency<Props, Dependencies>(
	(container) => ({
		repositoryService: container.get(RepositoryService),
	}),
)(
	({repositoryService, match}) => {
		const {repositoryId} = match.params;
		const rowsPerPageOptions = [5, 10, 15, 20];
		const [subject] = useState(new BehaviorSubject<QueryParams>({} as any));
		const [isLoading, setIsLoading] = useState(false);
		const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
		const [sortColumn, setSortColumn] = useState('created_at');
		const [rowsPerPage, setRowsPerPage] = useState(5);
		const [currentPage, setCurrentPage] = useState(0);
		const [search, setSearch] = useState('');
		const [page, setPage] = useState<RepositoryJobPage>({Count: 0, Jobs: []});

		const onSortChange = (sortColumn: string) => {
			setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
			setSortColumn(sortColumn);
		};

		useEffect(
			() => {
				setCurrentPage(0);

				const subscription = subject
					.pipe(
						debounceTime(200),
						tap((s) => {
							setIsLoading(true);
						}),
						delay(400),
						switchMap((state) => repositoryService
							.jobs(
								state.repositoryId,
								state.search,
								state.currentPage,
								state.rowsPerPage,
								state.sortOrder,
								state.sortColumn,
							),
						),
						tap((ss) => {
							setIsLoading(false);
						}),
					).subscribe(
						(jobs) => {
							setPage(jobs);
						},
					);

				return () => {
					subscription.unsubscribe();
				};
			},
			[repositoryId, subject, repositoryService],
		);

		useEffect(
			() => {
				subject.next({
					repositoryId,
					rowsPerPage,
					currentPage,
					sortOrder,
					sortColumn,
					search,
				});
			},
			[
				subject,
				repositoryId,
				rowsPerPage,
				currentPage,
				sortOrder,
				sortColumn,
				search,
			],
		);

		return <RepositoryJobsComponent
			isLoading={isLoading}
			sortColumn={sortColumn}
			sortOrder={sortOrder}
			rowsPerPageOptions={rowsPerPageOptions}
			rowsPerPage={rowsPerPage}
			currentPage={currentPage}
			page={page}
			onSortChange={onSortChange}
			onPageChange={(p) => {
				setCurrentPage(p);
			}}
			onRowsPerPageChange={(r) => {
				setRowsPerPage(r);
			}}
			onSearch={(t) => {
				setSearch(t);
			}}
		/>;
	},
);
