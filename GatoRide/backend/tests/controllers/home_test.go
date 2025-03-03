package controllers_test

// import (
// 	"backend/controllers"
// 	"backend/models"
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // ✅ Mock User Collection (Simulating MongoDB)
// type MockUserCollection struct {
// 	UserData *models.User
// 	FindErr  error
// }

// func (m *MockUserCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
// 	if m.FindErr != nil {
// 		return &mongo.SingleResult{}
// 	}
// 	if m.UserData == nil {
// 		return &mongo.SingleResult{}
// 	}
// 	return &mongo.SingleResult{} // Simulating a valid response
// }

// // ✅ Mock Rides Collection (Simulating MongoDB)
// type MockRidesCollection struct {
// 	RidesData []models.Ride
// 	FindErr   error
// }

// func (m *MockRidesCollection) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
// 	if m.FindErr != nil {
// 		return nil, m.FindErr
// 	}
// 	var docs []interface{}
// 	for _, ride := range m.RidesData {
// 		docs = append(docs, ride)
// 	}
// 	return mongo.NewCursorFromDocuments(docs, nil, nil)
// }

// // ✅ Test Setup: Gin Router with Mocks
// func setupRouter(mockUser *MockUserCollection, mockRides *MockRidesCollection) *gin.Engine {
// 	gin.SetMode(gin.TestMode)
// 	router := gin.Default()

// 	router.GET("/home", func(c *gin.Context) {
// 		// Simulate middleware injecting userID
// 		userID := primitive.NewObjectID().Hex()
// 		c.Set("userID", userID)

// 		// Call HomeHandler with mock collections
// 		controllers.HomeHandler(c)
// 	})

// 	return router
// }

// // ✅ Test 1: Valid User with Location Set
// func TestHomeHandler_ValidUser(t *testing.T) {
// 	mockUser := &MockUserCollection{
// 		UserData: &models.User{
// 			ID: primitive.NewObjectID(),
// 			Location: models.Location{
// 				Latitude:  40.7128,
// 				Longitude: -74.0060,
// 			},
// 		},
// 	}
// 	mockRides := &MockRidesCollection{
// 		RidesData: []models.Ride{
// 			{Pickup: models.Location{Latitude: 40.7128, Longitude: -74.0060}, Status: "open"},
// 		},
// 	}

// 	router := setupRouter(mockUser, mockRides)

// 	req, _ := http.NewRequest("GET", "/home", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")
// }

// // ✅ Test 2: Missing User ID
// func TestHomeHandler_MissingUserID(t *testing.T) {
// 	mockUser := &MockUserCollection{}
// 	mockRides := &MockRidesCollection{}

// 	router := setupRouter(mockUser, mockRides)

// 	req, _ := http.NewRequest("GET", "/home", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status 401 Unauthorized")
// }

// // ✅ Test 3: User Not Found
// func TestHomeHandler_UserNotFound(t *testing.T) {
// 	mockUser := &MockUserCollection{UserData: nil}
// 	mockRides := &MockRidesCollection{}

// 	router := setupRouter(mockUser, mockRides)

// 	req, _ := http.NewRequest("GET", "/home", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusNotFound, w.Code, "Expected status 404 Not Found")
// }
