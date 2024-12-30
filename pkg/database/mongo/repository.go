package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/brandonbraner/maas/pkg/database/mongo"

	"go.mongodb.org/mongo-driver/mongo"
)

const UserCollection = "users"

// Repository defines the interface for user repository operations
type Repository interface {
	Save(ctx context.Context, user *User) error
}

// MongoDBRepository implements the Repository interface using MongoDB
type MongoDBRepository struct {
	*mongo.MongoRepository[User]
	client     *mongo.Client
	database   string
	collection string
}

// NewMongoDBRepository creates a new MongoDB repository instance
func NewMongoDBRepository(client *mongo.Client, database string) *MongoDBRepository {
	collection := client.Database(database).Collection(UserCollection)
	return &MongoDBRepository{
		MongoRepository: mongo.NewMongoRepository[User](collection),
		client:          client,
		database:        database,
		collection:      UserCollection,
	}
}

// Save inserts a new user document into MongoDB
func (r *MongoDBRepository) Save(ctx context.Context, user *User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	// Check if username already exists
	filter := map[string]interface{}{"username": user.Username}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking username existence: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("username %s already exists", user.Username)
	}

	return r.MongoRepository.Save(ctx, user)
}
