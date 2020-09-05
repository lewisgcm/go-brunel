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
	repositoryCollectionName = "repository"
)

type RepositoryStore struct {
	Database *mongo.Database
}

type mongoRepository struct {
	ObjectID  primitive.ObjectID        `bson:"_id,omitempty"`
	Project   string                    `bson:"project"`
	Name      string                    `bson:"name"`
	URI       string                    `bson:"uri"`
	Triggers  []store.RepositoryTrigger `bson:",omitempty"`
	CreatedAt *time.Time                `bson:"created_at,omitempty"`
	UpdatedAt time.Time                 `bson:"updated_at"`
	DeletedAt *time.Time                `bson:"deleted_at" json:",omitempty"`
}

func (r *mongoRepository) ToRepository() *store.Repository {
	return &store.Repository{
		ID:        shared.RepositoryID(r.ObjectID.Hex()),
		Project:   r.Project,
		Name:      r.Name,
		URI:       r.URI,
		Triggers:  r.Triggers,
		CreatedAt: *r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}
}

func (r *RepositoryStore) AddOrUpdate(repository store.Repository) (*store.Repository, error) {
	repo := mongoRepository{
		Project:   repository.Project,
		Name:      repository.Name,
		URI:       repository.URI,
		Triggers:  repository.Triggers,
		UpdatedAt: time.Now(),
		DeletedAt: repository.DeletedAt,
	}

	upsert := true
	after := options.After
	if err := r.
		Database.
		Collection(repositoryCollectionName).
		FindOneAndUpdate(
			context.Background(),
			bson.M{"project": repository.Project, "name": repository.Name},
			bson.M{"$set": repo, "$setOnInsert": bson.M{"created_at": time.Now()}},
			&options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &after},
		).Decode(&repo); err != nil {
		return nil, errors.Wrap(err, "error saving repository")
	}

	return repo.ToRepository(), nil
}

func (r *RepositoryStore) SetTriggers(id shared.RepositoryID, triggers []store.RepositoryTrigger) error {
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
		).Err(); e != nil {
		return errors.Wrap(err, "error adding or updating repository")
	}

	return nil
}

func (r *RepositoryStore) Get(id shared.RepositoryID) (*store.Repository, error) {
	objectID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}

	var d mongoRepository
	err = r.Database.Collection(repositoryCollectionName).
		FindOne(
			context.Background(),
			bson.M{"_id": objectID},
		).
		Decode(&d)

	if err == mongo.ErrNoDocuments {
		return nil, store.ErrorNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "error getting repository")
	}

	return d.ToRepository(), nil
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
					{"name": bsonx.Regex(fmt.Sprintf(".*%s.*", filter), "i")},
					{"project": bsonx.Regex(fmt.Sprintf(".*%s.*", filter), "i")},
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
		repos = append(repos, *repo.ToRepository())
	}
	return repos, nil
}

func (r *RepositoryStore) Delete(id shared.RepositoryID, hard bool) error {
	objectID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return err
	}

	if hard {
		_, err = r.
			Database.
			Collection(repositoryCollectionName).
			DeleteOne(context.Background(), bson.M{"_id": objectID})

		return errors.Wrap(err, "error deleting")
	} else {
		err = r.
			Database.
			Collection(repositoryCollectionName).
			FindOneAndUpdate(context.Background(), bson.M{"_id": objectID}, bson.M{"deleted_at": time.Now()}).
			Err()

		return errors.Wrap(err, "error soft deleting")
	}
}
