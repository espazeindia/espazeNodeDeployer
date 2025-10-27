package repository

import (
	"context"
	"time"

	"github.com/espaze/espazeNodeDeployer/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GitHubTokenRepository interface {
	Create(ctx context.Context, token *entities.GitHubToken) error
	GetByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.GitHubToken, error)
	Update(ctx context.Context, userID primitive.ObjectID, token string) error
	Delete(ctx context.Context, userID primitive.ObjectID) error
}

type githubTokenRepository struct {
	collection *mongo.Collection
}

func NewGitHubTokenRepository(db *mongo.Database) GitHubTokenRepository {
	collection := db.Collection("github_tokens")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	
	collection.Indexes().CreateMany(ctx, indexes)
	
	return &githubTokenRepository{collection: collection}
}

func (r *githubTokenRepository) Create(ctx context.Context, token *entities.GitHubToken) error {
	token.ID = primitive.NewObjectID()
	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()
	
	_, err := r.collection.InsertOne(ctx, token)
	return err
}

func (r *githubTokenRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.GitHubToken, error) {
	var token entities.GitHubToken
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *githubTokenRepository) Update(ctx context.Context, userID primitive.ObjectID, token string) error {
	update := bson.M{
		"$set": bson.M{
			"token":      token,
			"updated_at": time.Now(),
		},
	}
	
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, bson.M{"user_id": userID}, update, opts)
	return err
}

func (r *githubTokenRepository) Delete(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID})
	return err
}

