package repo

import (
	"context"
	"log"
	"take-home-assignment/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// VisitRepository handles database operations for visits
type VisitRepository struct {
	db         *MongoDB
	collection *mongo.Collection
}

// NewVisitRepository creates a new visit repository
func NewVisitRepository(db *MongoDB) *VisitRepository {
	collection := db.Collection("visits")
	
	// Create indexes
	_, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "linkId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: 1}},
		},
	})
	
	if err != nil {
		log.Printf("Failed to create indexes: %v", err)
	}
	
	return &VisitRepository{
		db:         db,
		collection: collection,
	}
}

// Create records a new visit
func (r *VisitRepository) Create(ctx context.Context, visit models.Visit) error {
	_, err := r.collection.InsertOne(ctx, visit)
	return err
}

// GetVisitsByLinkID retrieves all visits for a specific link
func (r *VisitRepository) GetVisitsByLinkID(ctx context.Context, linkID string, limit, offset int64) ([]models.Visit, error) {
	objID, err := primitive.ObjectIDFromHex(linkID)
	if err != nil {
		return nil, err
	}
	
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "timestamp", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{"linkId": objID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var visits []models.Visit
	if err := cursor.All(ctx, &visits); err != nil {
		return nil, err
	}
	
	return visits, nil
}