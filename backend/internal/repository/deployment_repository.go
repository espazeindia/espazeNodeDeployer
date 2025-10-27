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

type DeploymentRepository interface {
	Create(ctx context.Context, deployment *entities.Deployment) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Deployment, error)
	GetByNodeID(ctx context.Context, nodeID primitive.ObjectID) ([]*entities.Deployment, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*entities.Deployment, error)
	GetAll(ctx context.Context, filters map[string]interface{}) ([]*entities.Deployment, error)
	Update(ctx context.Context, id primitive.ObjectID, update map[string]interface{}) error
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status entities.DeploymentStatus) error
	UpdateMetrics(ctx context.Context, id primitive.ObjectID, metrics *entities.DeploymentMetrics) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetDeploymentsByStatus(ctx context.Context, status entities.DeploymentStatus) ([]*entities.Deployment, error)
	GetDeploymentStats(ctx context.Context, nodeID *primitive.ObjectID) (map[string]interface{}, error)
}

type deploymentRepository struct {
	collection *mongo.Collection
}

func NewDeploymentRepository(db *mongo.Database) DeploymentRepository {
	collection := db.Collection("deployments")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "node_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "context_path", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}
	
	collection.Indexes().CreateMany(ctx, indexes)
	
	return &deploymentRepository{collection: collection}
}

func (r *deploymentRepository) Create(ctx context.Context, deployment *entities.Deployment) error {
	deployment.ID = primitive.NewObjectID()
	deployment.CreatedAt = time.Now()
	deployment.UpdatedAt = time.Now()
	
	if deployment.Status == "" {
		deployment.Status = entities.DeploymentStatusPending
	}
	
	_, err := r.collection.InsertOne(ctx, deployment)
	return err
}

func (r *deploymentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Deployment, error) {
	var deployment entities.Deployment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&deployment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &deployment, nil
}

func (r *deploymentRepository) GetByNodeID(ctx context.Context, nodeID primitive.ObjectID) ([]*entities.Deployment, error) {
	filter := bson.M{"node_id": nodeID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var deployments []*entities.Deployment
	if err = cursor.All(ctx, &deployments); err != nil {
		return nil, err
	}
	
	return deployments, nil
}

func (r *deploymentRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*entities.Deployment, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var deployments []*entities.Deployment
	if err = cursor.All(ctx, &deployments); err != nil {
		return nil, err
	}
	
	return deployments, nil
}

func (r *deploymentRepository) GetAll(ctx context.Context, filters map[string]interface{}) ([]*entities.Deployment, error) {
	filter := bson.M{}
	for key, value := range filters {
		filter[key] = value
	}
	
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var deployments []*entities.Deployment
	if err = cursor.All(ctx, &deployments); err != nil {
		return nil, err
	}
	
	return deployments, nil
}

func (r *deploymentRepository) Update(ctx context.Context, id primitive.ObjectID, update map[string]interface{}) error {
	update["updated_at"] = time.Now()
	
	updateDoc := bson.M{"$set": update}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, updateDoc)
	return err
}

func (r *deploymentRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status entities.DeploymentStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	
	if status == entities.DeploymentStatusRunning {
		update["$set"].(bson.M)["deployed_at"] = time.Now()
	}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *deploymentRepository) UpdateMetrics(ctx context.Context, id primitive.ObjectID, metrics *entities.DeploymentMetrics) error {
	update := bson.M{
		"$set": bson.M{
			"metrics":                 metrics,
			"updated_at":              time.Now(),
			"last_health_check_at":    time.Now(),
		},
	}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *deploymentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *deploymentRepository) GetDeploymentsByStatus(ctx context.Context, status entities.DeploymentStatus) ([]*entities.Deployment, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var deployments []*entities.Deployment
	if err = cursor.All(ctx, &deployments); err != nil {
		return nil, err
	}
	
	return deployments, nil
}

func (r *deploymentRepository) GetDeploymentStats(ctx context.Context, nodeID *primitive.ObjectID) (map[string]interface{}, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.M{}}}
	if nodeID != nil {
		matchStage = bson.D{{Key: "$match", Value: bson.M{"node_id": nodeID}}}
	}
	
	pipeline := mongo.Pipeline{
		matchStage,
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$status"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	stats := make(map[string]interface{})
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	for _, result := range results {
		status := result["_id"].(string)
		count := result["count"]
		stats[status] = count
	}
	
	// Get total count
	filter := bson.M{}
	if nodeID != nil {
		filter["node_id"] = nodeID
	}
	totalCount, _ := r.collection.CountDocuments(ctx, filter)
	stats["total"] = totalCount
	
	return stats, nil
}

