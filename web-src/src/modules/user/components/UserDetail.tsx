import React, {useState, useEffect} from 'react';
import {Grid, Switch, FormControlLabel, Avatar, Hidden} from '@material-ui/core';

import {User, UserRole, UserService} from '../../../services';
import { useDependency } from '../../../container';

interface Props {
	user: User;
}

export function UserDetail({user}: Props) {
	const userService = useDependency(UserService);
	const [isAdmin, setIsAdmin] = useState(user.Role === UserRole.Admin);

	useEffect(() => {
		setIsAdmin(user.Role === UserRole.Admin);
	}, [user]);

	const onAdminChange = (shouldBeAdmin: boolean) => {
		userService
			.update(
				user.Username,
				{
					Role: shouldBeAdmin ? UserRole.Admin : UserRole.Reader,
				},
			).subscribe(
				() => {
					setIsAdmin(shouldBeAdmin);
				},
				() => {},
			);
	};

	return <React.Fragment>
		<h1>{user.Username}</h1>
		<Grid container>
			<Hidden xsDown>
				<Grid item xs>
					<Avatar src={user.AvatarURL} style={{width: 150, height: 150}}/>
				</Grid>
			</Hidden>
			<Grid item xs>
				<Grid container>
					{
						user.Email &&
						<Grid item xs={12}>
							<Grid container>
								<Grid item xs={4} md={2}>
									<p><b>Email: </b></p>
								</Grid>
								<Grid item xs={8} md={10}>
									<p>{user.Email}</p>
								</Grid>
							</Grid>
						</Grid>
					}
					{
						user.Name &&
						<Grid item xs={12}>
							<Grid container>
								<Grid item xs={4} md={2}>
									<p><b>Name: </b></p>
								</Grid>
								<Grid item xs={8} md={10}>
									<p>{user.Name}</p>
								</Grid>
							</Grid>
						</Grid>
					}
					<Grid item xs={12}>
						<FormControlLabel
							control={
								<Switch
									checked={isAdmin}
									onChange={(e) => onAdminChange(e.target.checked)}
									color="primary"
								/>
							}
							label="Administrator"
						/>
					</Grid>
				</Grid>
			</Grid>
			<Grid item md={12} lg={4}></Grid>
		</Grid>
	</React.Fragment>;
}

