package mongo

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/store"
)

const (
	stageCollectionName = "job_stage"
)

type StageStore struct {
	Database *mongo.Database
}

func (r *StageStore) AddOrUpdate(stage store.Stage) error {
	upsert := true
	after := options.After
	err := r.
		Database.
		Collection(stageCollectionName).
		FindOneAndUpdate(
			context.Background(),
			bson.M{"id": stage.ID, "job_id": stage.JobID},
			bson.M{"$set": stage},
			&options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &after},
		).Err()

	return errors.Wrap(err, "error adding or updating stage")
}
