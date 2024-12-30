package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/brandonbraner/maas/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once // use this so we only get one instance
)

// GetClient returns a singleton MongoDB client instance
func GetClient() (*mongo.Client, error) {

	uri := config.AppConfig.MONGODB_URI

	var err error
	clientOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientOptions := options.Client().ApplyURI(uri)
		clientInstance, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			return
		}

		// Ping the database to verify connection
		err = clientInstance.Ping(ctx, nil)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return clientInstance, nil
}
