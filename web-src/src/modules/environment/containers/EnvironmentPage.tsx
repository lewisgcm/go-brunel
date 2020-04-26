import React, {useEffect, useState} from 'react';
import {BehaviorSubject, merge} from 'rxjs';
import {first, debounceTime, skip, tap, switchMap, distinctUntilChanged} from 'rxjs/operators';
import {match} from 'react-router';

import {Drawer} from '../../layout';
import {EnvironmentListComponent} from '../components/EnvironmentListComponent';
import {useDependency} from '../../../container';
import {EnvironmentService, Environment, EnvironmentList} from '../../../services';
import {EnvironmentDetail} from '../components/EnvironmentDetail';

interface Props {
	match: match<{environmentId: string}>;
}

export const EnvironmentPage = ({match}: Props) => {
	const environmentService = useDependency(EnvironmentService);
	const [environments, setEnvironments] = useState<EnvironmentList[]>([]);
	const [subject] = useState(new BehaviorSubject(''));
	const [isLoading, setLoading] = useState(false);
	const [selectedEnvironmentId, setSelectedEnvironmentId] = useState<string | undefined>(match.params.environmentId);
	const [detail, setDetail] = useState<Environment>();
	const [isEdit, setIsEdit] = useState(false);

	useEffect(() => {
		const subscription = merge(
			subject.pipe(first()),
			subject.pipe(
				distinctUntilChanged(),
				skip(1),
				debounceTime(200),
			),
		).pipe(
			tap(() => setLoading(true)),
			switchMap((term) => environmentService.list(term)),
			tap(() => setLoading(false)),
		).subscribe(
			(items) => {
				setEnvironments(items);
				if (items.length > 0) {
					setSelectedEnvironmentId(items[0].ID);
				}
			},
		);

		return () => {
			subscription.unsubscribe();
		};
	}, [environmentService, subject]);

	useEffect(() => {
		if (selectedEnvironmentId) {
			environmentService
				.get(selectedEnvironmentId)
				.subscribe((environment) => {
					setDetail(environment);
				});
		}
	}, [environmentService, selectedEnvironmentId]);

	return <Drawer sidebar={() => <EnvironmentListComponent
		isLoading={isLoading}
		environments={environments}
		onClick={(id) => {
			setIsEdit(false);
			setSelectedEnvironmentId(id);
		}}
		onSearch={(term) => subject.next(term)}
		selectedEnvironmentId={selectedEnvironmentId}
		onAdd={() => {
			setSelectedEnvironmentId(undefined);
			setIsEdit(true);
			setDetail({
				ID: '',
				Name: '',
				Variables: [],
			});
		}}
	/>}
	content={() => {
		return detail !== undefined ?
			<EnvironmentDetail
				isEdit={isEdit}
				onEdit={() => setIsEdit(true)}
				detail={Object.assign({}, detail)}
				onCancel={() => {
					setIsEdit(false);
					setSelectedEnvironmentId(selectedEnvironmentId ? selectedEnvironmentId : environments.length ? environments[0].ID : undefined);
				}}
				onSave={(environment) => {
					environmentService
						.save(environment)
						.subscribe((newEnvironment) => {
							setEnvironments(
								environments.filter((e) => e.ID !== newEnvironment.ID).concat([
									newEnvironment,
								]),
							);
							setIsEdit(false);
							setDetail(newEnvironment);
						});
				}} /> :
			<React.Fragment/>;
	}} />;
};
