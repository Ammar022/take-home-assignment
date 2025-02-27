package repo

import (
	"context"
	"log"
	"take-home-assignment/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LinkRepository handles database operations for links
type LinkRepository struct {
	db         *MongoDB
	collection *mongo.Collection
}

// NewLinkRepository creates a new link repository
func NewLinkRepository(db *MongoDB) *LinkRepository {
	collection := db.Collection("links")

	// Create indexes
	_, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "userId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "expiresAt", Value: 1}},
		},
	})

	if err != nil {
		log.Printf("Failed to create indexes: %v", err)
	}

	return &LinkRepository{
		db:         db,
		collection: collection,
	}
}

// Create adds a new link to the database
func (r *LinkRepository) Create(ctx context.Context, link models.Link) (models.Link, error) {
	if link.ID.IsZero() {
		link.ID = primitive.NewObjectID()
	}

	if link.CreatedAt.IsZero() {
		link.CreatedAt = time.Now()
	}

	_, err := r.collection.InsertOne(ctx, link)
	if err != nil {
		return models.Link{}, err
	}

	return link, nil
}

// GetByID retrieves a link by its ID
func (r *LinkRepository) GetByID(ctx context.Context, id primitive.ObjectID) (models.Link, error) {
	var link models.Link

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&link)
	if err != nil {
		return models.Link{}, err
	}

	return link, nil
}

// GetAll retrieves all links for a user
func (r *LinkRepository) GetAll(ctx context.Context, userID string, limit, offset int64) ([]models.Link, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"userId": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var links []models.Link
	if err := cursor.All(ctx, &links); err != nil {
		return nil, err
	}

	return links, nil
}

// Update updates an existing link
func (r *LinkRepository) Update(ctx context.Context, id primitive.ObjectID, link models.LinkUpdateDTO) error {
	update := bson.M{
		"$set": bson.M{},
	}

	if link.Title != "" {
		update["$set"].(bson.M)["title"] = link.Title
	}

	if link.URL != "" {
		update["$set"].(bson.M)["url"] = link.URL
	}

	if !link.ExpiresAt.IsZero() {
		update["$set"].(bson.M)["expiresAt"] = link.ExpiresAt
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete removes a link from the database
func (r *LinkRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteExpired removes all expired links
func (r *LinkRepository) DeleteExpired(ctx context.Context) (int64, error) {
	now := time.Now()
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"expiresAt": bson.M{"$lt": now},
	})

	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

// IncrementClicks increments the click count for a link
func (r *LinkRepository) IncrementClicks(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$inc": bson.M{"clicks": 1}},
	)
	return err
}
