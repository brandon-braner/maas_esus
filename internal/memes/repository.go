package memes

import (
	"context"
	"fmt"

	"github.com/brandonbraner/maas/config"
	"github.com/brandonbraner/maas/pkg/database"
	"go.mongodb.org/mongo-driver/mongo"
)

type memeRepository struct {
	collection *mongo.Collection
}

func NewMemeRepository(ctx context.Context) (*memeRepository, error) {
	client, err := database.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get database client: %w", err)
	}

	db := client.Database(config.AppConfig.MONGO_DB_NAME)
	collection := db.Collection(config.AppConfig.USER_COLLECTION_NAME)

	return &memeRepository{
		collection: collection,
	}, nil
}
