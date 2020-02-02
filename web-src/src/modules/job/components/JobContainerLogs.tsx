import React, {useEffect, useState} from 'react';
import {withDependency} from "../../../container";
import {JobService} from "../../../services";

interface Dependencies {
    jobService: JobService;
}

interface Props {
    containerId: string
}

require('./JobContainerLogs.css');

export const JobContainerLogs = withDependency<Props, Dependencies>((container) => ({
    jobService: container.get(JobService),
}))(({jobService, containerId}) => {

    const [content, setContent] = useState<null | HTMLDivElement>();

    useEffect(() => {
        if (content) {
            const subscription = jobService
                .containerLogs(containerId)
                .subscribe(
                    (progress) => {
                        content.innerHTML = progress;
                    }
                );

            return () => {
                return subscription.unsubscribe();
            };
        }
    }, [jobService, containerId, content]);

    return <div className={'term-container'} ref={(r) => setContent(r)} />;
});

