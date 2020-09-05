import React, { useState, useEffect } from "react";
import {
	Grid,
	TextField,
	FormControlLabel,
	InputAdornment,
	IconButton,
	Switch,
	Hidden,
	Divider,
	Button,
} from "@material-ui/core";
import DeleteIcon from "@material-ui/icons/Delete";
import Visibility from "@material-ui/icons/Visibility";
import VisibilityOff from "@material-ui/icons/VisibilityOff";

import {
	EnvironmentVariable,
	EnvironmentVariableType,
} from "../../../services";

interface VariableProps {
	isEdit: boolean;
	variable: EnvironmentVariable;
	onSave: (variable: EnvironmentVariable) => void;
	onRemove: (name: string) => void;
}

export function VariableEntry({
	isEdit,
	variable,
	onSave,
	onRemove,
}: VariableProps) {
	const [showPassword, setShowPassword] = useState(false);
	const isSensitive =
		variable && variable.Type === EnvironmentVariableType.Password
			? true
			: false;

	useEffect(() => {
		setShowPassword(false);
	}, [variable]);

	const save = (value: string, sensitive: boolean) => {
		onSave({
			Name: variable.Name,
			Value: value,
			Type: sensitive
				? EnvironmentVariableType.Password
				: EnvironmentVariableType.Text,
		});
	};

	return (
		<React.Fragment>
			<Grid item xs={12} md={!isEdit ? 6 : 5}>
				<TextField
					InputLabelProps={{ shrink: true, required: false }}
					variant="outlined"
					value={variable.Name}
					disabled={true}
					label="Name"
					size="small"
					fullWidth
				/>
			</Grid>
			<Grid item xs={12} md={!isEdit ? 6 : 4} xl={!isEdit ? 6 : 5}>
				<TextField
					InputLabelProps={{ shrink: true }}
					InputProps={{
						endAdornment: isSensitive ? (
							<InputAdornment position="end">
								<IconButton
									onClick={() =>
										setShowPassword(!showPassword)
									}
								>
									{showPassword ? (
										<Visibility />
									) : (
										<VisibilityOff />
									)}
								</IconButton>
							</InputAdornment>
						) : (
							<React.Fragment />
						),
					}}
					onChange={(e) => {
						save(e.target.value, isSensitive);
					}}
					variant="outlined"
					type={isSensitive && !showPassword ? "password" : "text"}
					value={variable.Value}
					disabled={!isEdit}
					label="Value"
					size="small"
					fullWidth
					multiline={!isSensitive}
				/>
			</Grid>
			{isEdit && (
				<Grid item xs={12} md={2} xl={1}>
					<FormControlLabel
						control={
							<Switch
								color="primary"
								checked={isSensitive}
								onChange={(e) => {
									setShowPassword(false);
									save(variable.Value, e.target.checked);
								}}
							/>
						}
						label="Secret"
					/>
				</Grid>
			)}
			{isEdit && (
				<React.Fragment>
					<Hidden smDown>
						<Grid item xs={2} md={1}>
							{
								<IconButton
									onClick={() => onRemove(variable.Name)}
								>
									<DeleteIcon />
								</IconButton>
							}
						</Grid>
					</Hidden>
					<Hidden mdUp>
						<Grid item xs={12}>
							{
								<Button
									onClick={() => onRemove(variable.Name)}
									color="secondary"
									variant="outlined"
									fullWidth
								>
									Remove
								</Button>
							}
						</Grid>
					</Hidden>
				</React.Fragment>
			)}
			<Hidden mdUp>
				<Grid item xs={12}>
					<Divider />
				</Grid>
			</Hidden>
		</React.Fragment>
	);
}
