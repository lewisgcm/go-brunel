import React, {useState} from 'react';
import {connect} from 'react-redux';
import {useHistory} from 'react-router';
import {Button, Container, Divider, Typography} from '@material-ui/core';
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';
import {FaGitlab, FaGithub} from 'react-icons/fa';

import {withDependency} from '../../container';
import {AuthService} from '../../services';
import {setAuthenticated, OAuthPopup} from '../layout';

interface Props {
	setLoggedIn: () => void;
}

interface Dependencies {
	authService: AuthService;
}

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		container: {
			height: '100%',
			display: 'flex',
			alignItems: 'center',
			justifyContent: 'center',
		},
		login: {
			textAlign: 'center',
		},
		gitlab: {
			background: 'linear-gradient(45deg, #e24228 30%, #fca326 90%)',
			display: 'block',
			marginBottom: theme.spacing(1),
		},
		github: {
			background: 'linear-gradient(45deg, #272727 30%, #464646 90%)',
			display: 'block',
			color: 'white',
		},
		divider: {
			margin: `${theme.spacing(1)}px 0px`,
		},
	}),
);

export const Login = connect(
	null,
	(dispatch) => ({
		setLoggedIn: () => dispatch(setAuthenticated(true)),
	}),
)(withDependency<Props, Dependencies>(
	(container) => ({
		authService: container.get(AuthService),
	}),
)(
	({authService, setLoggedIn}) => {
		const classes = useStyles();
		const history = useHistory();
		const [isOpen, setOpen] = useState(false);
		const [provider, setProvider] = useState<'github' | 'gitlab'>('github');

		return <Container className={classes.container}>
			<div className={classes.login}>
				<Typography variant={'h5'}>
					Brunel CI
				</Typography>

				<Divider className={classes.divider}/>

				<Button className={classes.gitlab} onClick={() => {
					setProvider('gitlab');
					setOpen(true);
				}}>
					<FaGitlab style={{paddingRight: '10px', verticalAlign: 'middle'}} />
					Login with GitLab
				</Button>

				<Button className={classes.github} onClick={() => {
					setProvider('github');
					setOpen(true);
				}}>
					<FaGithub style={{paddingRight: '10px', verticalAlign: 'middle'}} />
					Login with GitHub
				</Button>
			</div>
			<OAuthPopup
				isOpen={isOpen}
				provider={provider}
				onDone={(e) => {
					authService.setAuthentication(e);
					setLoggedIn();
					history.push('/repository/');
				}}
			/>
		</Container>;
	},
));
