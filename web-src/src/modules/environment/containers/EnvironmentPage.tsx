import React from "react";
import {Drawer} from "../../layout";
import {EnvironmentListComponent} from "../components/EnvironmentListComponent";

export function EnvironmentPage() {
	return <Drawer sidebar={() => <EnvironmentListComponent isLoading={false} environments={[]} onClick={() => {}} onSearch={() => {}} selectedEnvironmentId={''}/>}
				   content={() => <React.Fragment>
					   <h1>Hello World</h1>
				   </React.Fragment>}/>;
}
