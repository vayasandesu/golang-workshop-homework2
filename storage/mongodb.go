package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbConfiguration struct {
	ConnectionTimeout        time.Duration
	ConnectionStringTemplate string
}

type MongoDb struct {
	Config   MongoDbConfiguration
	Resource *mongo.Database
}

func getContext(timeout time.Duration) context.Context {
	ctx, err := context.WithTimeout(context.Background(), timeout*time.Second)
	if err != nil {
		return nil
	}

	return ctx
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
	config := db.Config
	resource := db.Resource

	ctx, err := context.WithTimeout(context.Background(), config.ConnectionTimeout*time.Second)
	if err != nil {
		return errors.New("cannot drop database")
	}

	return resource.Drop(ctx)
}

func (db *MongoDb) Close() {
	config := db.Config
	ctx := getContext(config.ConnectionTimeout)
	if ctx != nil {
		db.Resource.Client().Disconnect(ctx)
	}
}
