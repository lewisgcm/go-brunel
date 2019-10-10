package store

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go-brunel/internal/pkg/server/store"
	mongo2 "go-brunel/internal/pkg/server/store/mongo"
	"go-brunel/internal/pkg/shared"
	"testing"
)

func TestSomething(t *testing.T) {
	mongoClient, err := mongo.NewClient("mongodb://root:example@localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	err = mongoClient.Connect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	db := mongoClient.Database("brunel")
	cs := mongo2.ContainerStore{
		Database: db,
	}

	err = cs.Add(store.Container{
		JobID:       shared.JobID("5cdc9c2a2bdfc4db82159bac"),
		ContainerID: "asdasd",
		State:       shared.ContainerStateStopped,
	})
	if err != nil {
		t.Fatal(err)
	}
}
