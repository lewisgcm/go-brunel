import React, {useState} from 'react';
import {
	Grid,
	TextField,
	FormControlLabel,
	IconButton,
	Switch,
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

	return <React.Fragment>
		<Grid item xs={5}>
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
		<Grid item xs={5}>
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
		<Grid item xs={1}>
			<FormControlLabel
				control={<Switch color="primary" checked={sensitive} onChange={(e) => setSensitive(e.target.checked)} />}
				label="Secret"
			/>
		</Grid>
		<Grid item xs={1}>
			<IconButton disabled={!isNameValid(name)} onClick={() => {
				onSave({
					Name: name,
					Value: value,
					Type: sensitive ? EnvironmentVariableType.Password : EnvironmentVariableType.Text,
				});
				setName('');
				setValue('');
				setNameDirty(false);
				setSensitive(false);
			}}>
				<AddIcon />
			</IconButton>
		</Grid>
	</React.Fragment>;
}
