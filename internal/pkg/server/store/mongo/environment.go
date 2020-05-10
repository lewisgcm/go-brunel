package mongo

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"time"
)

const (
	environmentCollectionName = "environment"
)

type EnvironmentStore struct {
	Database *mongo.Database
}

type mongoEnvironment struct {
	ObjectID          primitive.ObjectID `bson:"_id,omitempty"`
	store.Environment `bson:",inline"`
}

func (s *EnvironmentStore) Get(id shared.EnvironmentID) (*store.Environment, error) {
	objectID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}

	var d mongoEnvironment
	err = s.Database.Collection(environmentCollectionName).
		FindOne(
			context.Background(),
			bson.M{"_id": objectID},
		).
		Decode(&d)

	if err == mongo.ErrNoDocuments {
		return nil, store.ErrorNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "error getting environment")
	}

	d.Environment.ID = shared.EnvironmentID(d.ObjectID.Hex())
	return &d.Environment, nil
}

func (s *EnvironmentStore) nameIsUnique(id shared.EnvironmentID, name string) (bool, error) {
	var mongoEntity mongoEnvironment
	if e := s.
		Database.
		Collection(environmentCollectionName).
		FindOne(
			context.Background(),
			bson.M{"name": name},
		).Decode(&mongoEntity); e != nil {
		if e == mongo.ErrNoDocuments {
			return true, nil
		}
		return false, errors.Wrap(e, "error checking if name is unique")
	}

	return mongoEntity.ObjectID.Hex() == string(id), nil
}

func (s *EnvironmentStore) AddOrUpdate(environment store.Environment) (*store.Environment, error) {
	unique, e := s.nameIsUnique(environment.ID, environment.Name)
	if e != nil {
		return nil, e
	}
	if !unique {
		return nil, errors.New("environment name must be unique")
	}

	var mongoEntity mongoEnvironment
	mongoEntity.UpdatedAt = time.Now()
	if environment.ID == "" {
		mongoEntity.CreatedAt = time.Now()
	}

	upsert := true
	after := options.After

	if err := s.
		Database.
		Collection(environmentCollectionName).
		FindOneAndUpdate(
			context.Background(),
			bson.M{"name": environment.Name},
			bson.M{"$set": environment},
			&options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &after},
		).Decode(&mongoEntity); err != nil {
		return nil, errors.Wrap(err, "error adding or updating environment")
	}

	mongoEntity.Environment.ID = shared.EnvironmentID(mongoEntity.ObjectID.Hex())
	return &mongoEntity.Environment, nil
}

func (s *EnvironmentStore) Filter(filter string) ([]store.EnvironmentList, error) {
	environments := []store.EnvironmentList{}
	decoder, err := s.
		Database.
		Collection(environmentCollectionName).Aggregate(
		context.Background(),
		[]bson.M{
			{"$match": bson.M{"name": bsonx.Regex(fmt.Sprintf(".*%s.*", filter), "")}},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching environment list")
	}

	for decoder.Next(context.Background()) {
		var environment mongoEnvironment
		err = decoder.Decode(&environment)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding environment list")
		}
		environment.ID = shared.EnvironmentID(environment.ObjectID.Hex())
		environments = append(environments, store.EnvironmentList{
			ID:   environment.ID,
			Name: environment.Name,
		})
	}
	return environments, nil
}

func (s *EnvironmentStore) GetVariable(id shared.EnvironmentID, name string) (*string, error) {
	env, err := s.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "error getting environment")
	}

	for _, v := range env.Variables {
		if v.Name == name {
			return &v.Value, nil
		}
	}

	return nil, errors.New("environment variable not found")
}
