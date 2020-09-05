import React, { useState } from "react";
import { connect } from "react-redux";
import { useHistory } from "react-router";
import { Button, Container, Divider, Typography } from "@material-ui/core";
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";
import { FaGitlab, FaGithub } from "react-icons/fa";

import { useDependency } from "../../../container";
import { AuthService, UserRole } from "../../../services";
import { setAuthenticated, setRole } from "../../layout";
import { OAuthPopup } from "../components/OAuthPopup";

interface Props {
	setLoggedIn: () => void;
	setUserRole: (role?: UserRole) => void;
}

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		container: {
			height: "100%",
			display: "flex",
			alignItems: "center",
			justifyContent: "center",
		},
		login: {
			textAlign: "center",
			maxWidth: "180px",
			marginRight: "auto",
			marginLeft: "auto",
		},
		gitlab: {
			background: "linear-gradient(45deg, #e24228 30%, #fca326 90%)",
			display: "block",
			marginBottom: theme.spacing(1),
		},
		github: {
			background: "linear-gradient(45deg, #272727 30%, #464646 90%)",
			display: "block",
			color: "white",
		},
		divider: {
			margin: `${theme.spacing(1)}px 0px`,
		},
	})
);

export const Login = connect(null, (dispatch) => ({
	setUserRole: (role?: UserRole) => dispatch(setRole(role)),
	setLoggedIn: () => dispatch(setAuthenticated(true)),
}))(({ setLoggedIn, setUserRole }: Props) => {
	const classes = useStyles();
	const history = useHistory();
	const authService = useDependency(AuthService);
	const [isOpen, setOpen] = useState(false);
	const [provider, setProvider] = useState<"github" | "gitlab">("github");
	const [error, setError] = useState<string | undefined>(undefined);

	return (
		<Container className={classes.container}>
			<div>
				<div className={classes.login}>
					<Typography variant={"h5"}>Brunel CI</Typography>

					<Divider className={classes.divider} />

					<Button
						className={classes.gitlab}
						onClick={() => {
							setProvider("gitlab");
							setOpen(true);
						}}
					>
						<FaGitlab
							style={{
								paddingRight: "10px",
								verticalAlign: "middle",
							}}
						/>
						Login with GitLab
					</Button>

					<Button
						className={classes.github}
						onClick={() => {
							setProvider("github");
							setOpen(true);
						}}
					>
						<FaGithub
							style={{
								paddingRight: "10px",
								verticalAlign: "middle",
							}}
						/>
						Login with GitHub
					</Button>
				</div>
				{error && (
					<Typography
						style={{
							maxWidth: "250px",
							textAlign: "center",
							marginTop: "10px",
						}}
						color="error"
					>
						{error}
					</Typography>
				)}
			</div>
			<OAuthPopup
				isOpen={isOpen}
				provider={provider}
				onDone={(e) => {
					setError(undefined);
					setOpen(false);
					authService.setAuthentication(e);
					setUserRole(authService.getRole());
					setLoggedIn();
					history.push("/repository/");
				}}
				onError={(e) => {
					setError(e);
					setOpen(false);
				}}
				onAbort={() => {
					setError("Login aborted.");
					setOpen(false);
				}}
			/>
		</Container>
	);
});
