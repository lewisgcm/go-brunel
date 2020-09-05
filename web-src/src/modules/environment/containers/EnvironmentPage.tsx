import React, { useEffect, useState, Dispatch } from "react";
import { BehaviorSubject, merge } from "rxjs";
import { first, debounceTime, skip, tap, switchMap } from "rxjs/operators";
import { match } from "react-router";
import { connect } from "react-redux";
import { makeStyles, Theme, createStyles, Button } from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import { useHistory } from "react-router-dom";

import {
	Drawer,
	ActionTypes,
	toggleSidebar,
	SearchableList,
	SearchListState,
} from "../../layout";
import { useDependency } from "../../../container";
import {
	EnvironmentService,
	Environment,
	EnvironmentList,
	SocketService,
	EventType,
} from "../../../services";
import { EnvironmentDetail } from "../components/EnvironmentDetail";

interface Props {
	match: match<{ environmentId: string }>;
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
	})
);

export const EnvironmentPage = connect(
	null,
	mapDispatchToProps
)(
	({
		match,
		hideMobileSidebar,
	}: Props & ReturnType<typeof mapDispatchToProps>) => {
		const classes = useStyles();
		const history = useHistory();
		const socketService = useDependency(SocketService);
		const environmentService = useDependency(EnvironmentService);
		const [environments, setEnvironments] = useState<EnvironmentList[]>([]);
		const [subject] = useState(new BehaviorSubject(""));
		const [listState, setListState] = useState(SearchListState.Loaded);
		const [detail, setDetail] = useState<Environment>();
		const [isEdit, setIsEdit] = useState(false);
		const [saveError, setSaveError] = useState<string | undefined>();

		useEffect(() => {
			const subscription = socketService
				.events(EventType.EnvironmentCreated)
				.subscribe(() => {
					subject.next(subject.getValue());
				});

			return () => {
				subscription.unsubscribe();
			};
		}, [socketService, subject]);

		useEffect(() => {
			setSaveError("");
		}, [isEdit]);

		useEffect(() => {
			setIsEdit(false);
		}, [detail]);

		useEffect(() => {
			if (match.params.environmentId) {
				environmentService.get(match.params.environmentId).subscribe(
					(environment) => {
						setDetail(environment);
					},
					() => {
						setDetail(undefined);
					}
				);
			}
		}, [environmentService, match.params.environmentId, history]);

		useEffect(() => {
			const subscription = merge(
				subject.pipe(first()),
				subject.pipe(skip(1), debounceTime(200))
			)
				.pipe(
					tap(() => setListState(SearchListState.Loading)),
					switchMap((term) => environmentService.list(term)),
					tap(() => setListState(SearchListState.Loaded))
				)
				.subscribe(
					(items) => {
						setEnvironments(items);
					},
					() => {
						setListState(SearchListState.Error);
					}
				);

			return () => {
				subscription.unsubscribe();
			};
		}, [environmentService, subject]);

		useEffect(() => {
			if (environments.length && !match.params.environmentId && !isEdit) {
				history.push(`/environment/${environments[0].ID}`);
			}
		}, [environments, match.params.environmentId, history, isEdit]);

		const onSave = (environment: Environment) => {
			environmentService.save(environment).subscribe(
				(newEnvironment) => {
					if (!environment.ID) {
						setEnvironments(
							environments
								.filter((e) => e.ID !== newEnvironment.ID)
								.concat([newEnvironment])
						);
						history.replace(`/environment/${newEnvironment.ID}`);
					} else {
						setDetail(newEnvironment);
						setEnvironments(
							environments.map((e) =>
								e.ID === newEnvironment.ID ? newEnvironment : e
							)
						);
					}
				},
				(error) => {
					setSaveError(
						`Failed to save environment: ${error.message}`
					);
				}
			);
		};

		const AddButton = () => (
			<Button
				onClick={() => {
					hideMobileSidebar();
					setDetail({
						ID: "",
						Name: "",
						Variables: [],
					});
				}}
				className={classes.addButton}
				variant="contained"
				color="primary"
				fullWidth
				startIcon={<AddIcon />}
			>
				Add Environment
			</Button>
		);

		const sidebar = () => (
			<SearchableList
				state={listState}
				errorPlaceholder="Error loading environments."
				emptyPlaceholder="No environments found."
				searchPlaceholder="Search for environments"
				items={environments}
				render={(item) => ({
					selected: item.ID === match.params.environmentId,
					text: item.Name,
					key: item.ID,
				})}
				onSearch={(term) => subject.next(term)}
				onClick={(item) => {
					hideMobileSidebar();
					history.push(`/environment/${item.ID}`);
				}}
				children={<AddButton />}
			/>
		);

		const content = () => {
			return detail !== undefined ? (
				<EnvironmentDetail
					isEdit={detail && !detail.ID ? true : isEdit}
					error={saveError}
					onEdit={() => setIsEdit(true)}
					detail={{
						...detail,
					}}
					onCancel={() => {
						if (!detail.ID) {
							setDetail(undefined);
							history.push(`/environment/`);
						} else {
							setIsEdit(false);
						}
					}}
					onSave={onSave}
				/>
			) : (
				<React.Fragment />
			);
		};

		return <Drawer sidebar={sidebar} content={content} />;
	}
);
