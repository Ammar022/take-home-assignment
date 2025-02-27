package repo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB represents a MongoDB client connection
type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDBConnection creates a new MongoDB connection
func NewMongoDBConnection(uri, database string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create MongoDB client
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(100).        // Configure connection pool
		SetMinPoolSize(10).         // Minimum connections to maintain
		SetMaxConnIdleTime(30 * time.Second) // Idle connection timeout

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Verify connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		client:   client,
		database: client.Database(database),
	}, nil
}

// Collection returns a handle to the specified collection
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

// Close disconnects from MongoDB
func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}


