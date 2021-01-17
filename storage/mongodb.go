package storage

import (
	"context"
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

func (db *MongoDb) InitializeResource() error {
	config := db.Config
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	dbName := os.Getenv("MONGODB_DB_NAME")
	clusterEndpoint := os.Getenv("MONGODB_ENDPOINT")

	connectionURI := fmt.Sprintf(config.ConnectionStringTemplate, username, password, clusterEndpoint)
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		logrus.Errorf("Failed to create client: %v", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logrus.Errorf("Failed to connect to server: %v", err)
		return err
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		logrus.Errorf("Failed to ping cluster: %v", err)
		return err
	}

	db.Resource = client.Database(dbName)
	return nil
}

func (db *MongoDb) Close() {
	logrus.Warning("Closing all db connections")
}
