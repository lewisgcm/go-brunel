import React from 'react';
import moment from 'moment';

import {
	makeStyles,
	Theme,
	createStyles,
} from '@material-ui/core/styles';

import {
	Paper,
	Table,
	TableRow,
	TablePagination,
	TableCell,
	TableBody,
	LinearProgress,
	CircularProgress,
	Icon,
	TableHead,
	TableSortLabel,
	TextField, Tooltip,
} from '@material-ui/core';

import {
	RepositoryJobPage,
	JobState, RepositoryJob,
} from '../../services';

interface Props {
    isLoading: boolean;
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
		sort: {
			'& .MuiTableSortLabel-icon': {
				opacity: 0,
				fontSize: '1.3em',
				fontWeight: 'bold',
			},
		},
		duration: {
			width: '10em',
		},
		hidden: {
			visibility: 'hidden',
		},
		headerRow: {
			'& .MuiTableCell-root': {
				fontSize: `12px`,
				opacity: 0.7,
			},
		},
		footer: {
			'fontSize': '12px !important',
			'opacity': 0.9,
			'& .MuiTypography-root': {
				fontSize: `12px`,
				opacity: 0.9,
			},
		},
		search: {
			width: '100%',
			paddingBottom: theme.spacing(3),
		},
	}),
);

function jobStatus(state: JobState): React.ReactNode {
	switch (state) {
	case JobState.Processing:
	case JobState.Waiting:
		return <CircularProgress size={24} thickness={4} />;
	case JobState.Failed:
		return <Icon color="error">error</Icon>;
	case JobState.Cancelled:
		return <Icon color="secondary">cancel</Icon>;
	case JobState.Success:
		return <Tooltip title={'Success'}><Icon color="primary"
					 style={{color: 'rgb(0, 100, 0)', position: 'relative', top: 3}} >
				check_circle
		</Icon></Tooltip>;
	}
}

function duration(job: RepositoryJob): string {
	if ( job.StartedAt && job.StoppedAt ) {
		return moment
			.duration(moment(job.StoppedAt).diff(moment(job.StartedAt))).humanize();
	}
	return 'N/A';
}

export function RepositoryJobsComponent(
	{
		isLoading,
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

	return (
		<div>
			<h1>Namepsace/Name</h1>
			<h4>https://github.com/lewisgcm/go-brunel.git</h4>
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
						{ page.Jobs
							.map((job) => {
								return (
									<TableRow
										hover
										// onClick={event => handleClick(event, row.name)}
										key={job.ID}
									>
										<TableCell align="center">{jobStatus(job.State)}</TableCell>
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
