// Package testutils provides HTTP testing utilities
package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

// HTTPTestHelper provides utilities for HTTP testing
type HTTPTestHelper struct {
	router *gin.Engine
}

// NewHTTPTestHelper creates a new HTTP test helper
func NewHTTPTestHelper(router *gin.Engine) *HTTPTestHelper {
	return &HTTPTestHelper{router: router}
}

// Request performs an HTTP request and returns the response
func (h *HTTPTestHelper) Request(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	h.router.ServeHTTP(w, req)

	return w
}

// Get performs a GET request
func (h *HTTPTestHelper) Get(path string, headers map[string]string) *httptest.ResponseRecorder {
	return h.Request(http.MethodGet, path, nil, headers)
}

// Post performs a POST request
func (h *HTTPTestHelper) Post(path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	return h.Request(http.MethodPost, path, body, headers)
}

// Put performs a PUT request
func (h *HTTPTestHelper) Put(path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	return h.Request(http.MethodPut, path, body, headers)
}

// Delete performs a DELETE request
func (h *HTTPTestHelper) Delete(path string, headers map[string]string) *httptest.ResponseRecorder {
	return h.Request(http.MethodDelete, path, nil, headers)
}

// AssertJSONResponse asserts that the response contains valid JSON
func (h *HTTPTestHelper) AssertJSONResponse(resp *httptest.ResponseRecorder) map[string]interface{} {
	if resp.Header().Get("Content-Type") != "application/json" {
		panic("Expected Content-Type to be application/json")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		panic(err)
	}

	return result
}

// AssertStatus asserts the HTTP status code
func (h *HTTPTestHelper) AssertStatus(resp *httptest.ResponseRecorder, expectedStatus int) {
	if resp.Code != expectedStatus {
		panic("Expected status " + string(rune(expectedStatus)) + ", got " + string(rune(resp.Code)))
	}
}

// AssertBodyContains asserts that the response body contains the expected string
func (h *HTTPTestHelper) AssertBodyContains(resp *httptest.ResponseRecorder, expected string) {
	if !strings.Contains(resp.Body.String(), expected) {
		panic("Response body does not contain expected string: " + expected)
	}
}
