
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
	StoppedBy: string;
	Duration: string | number;
	State: JobState;
	Commit: Commit;
	Repository: Repository;
}

export interface Log {
	Message: string;
	Type: number;
	Time: string;
	StageID: string;
}

export interface JobProgress {
	State: JobState;
	Stages: JobStage[];
}

export interface JobStage extends Stage {
	Containers: Container[];
	Logs: Log[];
}

export interface Stage {
	ID: string;
	JobID: string;
	StartedAt: string;
	State: number;
	StoppedAt: string;
}

export enum ContainerState {
	Starting = 0,
	Running = 1,
	Stopped = 2,
	Error = 3,
}

export interface Container {
	ID: string;
	JobID: string;
	ContainerID: string;
	State: ContainerState;
	Meta: {
		StageID: string;
		Service: boolean;
	};
	Spec: any;
	CreatedAt: string;
	StartedAt: string;
	StoppedAt: string;
}

export enum StageState {
	Running = 0,
	Success = 1,
	Error = 2
}

export interface User {
	Username: string;
	Email: string;
	Name: string;
	AvatarURL: string;
	Role: UserRole;
}

export interface UserList {
	Username: string;
	Role: UserRole;
}

export enum JobState {
	Waiting = 0,
	Processing = 1,
	Failed = 2,
	Success = 3,
	Cancelled = 4
}

export interface RepositoryJobPage {
	Count: number;
	Jobs: RepositoryJob[];
}

export interface RepositoryJob {
	ID: string;
	RepositoryID: string;
	Commit: {
		Branch: string;
		Revision: string;
	};
	State: JobState;
	StartedBy: string;
	CreatedAt: string;
	StartedAt: string;
	StoppedAt: string;
}

export interface Repository {
	ID: string;
	Project: string;
	Name: string;
	URI: string;
	Triggers: RepositoryTrigger[];
	CreatedAt: string;
}

export interface RepositoryTrigger {
	Type: RepositoryTriggerType;
	Pattern: string;
	EnvironmentID?: string;
}

export enum RepositoryTriggerType {
	Tag = 0,
	Branch = 1
}

export enum EnvironmentVariableType {
	Text = 0,
	Password = 1
}

export interface EnvironmentVariable {
	Type: EnvironmentVariableType;
	Name: string;
	Value: string;
}

export interface EnvironmentList {
	ID: string;
	Name: string;
}

export interface Environment extends EnvironmentList {
	Variables: EnvironmentVariable[];
}

export enum UserRole {
	Admin = 'admin',
	Reader = 'reader',
}
