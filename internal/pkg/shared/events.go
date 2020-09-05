package shared

type RepositoryCreated struct {
	RepositoryID RepositoryID
}

type JobCreated struct {
	JobID        JobID
	RepositoryID RepositoryID
}

type JobUpdated struct {
	JobID JobID
}

type EnvironmentCreated struct {
	EnvironmentID EnvironmentID
}

type EnvironmentUpdated struct {
	EnvironmentID EnvironmentID
}

type event struct {
	Type    string
	Payload interface{}
}

func NewRepositoryCreated(repositoryID RepositoryID) interface{} {
	return event{
		Type: "REPOSITORY_CREATED",
		Payload: &RepositoryCreated{
			RepositoryID: repositoryID,
		},
	}
}

func NewJobCreated(jobID JobID, repositoryID RepositoryID) interface{} {
	return event{
		Type: "JOB_CREATED",
		Payload: &JobCreated{
			JobID:        jobID,
			RepositoryID: repositoryID,
		},
	}
}

func NewJobUpdated(jobID JobID) interface{} {
	return event{
		Type: "JOB_UPDATED",
		Payload: &JobUpdated{
			JobID: jobID,
		},
	}
}

func NewEnvironmentCreated(id EnvironmentID) interface{} {
	return event{
		Type:    "ENVIRONMENT_CREATED",
		Payload: &EnvironmentCreated{EnvironmentID: id},
	}
}

func NewEnvironmentUpdated(id EnvironmentID) interface{} {
	return event{
		Type:    "ENVIRONMENT_UPDATED",
		Payload: &EnvironmentUpdated{EnvironmentID: id},
	}
}
