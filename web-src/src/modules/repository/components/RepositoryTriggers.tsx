import React, { useState } from 'react';
import {Grid, Select, MenuItem, TextField, FormControl, InputLabel} from '@material-ui/core';
import RepositoryEnvironmentSelection from './RepositoryEnvironmentSelection';

export function RepositoryTriggers() {
	const [referenceType, setReferenceType] = useState(0);

	return <Grid container justify="space-between" spacing={3}>
		<Grid item xs={4}>
			<FormControl fullWidth>
				<InputLabel id="demo-simple-select-label">Reference Type</InputLabel>
				<Select
					labelId="demo-simple-select-label"
					id="demo-simple-select"
					value={referenceType}
					onChange={(e) => setReferenceType(e.target.value as number)}
				>
					<MenuItem value={0}>Branch</MenuItem>
					<MenuItem value={1}>Tag</MenuItem>
				</Select>
			</FormControl>
		</Grid>
		<Grid item xs={4}>
			<TextField label="Reference" fullWidth></TextField>
		</Grid>
		<Grid item xs={4}>
			<RepositoryEnvironmentSelection />
		</Grid>
	</Grid>;
}
