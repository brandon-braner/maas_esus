package usersapi

import (
	"context"
	"fmt"
	"log"

	"github.com/brandonbraner/maas/config"
	"github.com/brandonbraner/maas/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUserName(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, id string, user *User) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(ctx context.Context) (*userRepository, error) {
	client, err := database.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get database client: %w", err)
	}

	db := client.Database(config.AppConfig.MONGO_DB_NAME)
	collection := db.Collection(config.AppConfig.USER_COLLECTION_NAME)

	return &userRepository{
		collection: collection,
	}, nil
}

func (r *userRepository) Create(ctx context.Context, user *User) (*mongo.InsertOneResult, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return result, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByUserName(ctx context.Context, username string) (*User, error) {
	var user User
	// err := r.collection.FindOne(ctx, map[string]string{"_id": id}).Decode(&user)
	err := r.collection.FindOne(ctx, bson.M{"username": username}, options.FindOne().SetProjection(bson.M{"password": 0})).Decode(&user)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, id string, user *User) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": user},
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal("Failed to convert hex string to ObjectID:", err)
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
