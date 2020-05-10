// +build storeIntegrationTests !unit

package store

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go-brunel/internal/pkg/server/store"
	mongo2 "go-brunel/internal/pkg/server/store/mongo"
	"testing"
)

type testSuite struct {
	environmentStores []store.EnvironmentStore
}

func getMongo(t *testing.T) *mongo.Database {
	mongoClient, err := mongo.NewClient("mongodb://root:example@localhost:27017")
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

	return testSuite{
		environmentStores: environmentStores,
	}
}
