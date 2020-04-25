import React from 'react';
import moment from 'moment';

import {
	makeStyles,
	Theme,
	createStyles,
	Paper,
	Table,
	TableRow,
	TablePagination,
	TableCell,
	TableBody,
	LinearProgress,
	Icon,
	TableHead,
	TableSortLabel,
	TextField, Tooltip,
} from '@material-ui/core';

import {RepositoryJobPage, JobState, RepositoryJob, Repository} from '../../../services';
import {useHistory} from 'react-router';

interface Props {
	isLoading: boolean;
	repository: Repository;
	page: RepositoryJobPage;
	sortColumn: string;
	sortOrder: 'asc' | 'desc';
	rowsPerPageOptions: number[];
	rowsPerPage: number;
	currentPage: number;
	onSortChange: (sortColumn: string) => void;
	onPageChange: (page: number) => void;
	onRowsPerPageChange: (rowsPerPage: number) => void;
	onSearch: (term: string) => void;
}

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		'sort': {
			'& .MuiTableSortLabel-icon': {
				opacity: 0,
				fontSize: '1.3em',
				fontWeight: 'bold',
			},
		},
		'duration': {
			width: '10em',
		},
		'hidden': {
			visibility: 'hidden',
		},
		'headerRow': {
			'& .MuiTableCell-root': {
				fontSize: `12px`,
				opacity: 0.7,
			},
		},
		'footer': {
			'fontSize': '12px !important',
			'opacity': 0.9,
			'& .MuiTypography-root': {
				fontSize: `12px`,
				opacity: 0.9,
			},
		},
		'search': {
			width: '100%',
			paddingBottom: theme.spacing(3),
		},
		'@keyframes spinRound': {
			from: {
				transform: 'rotate(0deg)',
			},
			to: {
				transform: 'rotate(360deg)',
			},
		},
		'inProgress': {
			color: theme.palette.grey.A200,
			position: 'relative',
			top: 3,
			animation: '$spinRound 2s linear infinite',
		},
		'cancelled': {
			color: theme.palette.grey.A200,
		},
	}),
);

function jobStatus(classes: any, state: JobState): React.ReactNode {
	switch (state) {
	case JobState.Processing:
	case JobState.Waiting:
		return <Tooltip title={'In Progress'}>
			<Icon className={classes.inProgress}>loop</Icon>
		</Tooltip>;
	case JobState.Failed:
		return <Tooltip title={'Failed'}>
			<Icon color="error">error</Icon>
		</Tooltip>;
	case JobState.Cancelled:
		return <Tooltip title={'Cancelled'}>
			<Icon className={classes.cancelled}>cancel</Icon>
		</Tooltip>;
	case JobState.Success:
		return <Tooltip title={'Success'}><Icon color="primary"
			style={{color: 'rgb(0, 100, 0)', position: 'relative', top: 3}} >
				check_circle
		</Icon></Tooltip>;
	}
}

function duration(job: RepositoryJob): string {
	if (job.StartedAt && job.StoppedAt) {
		return moment
			.duration(moment(job.StoppedAt).diff(moment(job.StartedAt))).humanize();
	}
	return 'N/A';
}

export function RepositoryJobsComponent(
	{
		isLoading,
		repository,
		page,
		sortOrder,
		sortColumn,
		rowsPerPageOptions,
		rowsPerPage,
		currentPage,
		onSortChange,
		onPageChange,
		onRowsPerPageChange,
		onSearch,
	}: Props,
) {
	const classes = useStyles();
	const history = useHistory();

	return (
		<div>
			<h1>{repository.Project}/{repository.Name}</h1>
			<h4>{repository.URI}</h4>
			<TextField className={classes.search}
				label="Search by branch, revision or user"
				onChange={(e) => onSearch(e.target.value)} />
			<Paper square>
				<LinearProgress className={isLoading ? '' : classes.hidden}/>
				<Table size={'medium'}>
					<TableHead>
						<TableRow className={classes.headerRow}>
							<TableCell align={'center'}
								onClick={() => onSortChange('state')}
								style={{width: '64px'}} >
								<TableSortLabel
									className={classes.sort}
									active={sortColumn === 'state'}
									direction={sortOrder} >
									State
								</TableSortLabel>
							</TableCell>
							<TableCell>Branch</TableCell>
							<TableCell>Duration</TableCell>
							<TableCell onClick={() => onSortChange('created_at')} >
								<TableSortLabel
									className={classes.sort}
									hideSortIcon={false}
									active={sortColumn === 'created_at'}
									direction={sortOrder} >
									Created
								</TableSortLabel>
							</TableCell>
							<TableCell>Started By</TableCell>
						</TableRow>
					</TableHead>
					<TableBody>
						{
							page.Jobs.length === 0 && <TableRow>
								<TableCell colSpan={5}>
									No jobs match that search criteria.
								</TableCell>
							</TableRow>
						}
						{ page.Jobs
							.map((job) => {
								return (
									<TableRow
										hover
										onClick={() => history.push(`/job/${job.ID}`)}
										key={job.ID}
										style={{cursor: 'pointer'}} >
										<TableCell align="center">{jobStatus(classes, job.State)}</TableCell>
										<TableCell align="left">
											{job.Commit.Branch.replace('refs/heads/', '')}
										</TableCell>
										<TableCell className={classes.duration} align="left">
											{duration(job)}
										</TableCell>
										<TableCell align="left">
											{ moment(job.CreatedAt).format('LLLL') }
										</TableCell>
										<TableCell align="left">{job.StartedBy}</TableCell>
									</TableRow>
								);
							})}
					</TableBody>
				</Table>
				<TablePagination
					className={classes.footer}
					rowsPerPageOptions={rowsPerPageOptions}
					component="div"
					count={page.Count}
					rowsPerPage={rowsPerPage}
					page={currentPage}
					backIconButtonProps={{
						'aria-label': 'previous page',
					}}
					nextIconButtonProps={{
						'aria-label': 'next page',
					}}
					onChangePage={(e, p) => onPageChange(p)}
					onChangeRowsPerPage={(e) =>
						onRowsPerPageChange(Number(e.target.value))
					}
				/>
			</Paper>
		</div>
	);
}
