import React, {useEffect, useState, Dispatch} from 'react';
import {BehaviorSubject, merge} from 'rxjs';
import {first, debounceTime, skip, tap, switchMap, distinctUntilChanged} from 'rxjs/operators';
import {match} from 'react-router';
import {connect} from 'react-redux';
import {makeStyles, Theme, createStyles, Button} from '@material-ui/core';
import AddIcon from '@material-ui/icons/Add';
import {useHistory} from 'react-router-dom';

import {Drawer, ActionTypes, toggleSidebar, SearchableList, SearchListState} from '../../layout';
import {useDependency} from '../../../container';
import {EnvironmentService, Environment, EnvironmentList} from '../../../services';
import {EnvironmentDetail} from '../components/EnvironmentDetail';

interface Props {
	match: match<{environmentId: string}>;
}

function mapDispatchToProps(dispatch: Dispatch<ActionTypes>) {
	return {
		hideMobileSidebar: () => {
			dispatch(toggleSidebar(false));
		},
	};
}

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		addButton: {
			marginBottom: theme.spacing(1),
		},
	}),
);

export const EnvironmentPage = connect(
	null,
	mapDispatchToProps,
)(({match, hideMobileSidebar}: Props & ReturnType<typeof mapDispatchToProps>) => {
	const classes = useStyles();
	const history = useHistory();
	const environmentService = useDependency(EnvironmentService);
	const [environments, setEnvironments] = useState<EnvironmentList[]>([]);
	const [subject] = useState(new BehaviorSubject(''));
	const [listState, setListState] = useState(SearchListState.Loaded);
	const [detail, setDetail] = useState<Environment>();
	const [isEdit, setIsEdit] = useState(false);
	const [saveError, setSaveError] = useState<string | undefined>();

	useEffect(
		() => {
			if (match.params.environmentId) {
				environmentService
					.get(match.params.environmentId)
					.subscribe((environment) => {
						setDetail(environment);
					});
			}
		},
		[environmentService, match.params.environmentId],
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
				switchMap((term) => environmentService.list(term)),
				tap(() => setListState(SearchListState.Loaded)),
			).subscribe(
				(items) => {
					setEnvironments(items);
				},
				() => {
					setListState(SearchListState.Error);
				},
			);

			return () => {
				subscription.unsubscribe();
			};
		},
		[environmentService, subject],
	);

	useEffect(
		() => {
			if (environments.length && !match.params.environmentId && !isEdit) {
				history.push(`/environment/${environments[0].ID}`);
			}
		},
		[environments, match.params.environmentId, history, isEdit],
	);

	const onSave = (environment: Environment) => {
		environmentService
			.save(environment)
			.subscribe(
				(newEnvironment) => {
					setEnvironments(
						environments.filter((e) => e.ID !== newEnvironment.ID).concat([
							newEnvironment,
						]),
					);
					setIsEdit(false);
					history.push(`/environment/${environment.ID}`);
					setSaveError(undefined);
				},
				(error) => {
					let message = 'Failed to save environment';
					if (error && error.Error) {
						message = `${message}: ${error.Error}`;
					}
					setSaveError(message);
				},
			);
	};

	const AddButton = () => (<Button
		onClick={() => {
			history.push(`/environment/`);
			hideMobileSidebar();
			setIsEdit(true);
			setDetail({
				ID: '',
				Name: '',
				Variables: [],
			});
		}}
		className={classes.addButton}
		variant="contained"
		color="primary"
		fullWidth
		startIcon={<AddIcon />}>
			Add Environment
	</Button>);

	const sidebar = () => <SearchableList
		state={listState}
		errorPlaceholder='Error loading environments.'
		emptyPlaceholder='No environments found.'
		searchPlaceholder='Search for environments'
		items={environments}
		render={(item) => ({
			selected: item.ID === match.params.environmentId,
			text: item.Name,
			key: item.ID,
		})}
		onSearch={(term) => subject.next(term)}
		onClick={(item) => {
			hideMobileSidebar();
			setIsEdit(false);
			history.push(`/environment/${item.ID}`);
		}}
		children={<AddButton/>}/>;

	const content = () => {
		return detail !== undefined ?
			<EnvironmentDetail
				isEdit={isEdit}
				error={saveError}
				onEdit={() => setIsEdit(true)}
				detail={{
					...detail,
				}}
				onCancel={() => {
					setIsEdit(false);
				}}
				onSave={onSave} /> :
			<React.Fragment/>;
	};

	return <Drawer sidebar={sidebar} content={content} />;
});
