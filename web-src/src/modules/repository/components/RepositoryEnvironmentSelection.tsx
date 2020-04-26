import React, {useState, useEffect} from 'react';
import TextField from '@material-ui/core/TextField';
import Autocomplete from '@material-ui/lab/Autocomplete';
import CircularProgress from '@material-ui/core/CircularProgress';
import {EnvironmentList, EnvironmentService} from '../../../services';
import {useDependency} from '../../../container';

export default function RepositoryEnvironmentSelection() {
	const environmentService = useDependency(EnvironmentService);
	const [open, setOpen] = useState(false);
	const [options, setOptions] = useState<EnvironmentList[]>([]);
	const loading = open && options.length === 0;

	useEffect(() => {
		let active = true;

		if (!loading) {
			return undefined;
		}

		environmentService
			.list('')
			.subscribe((items) => {
				if (active) {
					setOptions(items);
				}
			});

		return () => {
			active = false;
		};
	}, [environmentService, loading]);

	useEffect(() => {
		if (!open) {
			setOptions([]);
		}
	}, [open]);

	return (
		<Autocomplete
			open={open}
			onOpen={() => {
				setOpen(true);
			}}
			onClose={() => {
				setOpen(false);
			}}
			getOptionSelected={(option, value) => option.Name === value.Name}
			getOptionLabel={(option) => option.Name}
			options={options}
			loading={loading}
			renderInput={(params) => (
				<TextField
					{...params}
					label="Environment"
					InputProps={{
						...params.InputProps,
						endAdornment: (
							<React.Fragment>
								{loading ? <CircularProgress color="inherit" size={20} /> : null}
								{params.InputProps.endAdornment}
							</React.Fragment>
						),
					}}
				/>
			)}
		/>
	);
}
