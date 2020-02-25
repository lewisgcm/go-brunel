import React from 'react';

import {Drawer} from '../../layout';
import {EnvironmentListComponent} from '../components/EnvironmentListComponent';
import {withDependency} from '../../../container';
import {EnvironmentService} from '../../../services';

export const EnvironmentPage = withDependency(
	(container) => ({
		environmentService: container.get(EnvironmentService),
	}),
)(() => {
	return <Drawer sidebar={() => <EnvironmentListComponent isLoading={false} environments={[]} onClick={() => {}} onSearch={() => {}} selectedEnvironmentId={''}/>}
		content={() => <React.Fragment>
			<h1>Hello World</h1>
		</React.Fragment>}/>;
});
