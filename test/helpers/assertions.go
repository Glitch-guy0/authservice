// Package helpers provides assertion utilities for testing
package helpers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertionHelper provides custom assertion methods
type AssertionHelper struct {
	t *testing.T
}

// NewAssertionHelper creates a new assertion helper
func NewAssertionHelper(t *testing.T) *AssertionHelper {
	return &AssertionHelper{t: t}
}

// AssertJSONEqual asserts that two JSON strings are equal
func (a *AssertionHelper) AssertJSONEqual(expected, actual string) {
	var expectedJSON, actualJSON interface{}

	err := json.Unmarshal([]byte(expected), &expectedJSON)
	require.NoError(a.t, err, "Failed to parse expected JSON")

	err = json.Unmarshal([]byte(actual), &actualJSON)
	require.NoError(a.t, err, "Failed to parse actual JSON")

	assert.Equal(a.t, expectedJSON, actualJSON, "JSON objects are not equal")
}

// AssertJSONContains asserts that a JSON string contains specific fields
func (a *AssertionHelper) AssertJSONContains(jsonStr string, expectedFields map[string]interface{}) {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	require.NoError(a.t, err, "Failed to parse JSON")

	for key, expectedValue := range expectedFields {
		actualValue, exists := jsonData[key]
		assert.True(a.t, exists, "Expected field '%s' not found in JSON", key)
		if exists {
			assert.Equal(a.t, expectedValue, actualValue, "Field '%s' value mismatch", key)
		}
	}
}

// AssertHTTPStatus asserts HTTP response status code
func (a *AssertionHelper) AssertHTTPStatus(resp *http.Response, expectedStatus int) {
	assert.Equal(a.t, expectedStatus, resp.StatusCode,
		"Expected HTTP status %d, got %d", expectedStatus, resp.StatusCode)
}

// AssertHTTPHeader asserts HTTP response header
func (a *AssertionHelper) AssertHTTPHeader(resp *http.Response, key, expectedValue string) {
	actualValue := resp.Header.Get(key)
	assert.Equal(a.t, expectedValue, actualValue,
		"Expected header '%s' to be '%s', got '%s'", key, expectedValue, actualValue)
}

// AssertHTTPContentType asserts HTTP response content type
func (a *AssertionHelper) AssertHTTPContentType(resp *http.Response, expectedContentType string) {
	contentType := resp.Header.Get("Content-Type")
	assert.True(a.t, strings.Contains(contentType, expectedContentType),
		"Expected Content-Type to contain '%s', got '%s'", expectedContentType, contentType)
}

// AssertSliceContains asserts that a slice contains an element
func (a *AssertionHelper) AssertSliceContains(slice interface{}, element interface{}) {
	sliceValue := reflect.ValueOf(slice)
	elementValue := reflect.ValueOf(element)

	found := false
	for i := 0; i < sliceValue.Len(); i++ {
		if reflect.DeepEqual(sliceValue.Index(i).Interface(), elementValue.Interface()) {
			found = true
			break
		}
	}

	assert.True(a.t, found, "Slice does not contain element: %v", element)
}

// AssertMapContainsKey asserts that a map contains a key
func (a *AssertionHelper) AssertMapContainsKey(m interface{}, key interface{}) {
	mapValue := reflect.ValueOf(m)
	keyValue := reflect.ValueOf(key)

	assert.True(a.t, mapValue.MapIndex(keyValue).IsValid(),
		"Map does not contain key: %v", key)
}

// AssertMapValue asserts that a map has a specific value for a key
func (a *AssertionHelper) AssertMapValue(m interface{}, key, expectedValue interface{}) {
	mapValue := reflect.ValueOf(m)
	keyValue := reflect.ValueOf(key)

	actualValue := mapValue.MapIndex(keyValue)
	assert.True(a.t, actualValue.IsValid(), "Map does not contain key: %v", key)

	if actualValue.IsValid() {
		assert.Equal(a.t, expectedValue, actualValue.Interface(),
			"Map value for key '%v' mismatch", key)
	}
}

// AssertStringNotEmpty asserts that a string is not empty
func (a *AssertionHelper) AssertStringNotEmpty(str string) {
	assert.NotEmpty(a.t, str, "String should not be empty")
}

// AssertStringLength asserts string length
func (a *AssertionHelper) AssertStringLength(str string, expectedLength int) {
	assert.Equal(a.t, expectedLength, len(str),
		"Expected string length %d, got %d", expectedLength, len(str))
}
