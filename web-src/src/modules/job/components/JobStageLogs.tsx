import React from 'react';

import {JobStage} from '../../../services';

export function JobStageLogs({stage}: {stage: JobStage}) {
	if (stage.Logs.length > 0) {
		return <div key={stage.ID} className={'term-container'} >
			{stage.Logs.map((log, i) => <React.Fragment key={`${log.StageID}-${i}`}>
				{log.Message} <br/>
			</React.Fragment>)}
		</div>;
	}
	return <React.Fragment key={stage.ID} />;
}
