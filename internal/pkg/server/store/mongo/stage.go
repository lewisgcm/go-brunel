package mongo

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
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

func (r *StageStore) Get(jobID shared.JobID) ([]store.Stage, error) {
	stages := []store.Stage{}
	decoder, err := r.
		Database.
		Collection(stageCollectionName).Aggregate(
		context.Background(),
		[]bson.M{
			{"$match": bson.M{"job_id": jobID}},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching stage list")
	}

	for decoder.Next(context.Background()) {
		var stage store.Stage
		err = decoder.Decode(&stage)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding stage list")
		}
		stages = append(stages, stage)
	}
	return stages, nil
}
