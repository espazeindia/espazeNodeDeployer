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

type NodeRepository interface {
	Create(ctx context.Context, node *entities.Node) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Node, error)
	GetByMacAddress(ctx context.Context, macAddress string) (*entities.Node, error)
	GetAll(ctx context.Context, filters map[string]interface{}) ([]*entities.Node, error)
	Update(ctx context.Context, id primitive.ObjectID, update *entities.NodeUpdateRequest) error
	UpdateResources(ctx context.Context, id primitive.ObjectID, resources *entities.NodeResources) error
	UpdateLastSeen(ctx context.Context, id primitive.ObjectID) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetNodesByStatus(ctx context.Context, status entities.NodeStatus) ([]*entities.Node, error)
	GetNodesByLocation(ctx context.Context, latitude, longitude, radiusKm float64) ([]*entities.Node, error)
	GetNodeStats(ctx context.Context) (map[string]interface{}, error)
}

type nodeRepository struct {
	collection *mongo.Collection
}

func NewNodeRepository(db *mongo.Database) NodeRepository {
	collection := db.Collection("nodes")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "mac_address", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "node_name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "location.latitude", Value: 1}, {Key: "location.longitude", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "last_seen_at", Value: -1}},
		},
	}
	
	collection.Indexes().CreateMany(ctx, indexes)
	
	return &nodeRepository{collection: collection}
}

func (r *nodeRepository) Create(ctx context.Context, node *entities.Node) error {
	node.ID = primitive.NewObjectID()
	node.CreatedAt = time.Now()
	node.UpdatedAt = time.Now()
	node.LastSeenAt = time.Now()
	
	if node.Status == "" {
		node.Status = entities.NodeStatusOnline
	}
	
	_, err := r.collection.InsertOne(ctx, node)
	return err
}

func (r *nodeRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Node, error) {
	var node entities.Node
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&node)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &node, nil
}

func (r *nodeRepository) GetByMacAddress(ctx context.Context, macAddress string) (*entities.Node, error) {
	var node entities.Node
	err := r.collection.FindOne(ctx, bson.M{"mac_address": macAddress}).Decode(&node)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &node, nil
}

func (r *nodeRepository) GetAll(ctx context.Context, filters map[string]interface{}) ([]*entities.Node, error) {
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
	
	var nodes []*entities.Node
	if err = cursor.All(ctx, &nodes); err != nil {
		return nil, err
	}
	
	return nodes, nil
}

func (r *nodeRepository) Update(ctx context.Context, id primitive.ObjectID, update *entities.NodeUpdateRequest) error {
	updateDoc := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}
	
	if update.Status != "" {
		updateDoc["$set"].(bson.M)["status"] = update.Status
	}
	
	if update.Location != nil {
		updateDoc["$set"].(bson.M)["location"] = update.Location
	}
	
	if update.Resources != nil {
		updateDoc["$set"].(bson.M)["resources"] = update.Resources
	}
	
	if update.ClusterInfo != nil {
		updateDoc["$set"].(bson.M)["cluster_info"] = update.ClusterInfo
	}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, updateDoc)
	return err
}

func (r *nodeRepository) UpdateResources(ctx context.Context, id primitive.ObjectID, resources *entities.NodeResources) error {
	update := bson.M{
		"$set": bson.M{
			"resources":   resources,
			"updated_at":  time.Now(),
			"last_seen_at": time.Now(),
		},
	}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *nodeRepository) UpdateLastSeen(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"last_seen_at": time.Now(),
		},
	}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *nodeRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *nodeRepository) GetNodesByStatus(ctx context.Context, status entities.NodeStatus) ([]*entities.Node, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var nodes []*entities.Node
	if err = cursor.All(ctx, &nodes); err != nil {
		return nil, err
	}
	
	return nodes, nil
}

func (r *nodeRepository) GetNodesByLocation(ctx context.Context, latitude, longitude, radiusKm float64) ([]*entities.Node, error) {
	// Simple bounding box calculation (not perfect but good enough)
	latDelta := radiusKm / 111.0 // 1 degree latitude ≈ 111 km
	lonDelta := radiusKm / (111.0 * cosApprox(latitude))
	
	filter := bson.M{
		"location.latitude": bson.M{
			"$gte": latitude - latDelta,
			"$lte": latitude + latDelta,
		},
		"location.longitude": bson.M{
			"$gte": longitude - lonDelta,
			"$lte": longitude + lonDelta,
		},
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var nodes []*entities.Node
	if err = cursor.All(ctx, &nodes); err != nil {
		return nil, err
	}
	
	return nodes, nil
}

func (r *nodeRepository) GetNodeStats(ctx context.Context) (map[string]interface{}, error) {
	pipeline := mongo.Pipeline{
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
	totalCount, _ := r.collection.CountDocuments(ctx, bson.M{})
	stats["total"] = totalCount
	
	return stats, nil
}

// Helper function for approximate cosine
func cosApprox(degrees float64) float64 {
	radians := degrees * 3.14159265359 / 180.0
	// Simple approximation
	if radians < -1.5708 || radians > 1.5708 {
		return 0.0
	}
	return 1.0 - (radians * radians / 2.0)
}

