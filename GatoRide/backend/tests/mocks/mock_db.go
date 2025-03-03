package mocks

import (
	"backend/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CollectionInterface defines the interface that real and mock collections must implement.
type CollectionInterface interface {
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
}

// MockCollection simulates a MongoDB collection and implements CollectionInterface.
type MockCollection struct{}

// Find simulates a MongoDB Find operation.
func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	// Simulated ride data
	rides := []models.Ride{
		{
			Pickup: models.Location{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Address:   "New York, NY",
			},
			Dropoff: models.Location{
				Latitude:  40.7306,
				Longitude: -73.9352,
				Address:   "Brooklyn, NY",
			},
			Status: "open",
			Price:  25.5,
		},
	}

	// Convert rides into BSON format
	var docs []interface{}
	for _, ride := range rides {
		docs = append(docs, ride)
	}

	// Create a mock cursor
	cursor, err := mongo.NewCursorFromDocuments(docs, nil, bson.DefaultRegistry)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

// MockGetCollection returns a mock collection that implements CollectionInterface.
func MockGetCollection(name string) CollectionInterface {
	if name == "rides" {
		return &MockCollection{} // ✅ Correctly returns an interface implementation.
	}
	return nil
}

// MockCollectionEmpty simulates a MongoDB collection that returns no rides.
type MockCollectionEmpty struct{}

// Find returns an empty cursor to simulate no available rides.
func (m *MockCollectionEmpty) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	// Ensure it returns a valid empty slice
	emptyDocs := []interface{}{}
	return mongo.NewCursorFromDocuments(emptyDocs, nil, bson.DefaultRegistry)
}
