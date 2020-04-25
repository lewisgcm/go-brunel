import React, {useEffect, useState} from 'react';
import {merge, ReplaySubject} from 'rxjs';
import {debounceTime, first, skip, switchMap, tap} from 'rxjs/operators';

import {useDependency} from '../../../container';
import {RepositoryJobsComponent} from '../components/RepositoryJobsComponent';
import {Repository, RepositoryJobPage, RepositoryService} from '../../../services';

interface Props {
	repository: Repository;
}

interface QueryParams {
	selectedRepositoryId: string;
	rowsPerPage: number;
	currentPage: number;
	sortOrder: 'asc' | 'desc';
	sortColumn: string;
	search: string;
}

export const RepositoryJobs = ({repository}: Props) => {
	const rowsPerPageOptions = [5, 10, 15, 20];
	const [subject] = useState(new ReplaySubject<QueryParams>(1));
	const [isLoading, setIsLoading] = useState(false);
	const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
	const [sortColumn, setSortColumn] = useState('created_at');
	const [rowsPerPage, setRowsPerPage] = useState(5);
	const [currentPage, setCurrentPage] = useState(0);
	const [search, setSearch] = useState('');
	const [page, setPage] = useState<RepositoryJobPage>({Count: 0, Jobs: []});
	const repositoryService = useDependency(RepositoryService);

	const onSortChange = (sortColumn: string) => {
		setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
		setSortColumn(sortColumn);
	};

	useEffect(
		() => {
			setCurrentPage(0);

			const subscription = merge(
				subject.pipe(first()),
				subject.pipe(
					skip(1),
					debounceTime(200),
				),
			).pipe(
				tap((_) => {
					setIsLoading(true);
				}),
				switchMap((state) => repositoryService
					.jobs(
						state.selectedRepositoryId,
						state.search,
						state.currentPage,
						state.rowsPerPage,
						state.sortOrder,
						state.sortColumn,
					),
				),
				tap((_) => {
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
		[repositoryService, subject],
	);

	useEffect(
		() => {
			subject.next({
				selectedRepositoryId: repository.ID,
				rowsPerPage,
				currentPage,
				sortOrder,
				sortColumn,
				search,
			});
		},
		[
			subject,
			repository,
			rowsPerPage,
			currentPage,
			sortOrder,
			sortColumn,
			search,
		],
	);

	return <RepositoryJobsComponent
		isLoading={isLoading}
		repository={repository}
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
};
