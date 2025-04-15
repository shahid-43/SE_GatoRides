package mocks

import (
	"backend/models"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ✅ CollectionInterface for Mocking MongoDB
type CollectionInterfaces interface {
	FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
}

// ✅ Mock User Collection
type MockUserCollection struct {
	UserData *models.User
	FindErr  error
}

// ✅ FindOne Simulation
func (m *MockUserCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	if m.FindErr != nil {
		return &mongo.SingleResult{}
	}
	if m.UserData == nil {
		return &mongo.SingleResult{}
	}
	return &mongo.SingleResult{} // Simulate valid response
}

// ✅ Mock Rides Collection
type MockRidesCollection struct {
	RidesData []models.Ride
	FindErr   error
}

// ✅ Find Simulation
func (m *MockRidesCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	if m.FindErr != nil {
		return nil, m.FindErr
	}
	var docs []interface{}
	for _, ride := range m.RidesData {
		docs = append(docs, ride)
	}
	return mongo.NewCursorFromDocuments(docs, nil, nil)
}
