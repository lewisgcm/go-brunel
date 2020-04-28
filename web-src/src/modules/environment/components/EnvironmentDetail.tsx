import React, {useState, useEffect} from 'react';
import {
	TextField,
	makeStyles,
	Theme,
	createStyles,
	Button,
	Grid,
	Divider,
} from '@material-ui/core';

import {VariableEntry} from './VariableEntry';
import {VariableEntryCreate} from './VariableEntryCreate';
import {Environment} from '../../../services';

interface Props {
	isEdit: boolean;
	detail: Environment;
	onSave: (detail: Environment) => void;
	onCancel: () => void;
	onEdit: () => void;
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
		subheader: {
			marginTop: 0,
			color: theme.palette.grey[600],
		},
	}),
);

export function EnvironmentDetail({detail, onSave, onCancel, isEdit, onEdit}: Props) {
	const classes = useStyles({});
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
				{!isEdit && <h1 style={{marginBottom: 0}}>{name}</h1>}
				{isEdit && <TextField
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
					fullWidth/>}
			</Grid>

			{variables.length > 0 && <Grid item xs={12}>
				<h3 className={classes.subheader}>Variables</h3>
				<Divider />
			</Grid>}

			{variables.sort((a, b) => a.Name.localeCompare(b.Name)).map((variable) => {
				return <VariableEntry
					key={variable.Name}
					isEdit={isEdit}
					variable={variable}
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
				isEdit && <React.Fragment>
					<Grid item xs={12}>
						<h3 className={classes.subheader}>Add Variable</h3>
						<Divider />
					</Grid>
					<VariableEntryCreate
						allVariables={variables}
						onSave={(variable) => {
							setVariables(variables.concat([variable]));
						}}/>
				</React.Fragment>
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
						setNameIsDirty(false);
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
