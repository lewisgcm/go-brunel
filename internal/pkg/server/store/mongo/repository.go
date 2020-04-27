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
	"time"
)

const (
	repositoryCollectionName = "repository"
)

type RepositoryStore struct {
	Database *mongo.Database
}

type mongoRepository struct {
	ObjectID         primitive.ObjectID `bson:"_id,omitempty"`
	store.Repository `bson:",inline"`
}

func (r *RepositoryStore) AddOrUpdate(repository store.Repository) (store.Repository, error) {
	var repo mongoRepository
	repo.UpdatedAt = time.Now()
	upsert := true
	after := options.After
	err := r.
		Database.
		Collection(repositoryCollectionName).
		FindOneAndUpdate(
			context.Background(),
			bson.M{"project": repository.Project, "name": repository.Name},
			bson.M{"$set": repository},
			&options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &after},
		).Decode(&repo)
	if err != nil {
		return store.Repository{}, errors.Wrap(err, "error adding or updating repository")
	}

	repo.Repository.ID = store.RepositoryID(repo.ObjectID.Hex())
	return repo.Repository, nil
}

func (r *RepositoryStore) SetTriggers(id store.RepositoryID, triggers []store.RepositoryTrigger) error {
	objectID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return errors.Wrap(err, "error parsing object id")
	}

	if e := r.
		Database.
		Collection(repositoryCollectionName).
		FindOneAndUpdate(
			context.Background(),
			bson.M{"_id": objectID},
			bson.M{"$set": bson.M{
				"triggers":   triggers,
				"updated_at": time.Now(),
			}},
			&options.FindOneAndUpdateOptions{},
		).Err();
		e != nil {
		return errors.Wrap(err, "error adding or updating repository")
	}

	return nil
}

func (r *RepositoryStore) Get(id store.RepositoryID) (store.Repository, error) {
	objectID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return store.Repository{}, err
	}

	var d mongoRepository
	err = r.Database.Collection(repositoryCollectionName).
		FindOne(
			context.Background(),
			bson.M{"_id": objectID},
		).
		Decode(&d)

	if err == mongo.ErrNoDocuments {
		return store.Repository{}, store.ErrorNotFound
	}
	if err != nil {
		return store.Repository{}, errors.Wrap(err, "error getting repository")
	}

	d.Repository.ID = store.RepositoryID(d.ObjectID.Hex())
	return d.Repository, nil
}

func (r *RepositoryStore) Filter(filter string) ([]store.Repository, error) {
	repos := []store.Repository{}
	decoder, err := r.
		Database.
		Collection(repositoryCollectionName).Aggregate(
		context.Background(),
		[]bson.M{
			{"$match": bson.M{
				"$or": []bson.M{
					{"name": bsonx.Regex(fmt.Sprintf(".*%s.*", filter), "")},
					{"project": bsonx.Regex(fmt.Sprintf(".*%s.*", filter), "")},
				}}},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching repository list")
	}

	for decoder.Next(context.Background()) {
		var repo mongoRepository
		err = decoder.Decode(&repo)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding repository list")
		}
		repo.ID = store.RepositoryID(repo.ObjectID.Hex())
		repos = append(repos, repo.Repository)
	}
	return repos, nil
}
