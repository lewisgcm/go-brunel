// +build storeIntegrationTests !unit

package store

import (
	"context"
	"flag"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go-brunel/internal/pkg/server/store"
	mongo2 "go-brunel/internal/pkg/server/store/mongo"
	"testing"
)

type testSuite struct {
	environmentStores []store.EnvironmentStore
	repositoryStores  []store.RepositoryStore
	userStores        []store.UserStore
}

var mongoUri = ""

func init() {
	flag.StringVar(&mongoUri, "mongo-db-uri", "mongodb://root:example@localhost:27017", "Mongo Database URI")
}

func getMongo(t *testing.T) *mongo.Database {
	mongoClient, err := mongo.NewClient(mongoUri)
	if err != nil {
		t.Fatal(err)
	}

	err = mongoClient.Connect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return mongoClient.Database("brunel")
}

func setup(t *testing.T) testSuite {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	mongoDb := getMongo(t)

	// Initialize environment stores
	var environmentStores []store.EnvironmentStore
	environmentStores = append(environmentStores, &mongo2.EnvironmentStore{Database: mongoDb})

	// Initialize repository stores
	var repositoryStores []store.RepositoryStore
	repositoryStores = append(repositoryStores, &mongo2.RepositoryStore{Database: mongoDb})

	// Initialize user stores
	var userStores []store.UserStore
	userStores = append(userStores, &mongo2.UserStore{Database: mongoDb})

	return testSuite{
		environmentStores: environmentStores,
		repositoryStores:  repositoryStores,
		userStores:        userStores,
	}
}
