package mongo

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"time"
)

const (
	jobContainerCollectionName = "job_container"
)

type ContainerStore struct {
	Database *mongo.Database
}

type mongoContainer struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	JobID           primitive.ObjectID `bson:"job_id"`
	store.Container `bson:",inline"`
}

type mongoContainerUpdate struct {
	State     *shared.ContainerState `bson:"state,omitempty"`
	StoppedAt *time.Time             `bson:"stopped_at,omitempty"`
	StartedAt *time.Time             `bson:"started_at,omitempty"`
}

func (r *ContainerStore) Add(c store.Container) error {
	mContainer := mongoContainer{
		Container: c,
	}

	id, err := primitive.ObjectIDFromHex(string(c.JobID))
	if err != nil {
		return errors.Wrap(err, "invalid job id")
	}
	mContainer.JobID = id

	if _, err := r.
		Database.
		Collection(jobContainerCollectionName).
		InsertOne(
			context.Background(),
			mContainer,
		); err != nil {
		return errors.Wrap(err, "error adding container")
	}
	return nil
}

func (r *ContainerStore) update(id shared.ContainerID, update mongoContainerUpdate) error {
	_, err := r.
		Database.
		Collection(jobContainerCollectionName).
		UpdateOne(
			context.Background(),
			bson.M{"container_id": string(id)},
			bson.M{"$set": update},
		)
	return errors.Wrap(err, "error updating container")
}

func (r *ContainerStore) UpdateStateByContainerID(id shared.ContainerID, state shared.ContainerState) error {
	return r.update(id, mongoContainerUpdate{State: &state})
}

func (r *ContainerStore) UpdateStoppedAtByContainerID(id shared.ContainerID, t time.Time) error {
	return r.update(id, mongoContainerUpdate{StoppedAt: &t})
}

func (r *ContainerStore) UpdateStartedAtByContainerID(id shared.ContainerID, t time.Time) error {
	return r.update(id, mongoContainerUpdate{StartedAt: &t})
}

func (r *ContainerStore) FilterByJobID(i shared.JobID) ([]store.Container, error) {
	containers := []store.Container{}
	jobID, err := primitive.ObjectIDFromHex(string(i))
	if err != nil {
		return containers, errors.Wrap(err, "invalid job id")
	}

	decoder, err := r.
		Database.
		Collection(jobContainerCollectionName).
		Aggregate(
			context.Background(),
			[]bson.M{
				{"$match": bson.M{"job_id": jobID}},
				{"$sort": bson.M{"created_at": 1}},
			},
		)
	if err != nil {
		return containers, errors.Wrap(err, "error getting containers")
	}

	for decoder.Next(context.Background()) {
		var c mongoContainer
		err = decoder.Decode(&c)
		if err != nil {
			return containers, errors.Wrap(err, "error decoding container")
		}
		c.Container.ID = c.ID.Hex()
		c.Container.JobID = shared.JobID(c.JobID.Hex())
		containers = append(containers, c.Container)
	}
	return containers, nil
}
