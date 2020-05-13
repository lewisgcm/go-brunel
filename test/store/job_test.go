package store

import (
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"go-brunel/test"
	"testing"
	"time"
)

func addRepository(suite testSuite, t *testing.T) store.RepositoryID {
	r, e := suite.repositoryStores[0].AddOrUpdate(store.Repository{
		Project:   "project",
		Name:      "name",
		URI:       "uri",
		DeletedAt: nil,
	})
	if e != nil {
		t.Fatalf("error creating repository")
	}

	return r.ID
}

func removeRepository(suite testSuite, t *testing.T, id store.RepositoryID) {
	if e := suite.repositoryStores[0].Delete(id, true); e != nil {
		t.Fatalf("error deleting repository")
	}
}

func TestAddJob(t *testing.T) {
	suites := setup(t)
	repoId := addRepository(suites, t)
	defer removeRepository(suites, t, repoId)

	for _, jobStore := range suites.jobStores {
		now := time.Now()
		job, err := jobStore.Add(store.Job{
			RepositoryID:  repoId,
			EnvironmentID: nil,
			Commit: shared.Commit{
				Branch:   "branch",
				Revision: "revision",
			},
			State:     shared.JobStateSuccess,
			StartedBy: "startedBy",
		})
		if err != nil {
			t.Fatalf("could not create job: %e", err)
		}

		getJob, err := jobStore.Get(job.ID)
		if err != nil {
			t.Errorf("could not get job: %e", err)
		}

		if e := jobStore.Delete(job.ID); e != nil {
			t.Fatalf("error deleting job: %s", e)
		}

		test.ExpectString(t, string(repoId), string(getJob.RepositoryID))
		test.ExpectString(t, "startedBy", getJob.StartedBy)
		test.ExpectString(t, "branch", getJob.Commit.Branch)
		test.ExpectString(t, "revision", getJob.Commit.Revision)

		if getJob.State != shared.JobStateSuccess {
			t.Errorf("job state incorrect %d != %d", getJob.State, shared.JobStateSuccess)
		}

		if getJob.CreatedAt.Unix() < now.Unix() {
			t.Errorf("created date incorrect: %s < %s", getJob.CreatedAt.String(), now.String())
		}
	}
}

func TestUpdateJobStoppedAt(t *testing.T) {
	suites := setup(t)
	repoId := addRepository(suites, t)
	defer removeRepository(suites, t, repoId)

	for _, jobStore := range suites.jobStores {
		now := time.Now()
		job, err := jobStore.Add(store.Job{
			RepositoryID:  repoId,
			EnvironmentID: nil,
			Commit: shared.Commit{
				Branch:   "branch",
				Revision: "revision",
			},
			State:     shared.JobStateSuccess,
			StartedBy: "startedBy",
		})
		if err != nil {
			t.Fatalf("could not create job: %e", err)
		}

		if err := jobStore.UpdateStoppedAtByID(job.ID, now); err != nil {
			t.Errorf("error setting stopped at time: %s", err)
		}

		getJob, err := jobStore.Get(job.ID)
		if err != nil {
			t.Errorf("could not get job: %e", err)
		}

		if e := jobStore.Delete(job.ID); e != nil {
			t.Fatalf("error deleting job: %s", e)
		}

		if getJob.StoppedAt == nil || getJob.StoppedAt.Unix() < now.Unix() {
			t.Errorf("created date incorrect")
		}
	}
}

func TestUpdateJobCancelledBy(t *testing.T) {
	suites := setup(t)
	repoId := addRepository(suites, t)
	defer removeRepository(suites, t, repoId)

	for _, jobStore := range suites.jobStores {
		job, err := jobStore.Add(store.Job{
			RepositoryID:  repoId,
			EnvironmentID: nil,
			Commit: shared.Commit{
				Branch:   "branch",
				Revision: "revision",
			},
			State:     shared.JobStateSuccess,
			StartedBy: "startedBy",
		})
		if err != nil {
			t.Fatalf("could not create job: %e", err)
		}

		if err := jobStore.CancelByID(job.ID, "lewis"); err != nil {
			t.Errorf("error setting stopped at time: %s", err)
		}

		getJob, err := jobStore.Get(job.ID)
		if err != nil {
			t.Errorf("could not get job: %e", err)
		}

		if e := jobStore.Delete(job.ID); e != nil {
			t.Fatalf("error deleting job: %s", e)
		}

		test.ExpectString(t, "lewis", *getJob.StoppedBy)
	}
}

func TestUpdateJobUpdateState(t *testing.T) {
	suites := setup(t)
	repoId := addRepository(suites, t)
	defer removeRepository(suites, t, repoId)

	for _, jobStore := range suites.jobStores {
		job, err := jobStore.Add(store.Job{
			RepositoryID:  repoId,
			EnvironmentID: nil,
			Commit: shared.Commit{
				Branch:   "branch",
				Revision: "revision",
			},
			State:     shared.JobStateSuccess,
			StartedBy: "startedBy",
		})
		if err != nil {
			t.Fatalf("could not create job: %e", err)
		}

		if err := jobStore.UpdateStateByID(job.ID, shared.JobStateCancelled); err != nil {
			t.Errorf("error setting stopped at time: %s", err)
		}

		getJob, err := jobStore.Get(job.ID)
		if err != nil {
			t.Errorf("could not get job: %e", err)
		}

		if e := jobStore.Delete(job.ID); e != nil {
			t.Fatalf("error deleting job: %s", e)
		}

		if getJob.State != shared.JobStateCancelled {
			t.Errorf("incorrect job state")
		}
	}
}

func TestFilterJobsByRepositoryID(t *testing.T) {
	suites := setup(t)
	repoId := addRepository(suites, t)
	defer removeRepository(suites, t, repoId)

	for _, jobStore := range suites.jobStores {
		jobOne, err := jobStore.Add(store.Job{
			RepositoryID:  repoId,
			EnvironmentID: nil,
			Commit: shared.Commit{
				Branch:   "branch1",
				Revision: "revision1",
			},
			State:     shared.JobStateSuccess,
			StartedBy: "startedBy1",
		})
		if err != nil {
			t.Fatalf("could not create job: %e", err)
		}

		jobTwo, err := jobStore.Add(store.Job{
			RepositoryID:  repoId,
			EnvironmentID: nil,
			Commit: shared.Commit{
				Branch:   "branch2",
				Revision: "revision2",
			},
			State:     shared.JobStateSuccess,
			StartedBy: "startedBy2",
		})
		if err != nil {
			t.Errorf("could not create job: %e", err)
		}

		// Filter created at ascending
		page, err := jobStore.FilterByRepositoryID(repoId, "", 0, 1, "created_at", 1)
		if err != nil {
			t.Errorf("error filtering repository: %s", err)
		}
		if len(page.Jobs) != 1 || page.Jobs[0].Commit.Branch != "branch1" || page.Count != 2 {
			t.Errorf("unexpected page sorting by created ascending")
		}

		// Filter created at ascending skip first page
		page, err = jobStore.FilterByRepositoryID(repoId, "", 1, 1, "created_at", 1)
		if err != nil {
			t.Errorf("error filtering repository: %s", err)
		}
		if len(page.Jobs) != 1 || page.Jobs[0].Commit.Branch != "branch2" || page.Count != 2 {
			t.Errorf("unexpected page sorting by created ascending")
		}

		// Filter created at descending
		page, err = jobStore.FilterByRepositoryID(repoId, "", 0, 1, "created_at", -1)
		if err != nil {
			t.Errorf("error filtering repository: %s", err)
		}
		if len(page.Jobs) != 1 || page.Jobs[0].Commit.Branch != "branch2" || page.Count != 2 {
			t.Errorf("unexpected page sorting by created descending")
		}

		// Filter based on the branch name
		page, err = jobStore.FilterByRepositoryID(repoId, "branch2", 0, 10, "created_at", -1)
		if err != nil {
			t.Errorf("error filtering repository: %s", err)
		}
		if len(page.Jobs) != 1 || page.Jobs[0].Commit.Branch != "branch2" || page.Count != 1 {
			t.Errorf("unexpected page when filtering branch")
		}

		// Filter based on the revision
		page, err = jobStore.FilterByRepositoryID(repoId, "revision1", 0, 10, "created_at", -1)
		if err != nil {
			t.Errorf("error filtering repository: %s", err)
		}
		if len(page.Jobs) != 1 || page.Jobs[0].Commit.Revision != "revision1" || page.Count != 1 {
			t.Errorf("unexpected page when filtering revision")
		}

		// Filter based on the started by
		page, err = jobStore.FilterByRepositoryID(repoId, "startedBy1", 0, 10, "created_at", -1)
		if err != nil {
			t.Errorf("error filtering repository: %s", err)
		}
		if len(page.Jobs) != 1 || page.Jobs[0].StartedBy != "startedBy1" || page.Count != 1 {
			t.Errorf("unexpected page when filtering started by")
		}

		// Clean up
		if e := jobStore.Delete(jobOne.ID); e != nil {
			t.Errorf("could not delete job")
		}

		if e := jobStore.Delete(jobTwo.ID); e != nil {
			t.Errorf("could not delete job")
		}
	}
}
