package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DEFAULT_CONNECTION_STRING_FORMAT string = "mongodb://%s:%s@%s"

type MongoDbConfiguration struct {
	ConnectionTimeout        time.Duration
	ConnectionStringTemplate string
}

type MongoDb struct {
	Config   MongoDbConfiguration
	Resource *mongo.Database
}

func CreateDatabase(config *MongoDbConfiguration) (MongoDb, error) {
	db := MongoDb{
		Config: *config,
	}
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	dbName := os.Getenv("MONGODB_DB_NAME")
	clusterEndpoint := os.Getenv("MONGODB_ENDPOINT")

	connectionURI := fmt.Sprintf(config.ConnectionStringTemplate, username, password, clusterEndpoint)
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		logrus.Errorf("Failed to create client: %v", err)
		return db, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logrus.Errorf("Failed to connect to server: %v", err)
		return db, err
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		logrus.Errorf("Failed to ping cluster: %v", err)
		return db, err
	}

	db.Resource = client.Database(dbName)
	return db, nil
}

func (db *MongoDb) DropAll() error {
	resource := db.Resource

	collections, _ := resource.ListCollectionNames(context.TODO(), bson.D{{}})
	for _, name := range collections {
		_, err := resource.Collection(name).DeleteMany(context.TODO(), bson.M{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *MongoDb) Close() {
	db.Resource.Client().Disconnect(context.TODO())
}
