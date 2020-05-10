import React, {useState, useEffect} from 'react';
import {
	Grid,
	TextField,
	FormControlLabel,
	InputAdornment,
	IconButton,
	Switch,
	Hidden,
	Divider,
} from '@material-ui/core';
import DeleteIcon from '@material-ui/icons/Delete';
import Visibility from '@material-ui/icons/Visibility';
import VisibilityOff from '@material-ui/icons/VisibilityOff';

import {EnvironmentVariable, EnvironmentVariableType} from '../../../services';

interface VariableProps {
	isEdit: boolean;
	variable: EnvironmentVariable;
	onSave: (variable: EnvironmentVariable) => void;
	onRemove: (name: string) => void;
}

export function VariableEntry({isEdit, variable, onSave, onRemove}: VariableProps) {
	const [value, setValue] = useState((variable && variable.Value) || '');
	const [showPassword, setShowPassword] = useState(false);
	const [sensitive, setSensitive] = useState((variable && variable.Type === EnvironmentVariableType.Password) ? true : false);

	useEffect(() => {
		setValue((variable && variable.Value) || '');
	}, [variable]);

	const save = () => {
		onSave({
			Name: variable.Name,
			Value: value,
			Type: sensitive ? EnvironmentVariableType.Password : EnvironmentVariableType.Text,
		});
	};

	return <React.Fragment>
		<Grid item xs={12} md={!isEdit ? 6 : 5}>
			<TextField
				InputLabelProps={{shrink: true, required: false}}
				variant="outlined"
				value={variable.Name}
				disabled={true}
				label="Name"
				size="small"
				fullWidth/>
		</Grid>
		{
			isEdit &&
			<Grid item xs={12} md={2} xl={1}>
				<FormControlLabel
					control={<Switch color="primary" checked={sensitive} onChange={(e) => {
						setSensitive(e.target.checked);
						setShowPassword(false);
						save();
					}}/>}
					label="Secret"
				/>
			</Grid>
		}
		<Grid item xs={12} md={!isEdit ? 6 : 4} xl={!isEdit ? 6 : 5}>
			<TextField
				InputLabelProps={{shrink: true}}
				InputProps={{
					endAdornment: sensitive ?
						(<InputAdornment position="end">
							<IconButton
								onClick={() => setShowPassword(!showPassword)}
							>
								{showPassword ? <Visibility /> : <VisibilityOff />}
							</IconButton>
						</InputAdornment>) :
						(<React.Fragment />),
				}}
				onChange={(e) => {
					setValue(e.target.value);
					save();
				}}
				variant="outlined"
				type={(sensitive && !showPassword) ? 'password' : 'text'}
				value={value}
				disabled={!isEdit}
				label="Value"
				size="small"
				fullWidth
				multiline={!sensitive} />
		</Grid>
		{
			isEdit &&
			<Grid item xs={12} md={1}>
				{
					<IconButton onClick={() => onRemove(variable.Name)}>
						<DeleteIcon />
					</IconButton>
				}
			</Grid>
		}
		<Hidden mdUp>
			<Grid item xs={12} >
				<Divider/>
			</Grid>
		</Hidden>
	</React.Fragment>;
}
