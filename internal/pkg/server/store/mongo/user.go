package mongo

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/store"
)

const (
	userCollectionName = "user"
)

type UserStore struct {
	Database *mongo.Database
}

type mongoUser struct {
	ObjectID   primitive.ObjectID `bson:"_id,omitempty"`
	store.User `bson:",inline"`
}

func (r *UserStore) AddOrUpdate(user store.User) (store.User, error) {
	var f mongoUser
	upsert := true
	after := options.After
	err := r.
		Database.
		Collection(userCollectionName).
		FindOneAndUpdate(
			context.Background(),
			bson.M{"username": user.Username},
			bson.M{"$set": bson.M{"username": user.Username, "avatar_url": user.AvatarURL, "name": user.Name, "email": user.Email}},
			&options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &after},
		).Decode(&f)
	if err != nil {
		return f.User, errors.Wrap(err, "error upserting user")
	}
	f.User.ID = f.ObjectID.Hex()
	return f.User, nil
}

func (r *UserStore) GetByUsername(username string) (store.User, error) {
	var f mongoUser
	err := r.
		Database.
		Collection(userCollectionName).
		FindOne(
			context.Background(),
			bson.M{"username": username},
		).
		Decode(&f)
	if err == mongo.ErrNoDocuments {
		return f.User, store.ErrorNotFound
	}
	if err != nil {
		return f.User, errors.Wrap(err, "error getting user")
	}
	f.User.ID = f.ObjectID.Hex()
	return f.User, nil
}
