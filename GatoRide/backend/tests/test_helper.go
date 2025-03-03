package tests

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ParseResponse - Helper function to parse HTTP response JSON
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder, result interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), result)
	assert.Nil(t, err, "Failed to parse JSON response")
}
