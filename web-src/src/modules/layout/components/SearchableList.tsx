import React, {ReactNode} from 'react';
import {makeStyles, Theme, createStyles} from '@material-ui/core/styles';
import {
	List,
	ListItem,
	ListItemText,
	Typography,
	TextField,
	LinearProgress,
} from '@material-ui/core';
import {Alert} from '@material-ui/lab';

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		list: {
			padding: theme.spacing(2),
		},
		listItem: {
			borderBottom: `1px solid ${theme.palette.grey[300]}`,
		},
		input: {
			width: '100%',
		},
		hidden: {
			visibility: 'hidden',
		},
		empty: {
			textAlign: 'center',
			paddingTop: theme.spacing(1),
		},
		selectedItem: {
			backgroundColor: theme.palette.grey[300],
		},
	}),
);

export interface RenderItem {
	selected: boolean;
	text: string;
	key: string;
}

export enum SearchListState {
	Loading = 0,
	Loaded = 1,
	Error = 2,
}

interface Props<T> {
	state: SearchListState;
	emptyPlaceholder: string;
	searchPlaceholder: string;
	errorPlaceholder: string;
	items: T[];
	render: (item: T) => RenderItem;
	children?: ReactNode;
	onSearch: (term: string) => void;
	onClick: (item: T) => void;
}

export function SearchableList<T>({
	state,
	emptyPlaceholder,
	searchPlaceholder,
	errorPlaceholder,
	items,
	render,
	children,
	onSearch,
	onClick,
}: Props<T>) {
	const classes = useStyles({});

	return <List className={classes.list}>
		{children}

		<TextField className={classes.input}
			label={searchPlaceholder}
			onChange={(e) => onSearch(e.target.value)} />

		<LinearProgress className={state === SearchListState.Loading ? '' : classes.hidden} />

		{state === SearchListState.Loaded && items.map((i) => ({r: render(i), item: i})).map(
			(item) => {
				return <ListItem
					className={`${classes.listItem} ${item.r.selected ? classes.selectedItem : ''}`}
					button
					key={item.r.key}
					onClick={() => onClick(item.item)} >
					<ListItemText>{item.r.text}</ListItemText>
				</ListItem>;
			},
		)}

		{state === SearchListState.Loaded && items.length === 0 && <Typography className={classes.empty}>
			{emptyPlaceholder}
		</Typography>}

		{state === SearchListState.Error && <Alert severity="error">
			{errorPlaceholder}
		</Alert>}
	</List>;
}
