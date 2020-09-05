import React, { useState, useEffect } from "react";
import {
	Grid,
	Select,
	MenuItem,
	TextField,
	FormControl,
	InputLabel,
	IconButton,
	InputAdornment,
	Hidden,
	Button,
} from "@material-ui/core";
import DeleteIcon from "@material-ui/icons/Delete";

import RepositoryEnvironmentSelection from "./RepositoryEnvironmentSelection";
import { RepositoryTrigger, RepositoryTriggerType } from "../../../services";

interface TriggerProps {
	trigger: RepositoryTrigger;
	onRemove: () => void;
	onChange: (trigger: RepositoryTrigger) => void;
	isValid: boolean;
}

export function Trigger({
	trigger,
	onRemove,
	onChange,
	isValid,
}: TriggerProps) {
	const [reference, setReference] = useState(trigger.Pattern);
	const [referenceType, setReferenceType] = useState(trigger.Type);
	const [environmentId, setEnvironmentId] = useState<string | undefined>(
		trigger.EnvironmentID
	);

	useEffect(() => {
		setReference(trigger.Pattern);
		setReferenceType(trigger.Type);
		setEnvironmentId(trigger.EnvironmentID);
	}, [trigger]);

	return (
		<React.Fragment>
			<Grid item xs={12} md={3}>
				<FormControl fullWidth>
					<InputLabel>Type</InputLabel>
					<Select
						value={referenceType}
						onChange={(e) => {
							setReferenceType(e.target.value as number);
							onChange({
								Type: e.target.value as number,
								Pattern: reference,
								EnvironmentID: environmentId,
							});
						}}
					>
						<MenuItem value={RepositoryTriggerType.Branch}>
							Branch
						</MenuItem>
						<MenuItem value={RepositoryTriggerType.Tag}>
							Tag
						</MenuItem>
					</Select>
				</FormControl>
			</Grid>
			<Grid item xs={12} md={4}>
				<TextField
					InputProps={{
						startAdornment: (
							<InputAdornment position="end">/</InputAdornment>
						),
						endAdornment: (
							<InputAdornment position="end">/</InputAdornment>
						),
					}}
					label="Pattern"
					required
					value={reference}
					fullWidth
					error={!isValid}
					helperText={
						!isValid
							? "You must enter a pattern for matching branches or tags"
							: undefined
					}
					onChange={(e) => {
						setReference(e.target.value);
						onChange({
							Type: referenceType,
							Pattern: e.target.value,
							EnvironmentID: environmentId,
						});
					}}
				/>
			</Grid>
			<Grid item xs={12} md={4}>
				<RepositoryEnvironmentSelection
					value={environmentId}
					onChange={(e) => {
						setEnvironmentId(e);
						onChange({
							Type: referenceType,
							Pattern: reference,
							EnvironmentID: e,
						});
					}}
				/>
			</Grid>
			<Hidden mdUp>
				<Grid item xs={12}>
					<Button
						color="secondary"
						variant="outlined"
						fullWidth
						onClick={() => onRemove()}
					>
						Remove
					</Button>
				</Grid>
			</Hidden>
			<Hidden smDown>
				<Grid item xs={1}>
					<IconButton onClick={() => onRemove()}>
						<DeleteIcon />
					</IconButton>
				</Grid>
			</Hidden>
		</React.Fragment>
	);
}
