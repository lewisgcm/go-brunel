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
	ObjectID  primitive.ObjectID          `bson:"_id,omitempty"`
	Name      string                      `bson:"name"`
	Variables []store.EnvironmentVariable `bson:"variables"`
	CreatedAt *time.Time                  `bson:"created_at,omitempty"`
	UpdatedAt time.Time                   `bson:"updated_at" json:",omitempty"`
	DeletedAt *time.Time                  `bson:"deleted_at" json:",omitempty"`
}

func (e *mongoEnvironment) ToEnvironment() *store.Environment {
	return &store.Environment{
		ID:        shared.EnvironmentID(e.ObjectID.Hex()),
		Name:      e.Name,
		Variables: e.Variables,
		CreatedAt: *e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
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

	return d.ToEnvironment(), nil
}

func (s *EnvironmentStore) nameIsUnique(id shared.EnvironmentID, name string) error {
	var mongoEntity mongoEnvironment
	if e := s.
		Database.
		Collection(environmentCollectionName).
		FindOne(
			context.Background(),
			bson.M{"name": name},
		).Decode(&mongoEntity); e != nil {
		if e == mongo.ErrNoDocuments {
			return nil
		}
		return errors.Wrap(e, "error checking if name is unique")
	}

	if mongoEntity.ObjectID.Hex() != string(id) {
		return errors.New("environment name must be unique")
	}

	return nil
}

func (s *EnvironmentStore) AddOrUpdate(environment store.Environment) (*store.Environment, error) {
	if e := s.nameIsUnique(environment.ID, environment.Name); e != nil {
		return nil, e
	}

	mongoEntity := mongoEnvironment{
		Name:      environment.Name,
		Variables: environment.Variables,
		UpdatedAt: time.Now(),
	}

	criteria := bson.M{"name": environment.Name}
	if environment.ID != "" {
		oid, err := primitive.ObjectIDFromHex(string(environment.ID))
		if err != nil {
			return nil, errors.Wrap(err, "error parsing object id")
		}
		criteria = bson.M{"_id": oid}
	}

	upsert := true
	after := options.After
	if err := s.
		Database.
		Collection(environmentCollectionName).
		FindOneAndUpdate(
			context.Background(),
			criteria,
			bson.M{
				"$set":         mongoEntity,
				"$setOnInsert": bson.M{"created_at": time.Now()},
			},
			&options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &after},
		).Decode(&mongoEntity); err != nil {
		return nil, errors.Wrap(err, "error adding or updating environment")
	}

	return mongoEntity.ToEnvironment(), nil
}

func (s *EnvironmentStore) Filter(filter string) ([]store.EnvironmentList, error) {
	environments := []store.EnvironmentList{}
	decoder, err := s.
		Database.
		Collection(environmentCollectionName).Aggregate(
		context.Background(),
		[]bson.M{
			{"$match": bson.M{"name": bsonx.Regex(fmt.Sprintf(".*%s.*", filter), "i")}},
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
		environments = append(environments, store.EnvironmentList{
			ID:   shared.EnvironmentID(environment.ObjectID.Hex()),
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

func (s *EnvironmentStore) Delete(id shared.EnvironmentID, hard bool) error {
	objectID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return err
	}

	if hard {
		_, err = s.
			Database.
			Collection(environmentCollectionName).
			DeleteOne(context.Background(), bson.M{"_id": objectID})

		return errors.Wrap(err, "error deleting")
	} else {
		err = s.
			Database.
			Collection(environmentCollectionName).
			FindOneAndUpdate(context.Background(), bson.M{"_id": objectID}, bson.M{"deleted_at": time.Now()}).
			Err()

		return errors.Wrap(err, "error soft deleting")
	}
}
