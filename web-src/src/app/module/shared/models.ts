export interface User {
    Username: string;
    Email: string;
    Name: string;
    AvatarURL: string;
}

export interface JobPageable {
    Jobs: Job[];
    Count: number;
}

export interface Repository {
    ID: string;
    Project: string;
    Name: string;
    URI: string;
    Jobs: JobPageable;
}

export interface Commit {
    Branch: string;
    Revision: string;
}

export interface Job {
    ID: string;
    StartedBy: string;
    CreatedAt: string;
    StartedAt: string;
    StoppedAt: string;
    Duration: string | number;
    State: JobState;
    Commit: Commit;
}

export interface Log {
    Message: string;
    Type: number;
    Time: string;
    StageID: string;
}

export interface ContainerLog extends Log {
    ContainerID: string;
}

export interface JobProgress {
    Stages: JobStage[];
}

export interface JobContainer extends Container {
    Logs: ContainerLog[];
}

export interface JobStage extends Stage {
    Containers: JobContainer[];
    Logs: Log[];
}

export interface Stage {
    ID: string;
    JobID: string;
    StartedAt: string;
    State: number;
    StoppedAt: string;
}

export interface Container {
    ID: string;
    JobID: string;
    ContainerID: string;
    State: number;
    Meta: {
        StageID: string;
        Service: boolean;
    };
    Spec: any;
    CreatedAt: string;
    StartedAt: string;
    StoppedAt: string;
}

export enum JobState {
    Waiting = 0,
    Processing = 1,
    Failed = 2,
    Success = 3,
    Cancelled = 4,
}

export enum StageState {
    Running = 0,
    Success = 1,
    Error = 2
}

