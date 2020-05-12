import React, {useState} from 'react';
import {
	Grid,
	TextField,
	FormControlLabel,
	IconButton,
	Switch,
	Hidden,
	Button,
} from '@material-ui/core';
import AddIcon from '@material-ui/icons/Add';

import {
	EnvironmentVariable,
	EnvironmentVariableType,
} from '../../../services';

interface Props {
	allVariables: EnvironmentVariable[];
	onSave: (variable: EnvironmentVariable) => void;
}

export function VariableEntryCreate({allVariables, onSave}: Props) {
	const [nameDirty, setNameDirty] = useState(false);
	const [name, setName] = useState('');
	const [value, setValue] = useState('');
	const [sensitive, setSensitive] = useState(false);

	const isNameValid = (name: string) => {
		const existing = allVariables
			.find((r) => r.Name.trim() === name.trim());

		return !existing && (name && name.trim().length > 0 && name.trim().length < 10000);
	};

	const onSaveClick = () => {
		onSave({
			Name: name,
			Value: value,
			Type: sensitive ? EnvironmentVariableType.Password : EnvironmentVariableType.Text,
		});
		setName('');
		setValue('');
		setNameDirty(false);
		setSensitive(false);
	};

	return <React.Fragment>
		<Grid item xs={12} md={5}>
			<TextField
				InputLabelProps={{shrink: true}}
				onChange={(e) => {
					setNameDirty(true);
					setName(e.target.value);
				}}
				variant="outlined"
				value={name}
				label="Variable Name"
				size="small"
				error={!isNameValid(name) && nameDirty}
				helperText={'Variable name is required and must be unique'}
				required
				fullWidth/>
		</Grid>
		<Grid item xs={12} md={4} xl={5}>
			<TextField
				InputLabelProps={{shrink: true}}
				onChange={(e) => setValue(e.target.value)}
				variant="outlined"
				value={value}
				label="Variable Value"
				size="small"
				fullWidth
				multiline/>
		</Grid>
		<Grid item md={2} xl={1} xs={10}>
			<FormControlLabel
				control={<Switch color="primary" checked={sensitive} onChange={(e) => setSensitive(e.target.checked)} />}
				label="Secret"
			/>
		</Grid>
		<Hidden smDown>
			<Grid item md={1}>
				<IconButton disabled={!isNameValid(name)} onClick={() => onSaveClick()}>
					<AddIcon />
				</IconButton>
			</Grid>
		</Hidden>
		<Hidden smUp>
			<Grid item xs={12}>
				<Button disabled={!isNameValid(name)} onClick={() => onSaveClick()} color='primary' variant='outlined' fullWidth>
					Add
				</Button>
			</Grid>
		</Hidden>
	</React.Fragment>;
}
