import React, {useState, useEffect} from 'react';
import {TextField, CircularProgress} from '@material-ui/core';
import {Autocomplete} from '@material-ui/lab';

import {EnvironmentList, EnvironmentService, Environment} from '../../../services';
import {useDependency} from '../../../container';

interface RepositoryEnvironmentSelectionProps {
	value: string | undefined;
	onChange: (value: string | undefined) => void;
	required?: boolean;
	error?: boolean;
	helperText?: string;
}

export default function RepositoryEnvironmentSelection({
	value,
	onChange,
	required,
	error,
	helperText,
}: RepositoryEnvironmentSelectionProps) {
	const environmentService = useDependency(EnvironmentService);
	const [open, setOpen] = useState(false);
	const [options, setOptions] = useState<EnvironmentList[]>([]);
	const [selectedEnvironment, setSelectedEnvironment] = useState<EnvironmentList | null>(null);
	const loading = (open && options.length === 0) || (!!value && selectedEnvironment === null);

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
				if (value) {
					const item = items.find((i) => i.ID === value);
					if (item) {
						setSelectedEnvironment(item);
					}
				}
			});

		return () => {
			active = false;
		};
	}, [environmentService, value, selectedEnvironment, loading]);

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
			autoHighlight
			value={selectedEnvironment}
			onChange={(event: any, newValue: EnvironmentList | null) => {
				setSelectedEnvironment(newValue);
				onChange(newValue ? newValue.ID : undefined);
			}}
			getOptionSelected={(option, value) => option.ID === value.ID}
			getOptionLabel={(option) => option.Name}
			options={options}
			loading={loading}
			renderInput={(params) => (
				<TextField
					{...params}
					error={error}
					helperText={helperText}
					required={required}
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
