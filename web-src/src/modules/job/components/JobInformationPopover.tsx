import React, { useState } from "react";
import {
	Tooltip,
	Typography,
	Hidden,
	makeStyles,
	Theme,
	createStyles,
	Popover,
	IconButton,
} from "@material-ui/core";
import { IconType } from "react-icons/lib";

const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		typography: {
			padding: theme.spacing(2),
		},
		titleJobInfo: {
			paddingLeft: theme.spacing(2),
			fontSize: theme.typography.body2.fontSize,
			"& svg": {
				verticalAlign: "middle",
				height: "1.3em",
				width: "1.3em",
				marginRight: "8px",
			},
		},
	})
);

interface Props {
	icon: IconType;
	information: string;
	tooltipText: string;
	popover: React.ReactNode;
}

export function JobInformationPopover({
	icon,
	information,
	tooltipText,
	popover,
}: Props) {
	const classes = useStyles();
	const Icon = icon;
	const [anchorEl, setAnchorEl] = React.useState<HTMLButtonElement | null>(
		null
	);
	const [isOpen, setIsOpen] = useState(false);

	return (
		<React.Fragment>
			<Hidden mdDown>
				<Tooltip title={tooltipText}>
					<Typography className={classes.titleJobInfo}>
						{<Icon />}
						{information}
					</Typography>
				</Tooltip>
			</Hidden>
			<Hidden lgUp>
				<IconButton
					color="inherit"
					ref={(e) => setAnchorEl(e)}
					onClick={() => setIsOpen(true)}
				>
					{<Icon />}
				</IconButton>
				<Popover
					open={isOpen}
					anchorEl={anchorEl}
					onClose={() => setIsOpen(false)}
					anchorOrigin={{
						vertical: "bottom",
						horizontal: "center",
					}}
					transformOrigin={{
						vertical: "top",
						horizontal: "center",
					}}
				>
					<Typography className={classes.typography}>
						{popover}
					</Typography>
				</Popover>
			</Hidden>
		</React.Fragment>
	);
}
