package mongo

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"time"
)

const (
	jobLogCollectionName          = "job_log"
	jobContainerLogCollectionName = "job_container_log"
)

type mongoLog struct {
	JobID     primitive.ObjectID `bson:"job_id"`
	store.Log `bson:",inline"`
}

type LogStore struct {
	Database *mongo.Database
}

func (r *LogStore) Log(l store.Log) error {
	mLog := mongoLog{Log: l}
	objectId, err := primitive.ObjectIDFromHex(string(l.JobID))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid job id '%s'", l.JobID))
	}
	mLog.JobID = objectId

	_, err = r.
		Database.
		Collection(jobLogCollectionName).
		InsertOne(context.Background(), mLog)
	return errors.Wrap(err, "error logging job message")
}

func (r *LogStore) ContainerLog(l store.ContainerLog) error {
	_, err := r.
		Database.
		Collection(jobContainerLogCollectionName).
		InsertOne(context.Background(), l)
	return errors.Wrap(err, "error logging job message")
}

func (r *LogStore) FilterLogByJobIDFromTime(id shared.JobID, t time.Time) ([]store.Log, error) {
	logs := []store.Log{}
	objectId, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return logs, errors.Wrap(err, fmt.Sprintf("invalid job id '%s'", id))
	}

	decoder, err := r.
		Database.
		Collection(jobLogCollectionName).
		Aggregate(
			context.Background(),
			[]bson.M{
				{"$match": bson.M{"job_id": objectId, "time": bson.M{"$gte": primitive.DateTime(t.Unix())}}},
				{"$sort": bson.M{"time": 1}},
			},
		)
	if err != nil {
		return logs, errors.Wrap(err, "error getting logs")
	}

	for decoder.Next(context.Background()) {
		var c mongoLog
		err = decoder.Decode(&c)
		if err != nil {
			return logs, errors.Wrap(err, "error decoding log")
		}
		c.Log.JobID = shared.JobID(c.JobID.Hex())
		logs = append(logs, c.Log)
	}
	return logs, nil
}

func (r *LogStore) FilterContainerLogByContainerIDFromTime(id shared.ContainerID, t time.Time) ([]store.ContainerLog, error) {
	logs := []store.ContainerLog{}

	decoder, err := r.
		Database.
		Collection(jobContainerLogCollectionName).
		Aggregate(
			context.Background(),
			[]bson.M{
				{"$match": bson.M{"container_id": id, "time": bson.M{"$gte": primitive.DateTime(t.Unix())}}},
				{"$sort": bson.M{"time": 1}},
			},
		)
	if err != nil {
		return logs, errors.Wrap(err, "error getting container logs")
	}

	for decoder.Next(context.Background()) {
		var c store.ContainerLog
		err = decoder.Decode(&c)
		if err != nil {
			return logs, errors.Wrap(err, "error decoding container log")
		}
		logs = append(logs, c)
	}
	return logs, nil
}
