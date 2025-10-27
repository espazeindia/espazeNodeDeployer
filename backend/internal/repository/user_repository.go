package repository

import (
	"context"
	"time"

	"github.com/espazeindia/espazeNodeDeployer/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, id primitive.ObjectID, update map[string]interface{}) error
	UpdateLastLogin(ctx context.Context, id primitive.ObjectID) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	collection := db.Collection("users")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	
	collection.Indexes().CreateMany(ctx, indexes)
	
	return &userRepository{collection: collection}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Status = entities.UserStatusActive
	
	if user.Role == "" {
		user.Role = entities.UserRoleUser
	}
	
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, id primitive.ObjectID, update map[string]interface{}) error {
	update["updated_at"] = time.Now()
	
	updateDoc := bson.M{"$set": update}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, updateDoc)
	return err
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"last_login": time.Now(),
		},
	}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

