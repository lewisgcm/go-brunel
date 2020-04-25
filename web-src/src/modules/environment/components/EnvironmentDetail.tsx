import React, {useState, useEffect} from 'react';
import {TextField, makeStyles, Theme, createStyles, Button, Grid, FormControlLabel, Switch, ButtonGroup} from '@material-ui/core';

import {Environment, EnvironmentVariable, EnvironmentVariableType} from '../../../services';

interface Props {
	isEdit: boolean;
	detail: Environment;
	onSave: (detail: Environment) => void;
	onCancel: () => void;
	onEdit: () => void;
}

interface VariableProps {
	canEdit: boolean;
	variable: EnvironmentVariable | undefined;
	allVariables: EnvironmentVariable[];
	onSave: (variable: EnvironmentVariable) => void;
	onRemove: (name: string) => void;
}

enum VariableState {
	Read,
	Create,
	Edit
}

function VariableEntry({canEdit, variable, allVariables, onSave, onRemove}: VariableProps) {
	const [state, setState] = useState(!variable ? VariableState.Create : VariableState.Read);
	const [nameDirty, setNameDirty] = useState(false);
	const [name, setName] = useState((variable && variable.Name) || '');
	const [value, setValue] = useState((variable && variable.Value) || '');
	const [sensitive, setSensitive] = useState((variable && variable.Type === EnvironmentVariableType.Password) ? true : false);

	const isNameValid = (name: string) => {
		const existing = allVariables
			.find((r) => r.Name.trim() === name.trim());

		return state !== VariableState.Create || (!existing && (name && name.trim().length > 0 && name.trim().length < 10000));
	};

	return <React.Fragment>
		<Grid item xs={!canEdit ? 6 : 4}>
			<TextField
				onChange={(e) => {
					setNameDirty(true);
					setName(e.target.value);
				}}
				value={name}
				disabled={state === VariableState.Read || state === VariableState.Edit}
				label="Variable Name"
				error={!isNameValid(name) && nameDirty}
				helperText={'Variable name is required and must be unique'}
				required
				fullWidth/>
		</Grid>
		{ canEdit && <Grid item xs={1}>
			<FormControlLabel
				disabled={state === VariableState.Read}
				style={{height: '100%'}}
				control={<Switch color="primary" checked={sensitive} onChange={(e) => setSensitive(e.target.checked)} />}
				label="Secret"
			/>
		</Grid> }
		<Grid item xs={!canEdit ? 6 : 5}>
			<TextField
				onChange={(e) => setValue(e.target.value)}
				type={sensitive ? 'password' : 'text'}
				value={value} disabled={state === VariableState.Read} label="Variable Value" fullWidth multiline={!sensitive}/>
		</Grid>
		{ canEdit && <Grid item xs={2}>
			<ButtonGroup fullWidth>
				{
					state === VariableState.Create && <Button disabled={!isNameValid(name)} onClick={() => {
						onSave({
							Name: name,
							Value: value,
							Type: sensitive ? EnvironmentVariableType.Password : EnvironmentVariableType.Text,
						});
						setName('');
						setValue('');
						setSensitive(false);
						setNameDirty(false);
					}}
					variant="outlined"
					color="primary"
					fullWidth>
						Add
					</Button>
				}
				{
					state === VariableState.Edit && <Button disabled={!isNameValid(name)} onClick={() => {
						setState(VariableState.Read);
						onSave({
							Name: name,
							Value: value,
							Type: sensitive ? EnvironmentVariableType.Password : EnvironmentVariableType.Text,
						});
					}}
					variant="outlined"
					fullWidth>
						Save
					</Button>
				}
				{
					state === VariableState.Edit && <Button onClick={() => {
						setValue((variable && variable.Value) || '');
						setSensitive((variable && variable.Type === EnvironmentVariableType.Password) ? true : false);
						setState(VariableState.Read);
					}}
					variant="outlined"
					fullWidth>
						Cancel
					</Button>
				}
				{
					canEdit && state === VariableState.Read && <Button onClick={() => {
						setState(VariableState.Edit);
					}}
					variant="outlined"
					fullWidth>
						Edit
					</Button>
				}
				{
					canEdit && state === VariableState.Read && <Button onClick={() => onRemove(name)}
						variant="outlined"
						fullWidth>
						Remove
					</Button>
				}
			</ButtonGroup>
		</Grid>}
	</React.Fragment>;
}

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		titleEdit: {
			marginTop: theme.spacing(2),
		},
		environmentVariableTitle: {
			margin: 0,
		},
		sensitiveSwitch: {
			height: '100%',
		},
	}),
);

export function EnvironmentDetail({detail, onSave, onCancel, isEdit, onEdit}: Props) {
	const classes = useStyles();
	const [isNameDirty, setNameIsDirty] = useState(false);
	const [name, setName] = useState(detail.Name);
	const [variables, setVariables] = useState(detail.Variables);

	const isNameValid = (name: string) => {
		return name && name.trim().length > 0 && name.trim().length < 100;
	};

	useEffect(() => {
		setName(detail.Name);
		setVariables(detail.Variables);
	}, [detail]);

	return <div>
		<Grid container spacing={3}>
			<Grid item xs={12}>
				<TextField
					required
					onChange={(e) => {
						setName(e.target.value);
						setNameIsDirty(true);
					}}
					helperText={'Environment name is required and must be unique'}
					className={classes.titleEdit}
					error={!isNameValid(name) && isNameDirty}
					value={name}
					disabled={!isEdit}
					label="Environment Name"
					fullWidth/>
			</Grid>

			{variables.sort((a, b) => a.Name.localeCompare(b.Name)).map((variable) => {
				return <VariableEntry
					key={variable.Name}
					canEdit={isEdit}
					variable={variable}
					allVariables={variables}
					onSave={(newVariable) => {
						setVariables(
							variables
								.filter((v) => v.Name !== variable.Name)
								.concat([newVariable]),
						);
					}}
					onRemove={(name) => {
						setVariables(variables.filter((v) => v.Name !== name));
					}} />;
			})}

			{
				isEdit && <VariableEntry
					canEdit={isEdit}
					variable={undefined}
					allVariables={variables}
					onSave={(variable) => {
						setVariables(variables.concat([variable]));
					}}
					onRemove={() => {/* Not required */}} />
			}

			{isEdit && <React.Fragment>
				<Grid item xs={8}></Grid>
				<Grid item xs={2}>
					<Button variant="contained" fullWidth onClick={() => onCancel()}>
						Cancel
					</Button>
				</Grid>
				<Grid item xs={2}>
					<Button variant="contained" color="primary" disabled={!isNameValid(name)} fullWidth onClick={() => {
						onSave({
							ID: detail ? detail.ID : '',
							Name: name,
							Variables: variables,
						});
					}}>Save</Button>
				</Grid>
			</React.Fragment>}

			{!isEdit && <React.Fragment>
				<Grid item xs={10}></Grid>
				<Grid item xs={2}>
					<Button color="primary" variant="contained" fullWidth onClick={() => onEdit()}>
						Edit
					</Button>
				</Grid>
			</React.Fragment>}
		</Grid>
	</div>;
}
