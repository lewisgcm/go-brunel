import React from "react";
import { JobStage } from "../../../services";
import {createStyles, makeStyles, Theme} from "@material-ui/core/styles";

interface StageGraphProps {
    stages: JobStage[];
    selectedStageId: string;
    onStageSelect: (stage: JobStage) => void;
}

const stageSpacing = 100;

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        svg: {
            marginBottom: '10px',
        },
        line: {
            stroke: 'lightslategrey',
            strokeWidth: 5,
        },
        stageText: {
            fill: 'lightslategrey',
        },
        default: {
            fill: 'lightslategrey',
        },
        stage: {
            fill: 'lightgreen',
            '&:hover': {
                strokeWidth: '2px',
                stroke: 'darkgreen',
                cursor: 'pointer',
            },
            '&:global(.selected)': {
                strokeWidth: '5px',
                stroke: 'darkgreen',
            },
        },
    }),
);

export const JobProgressGraph = (props: StageGraphProps) => {
    const classes = useStyles();
    return <svg width="100%" height="100" viewBox={`0 0 300 100`} className={classes.svg} >

        {/* Render the starting point in our graph. */}
        <g transform={`translate(${-stageSpacing}, 50) rotate(0)`} >
            <line x1={0} y1="0" x2={stageSpacing} y2="0" className={classes.line} />
            <text x="0" y="35" textAnchor="middle" className={classes.stageText}>start</text>
            <circle cx="0" cy="0" r="10" className={classes.default} />
        </g>
        {
            props.stages.map(
                (stage, index) => <g key={stage.ID}
                                     transform={`translate(${index * stageSpacing}, 50) rotate(0)`}
                                     onClick={() => props.onStageSelect(stage)} >
                    <line x1={0} y1="0" x2={stageSpacing} y2="0" className={classes.line} />
                    <text x="0" y="35" textAnchor="middle" className={classes.stageText}>{stage.ID}</text>
                    <circle cx="0" cy="0" r="20" className={`${classes.stage} ${props.selectedStageId === stage.ID ? "selected" : ""}`} />
                </g>,
            )
        }
        {/* Render the ending point in our graph. */}
        <g transform={`translate(${props.stages.length * stageSpacing},50) rotate(0)`} >
            <text x="0" y="35" textAnchor="middle" className={classes.stageText}>end</text>
            <circle cx="0" cy="0" r="10" className={classes.default} />
        </g>
    </svg>;
};
