import React, {ReactNode, useState, useEffect} from 'react';
import {Link, match} from 'react-router-dom';
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

import {Observable, BehaviorSubject, merge} from 'rxjs';
import {first, distinctUntilChanged, skip, debounceTime, tap, switchMap} from 'rxjs/operators';

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

interface Props<T> {
	emptyPlaceholder: string;
	searchPlaceholder: string;
	errorPlaceholder: string;
	provider: (term: string) => Observable<T[]>;
	render: (item: T) => RenderItem;
	onClick: (item: T) => void;
	children?: ReactNode;
}

export function SearchableList<T>({
	emptyPlaceholder,
	searchPlaceholder,
	errorPlaceholder,
	provider,
	onClick,
	children,
	render,
}: Props<T>) {
	const classes = useStyles({});
	const [subject] = useState(new BehaviorSubject(''));
	const [isLoading, setLoading] = useState(false);
	const [items, setItems] = useState<T[]>([]);
	const [hasError, setHasError] = useState<boolean>(false);

	useEffect(
		() => {
			const subscription = merge(
				subject.pipe(first()),
				subject.pipe(
					distinctUntilChanged(),
					skip(1),
					debounceTime(200),
				),
			).pipe(
				tap(() => setLoading(true)),
				switchMap(provider),
				tap(() => setLoading(false)),
			).subscribe(
				(items) => {
					setHasError(false);
					setItems(items);
				},
				() => {
					setHasError(true);
					setItems([]);
					setLoading(false);
				},
			);

			return () => {
				subscription.unsubscribe();
			};
		},
		[provider, subject],
	);

	return <List className={classes.list}>
		{children}

		<TextField className={classes.input}
			label={searchPlaceholder}
			onChange={(e) => subject.next(e.target.value)} />

		<LinearProgress className={isLoading ? '' : classes.hidden} />

		{!hasError && items.map((i) => ({r: render(i), item: i})).map(
			(item) => {
				return <ListItem
					className={`${classes.listItem} ${item.r.selected ? classes.selectedItem : ''}`}
					button
					// component={Link}
					key={item.r.key}
					// to={`/user/${user.Username}`}
					onClick={() => onClick(item.item)} >
					<ListItemText>{item.r.text}</ListItemText>
				</ListItem>;
			},
		)}

		{!hasError && items.length === 0 && <Typography className={classes.empty}>
			{emptyPlaceholder}
		</Typography>}

		{hasError && <Alert severity="error">
			{errorPlaceholder}
		</Alert>}
	</List>;
}
