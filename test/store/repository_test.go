package store

import (
	"go-brunel/internal/pkg/server/store"
	"go-brunel/test"
	"testing"
	"time"
)

func TestAddRepository(t *testing.T) {
	suite := setup(t)

	for _, repositoryStore := range suite.repositoryStores {
		now := time.Now()
		repo, err := repositoryStore.AddOrUpdate(store.Repository{
			Project: "project",
			Name:    "name",
			URI:     "http://uri.com",
			Triggers: []store.RepositoryTrigger{
				{
					Type:    store.RepositoryTriggerTypeBranch,
					Pattern: "pattern",
				},
			},
			DeletedAt: nil,
		})

		if err != nil {
			t.Fatalf("error saving repository: %s", err)
		}

		if e := repositoryStore.Delete(repo.ID, true); e != nil {
			t.Fatalf("error deleting repository: %s", e)
		}

		if repo.ID == "" {
			t.Errorf("repository id not returned")
		}

		test.ExpectString(t, "project", repo.Project)
		test.ExpectString(t, "name", repo.Name)
		test.ExpectString(t, "http://uri.com", repo.URI)

		if repo.DeletedAt != nil {
			t.Errorf("repository delated at should be nil")
		}

		if now.UnixNano() < repo.CreatedAt.UnixNano() {
			t.Errorf("repository created at date incorrect: %s before %s", repo.CreatedAt.String(), now.String())
		}

		if now.UnixNano() < repo.UpdatedAt.UnixNano() {
			t.Errorf("repository updated at date incorrect: %s before %s", repo.UpdatedAt.String(), now.String())
		}

		if len(repo.Triggers) != 1 || repo.Triggers[0].Pattern != "pattern" || repo.Triggers[0].Type != store.RepositoryTriggerTypeBranch {
			t.Errorf("repository triggers do not match")
		}
	}
}

func TestUpdateRepository(t *testing.T) {
	suite := setup(t)

	for _, repositoryStore := range suite.repositoryStores {

		repo, err := repositoryStore.AddOrUpdate(store.Repository{
			Project: "project",
			Name:    "name",
			URI:     "http://uri.com",
			Triggers: []store.RepositoryTrigger{
				{
					Type:    store.RepositoryTriggerTypeBranch,
					Pattern: "pattern",
				},
			},
			DeletedAt: nil,
		})

		if err != nil {
			t.Fatalf("error saving repository: %s", err)
		}

		now := time.Now()
		updatedRepo, err := repositoryStore.AddOrUpdate(store.Repository{
			Project: "project",
			Name:    "name",
			URI:     "http://uri2.com",
			Triggers: []store.RepositoryTrigger{
				{
					Type:    store.RepositoryTriggerTypeTag,
					Pattern: "patternz",
				},
			},
			CreatedAt: time.Now(),
		})

		if err != nil {
			t.Errorf("error saving updated repository: %s", err)
		}

		if e := repositoryStore.Delete(repo.ID, true); e != nil {
			t.Fatalf("error deleting repository: %s", e)
		}

		test.ExpectString(t, string(repo.ID), string(updatedRepo.ID))
		test.ExpectString(t, "name", updatedRepo.Name)
		test.ExpectString(t, "project", updatedRepo.Project)
		test.ExpectString(t, "http://uri2.com", updatedRepo.URI)

		if updatedRepo.DeletedAt != nil {
			t.Errorf("repository delated at should be nil")
		}

		if !repo.CreatedAt.Equal(updatedRepo.CreatedAt) {
			t.Errorf("repository created at date incorrect: %s != %s", repo.CreatedAt.String(), updatedRepo.CreatedAt.String())
		}

		if now.UnixNano() < updatedRepo.UpdatedAt.UnixNano() {
			t.Errorf("repository updated at date incorrect: %s before %s", updatedRepo.UpdatedAt.String(), now.String())
		}

		if len(updatedRepo.Triggers) != 1 || updatedRepo.Triggers[0].Pattern != "patternz" || updatedRepo.Triggers[0].Type != store.RepositoryTriggerTypeTag {
			t.Errorf("repository triggers do not match")
		}
	}
}

func TestGetRepository(t *testing.T) {
	suite := setup(t)

	for _, repositoryStore := range suite.repositoryStores {
		repo, err := repositoryStore.AddOrUpdate(store.Repository{
			Project: "project",
			Name:    "name",
			URI:     "http://uri.com",
			Triggers: []store.RepositoryTrigger{
				{
					Type:    store.RepositoryTriggerTypeBranch,
					Pattern: "pattern",
				},
			},
			DeletedAt: nil,
		})

		if err != nil {
			t.Fatalf("error saving repository: %s", err)
		}

		getRepo, err := repositoryStore.Get(repo.ID)

		if e := repositoryStore.Delete(repo.ID, true); e != nil {
			t.Fatalf("error deleting repository: %s", e)
		}

		if getRepo == nil || err != nil {
			t.Fatalf("error getting repository")
		}

		test.ExpectString(t, string(repo.ID), string(getRepo.ID))
		test.ExpectString(t, "name", getRepo.Name)
		test.ExpectString(t, "project", getRepo.Project)
		test.ExpectString(t, "http://uri.com", getRepo.URI)

		if repo.DeletedAt != getRepo.DeletedAt {
			t.Errorf("repository delated at should be nil")
		}

		if !repo.CreatedAt.Equal(getRepo.CreatedAt) {
			t.Errorf("repository created at date incorrect: %s != %s", repo.CreatedAt.String(), getRepo.CreatedAt.String())
		}

		if !repo.UpdatedAt.Equal(getRepo.UpdatedAt) {
			t.Errorf("repository updated at date incorrect: %s != %s", repo.UpdatedAt.String(), getRepo.UpdatedAt.String())
		}

		if len(getRepo.Triggers) != 1 || getRepo.Triggers[0].Pattern != "pattern" || getRepo.Triggers[0].Type != store.RepositoryTriggerTypeBranch {
			t.Errorf("repository triggers do not match")
		}
	}
}

func TestGetRepositoryNotFound(t *testing.T) {
	suite := setup(t)

	for _, repositoryStore := range suite.repositoryStores {
		_, err := repositoryStore.Get(store.RepositoryID("5eb1d158a610b1d1024f0d59"))
		test.ExpectError(t, store.ErrorNotFound, err)
	}
}

func TestFilterRepository(t *testing.T) {
	suite := setup(t)

	for _, repositoryStore := range suite.repositoryStores {
		repo, err := repositoryStore.AddOrUpdate(store.Repository{
			Project: "project",
			Name:    "name",
			URI:     "http://uri.com",
			Triggers: []store.RepositoryTrigger{
				{
					Type:    store.RepositoryTriggerTypeBranch,
					Pattern: "pattern",
				},
			},
			DeletedAt: nil,
		})

		if err != nil {
			t.Fatalf("error saving repository: %s", err)
		}

		if repos, err := repositoryStore.Filter(""); len(repos) == 0 || err != nil {
			t.Errorf("unexpected search result for '%s'", "")
		}

		if repos, err := repositoryStore.Filter("PROJECT"); len(repos) == 0 || err != nil {
			t.Errorf("unexpected search result for '%s'", "PROJECT")
		}

		if repos, err := repositoryStore.Filter("NAME"); len(repos) == 0 || err != nil {
			t.Errorf("unexpected search result for '%s'", "NAME")
		}

		if repos, err := repositoryStore.Filter("sdsdsdsdsd"); len(repos) != 0 || err != nil {
			t.Errorf("unexpected search result for '%s'", "sdsdsdsdsd")
		}

		if e := repositoryStore.Delete(repo.ID, true); e != nil {
			t.Fatalf("error deleting repository: %s", e)
		}
	}
}

func TestSetTriggerRepository(t *testing.T) {
	suite := setup(t)

	for _, repositoryStore := range suite.repositoryStores {
		repo, err := repositoryStore.AddOrUpdate(store.Repository{
			Project: "project",
			Name:    "name",
			URI:     "http://uri.com",
			Triggers: []store.RepositoryTrigger{
				{
					Type:    store.RepositoryTriggerTypeBranch,
					Pattern: "pattern",
				},
			},
			DeletedAt: nil,
		})

		if err != nil {
			t.Fatalf("error saving repository: %s", err)
		}

		if e := repositoryStore.SetTriggers(repo.ID, []store.RepositoryTrigger{
			{
				Type:    store.RepositoryTriggerTypeTag,
				Pattern: "patternz",
			},
		}); e != nil {
			t.Errorf("error setting triggers")
		}

		getRepo, err := repositoryStore.Get(repo.ID)
		if err != nil {
			t.Errorf("error getting repository: %s", err)
		}

		if len(getRepo.Triggers) != 1 || getRepo.Triggers[0].Pattern != "patternz" || getRepo.Triggers[0].Type != store.RepositoryTriggerTypeTag {
			t.Errorf("repository triggers do not match")
		}

		if e := repositoryStore.Delete(repo.ID, true); e != nil {
			t.Fatalf("error deleting repository: %s", e)
		}
	}
}
