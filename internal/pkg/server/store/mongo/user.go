package mongo

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/security"
	"go-brunel/internal/pkg/server/store"
	"time"
)

const (
	userCollectionName = "user"
)

type UserStore struct {
	Database *mongo.Database
}

type mongoUser struct {
	Username  string            `bson:"username"`
	Email     string            `bson:"email"`
	Name      string            `bson:"name"`
	AvatarURL string            `bson:"avatar_url"`
	Role      security.UserRole `bson:"role"`
	CreatedAt *time.Time        `bson:"created_at,omitempty"`
}

func (u *mongoUser) ToUser() *store.User {
	return &store.User{
		Username:  u.Username,
		Email:     u.Email,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Role:      u.Role,
		CreatedAt: *u.CreatedAt,
	}
}

func (r *UserStore) Filter(filter string) ([]store.UserList, error) {
	users := []store.UserList{}
	decoder, err := r.
		Database.
		Collection(userCollectionName).Aggregate(
		context.Background(),
		[]bson.M{
			{"$match": bson.M{
				"$or": []bson.M{
					{"username": bsonx.Regex(fmt.Sprintf(".*%s.*", filter), "i")},
					{"email": bsonx.Regex(fmt.Sprintf(".*%s.*", filter), "i")},
				}}},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching user list")
	}

	for decoder.Next(context.Background()) {
		var user mongoUser
		err = decoder.Decode(&user)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding user list")
		}
		users = append(users, store.UserList{
			Username: user.Username,
			Role:     user.Role,
		})
	}

	return users, nil
}

func (r *UserStore) AddOrUpdate(user store.User) (*store.User, error) {
	entity := mongoUser{
		Username:  user.Username,
		Email:     user.Email,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
	}

	upsert := true
	after := options.After
	if err := r.
		Database.
		Collection(userCollectionName).
		FindOneAndUpdate(
			context.Background(),
			bson.M{"username": user.Username},
			bson.M{
				"$set":         entity,
				"$setOnInsert": bson.M{"created_at": time.Now()},
			},
			&options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &after},
		).Decode(&entity); err != nil {
		return nil, errors.Wrap(err, "error saving user")
	}

	return entity.ToUser(), nil
}

func (r *UserStore) GetByUsername(username string) (*store.User, error) {
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
		return nil, store.ErrorNotFound
	}

	if err != nil {
		return nil, errors.Wrap(err, "error getting user")
	}

	return f.ToUser(), nil
}

func (r *UserStore) Delete(username string, hard bool) error {
	if hard {
		_, err := r.
			Database.
			Collection(userCollectionName).
			DeleteOne(context.Background(), bson.M{"username": username})

		return errors.Wrap(err, "error deleting")
	} else {
		err := r.
			Database.
			Collection(userCollectionName).
			FindOneAndUpdate(context.Background(), bson.M{"username": username}, bson.M{"deleted_at": time.Now()}).
			Err()

		return errors.Wrap(err, "error soft deleting")
	}
}
