package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAssertionHelper_AssertJSONEqual(t *testing.T) {
	h := NewAssertionHelper(t)

	t.Run("equal JSON objects", func(t *testing.T) {
		h.AssertJSONEqual(`{"key": "value"}`, `{ "key": "value" }`)
	})

	t.Run("equal JSON arrays", func(t *testing.T) {
		h.AssertJSONEqual(`[1, 2, 3]`, `[1, 2, 3]`)
	})
}

func TestAssertionHelper_AssertJSONContains(t *testing.T) {
	h := NewAssertionHelper(t)
	jsonStr := `{"name":"test", "value": 123, "nested": {"key": "value"}}`

	h.AssertJSONContains(jsonStr, map[string]interface{}{
		"name":  "test",
		"value": float64(123),
	})
}

func TestAssertionHelper_AssertHTTPStatus(t *testing.T) {
	h := NewAssertionHelper(t)
	rr := httptest.NewRecorder()
	rr.WriteHeader(http.StatusOK)

	h.AssertHTTPStatus(rr.Result(), http.StatusOK)
}

func TestAssertionHelper_AssertHTTPHeader(t *testing.T) {
	h := NewAssertionHelper(t)
	rr := httptest.NewRecorder()
	rr.Header().Set("Content-Type", "application/json")

	h.AssertHTTPHeader(rr.Result(), "Content-Type", "application/json")
}

func TestAssertionHelper_AssertHTTPContentType(t *testing.T) {
	h := NewAssertionHelper(t)
	rr := httptest.NewRecorder()
	rr.Header().Set("Content-Type", "application/json; charset=utf-8")

	h.AssertHTTPContentType(rr.Result(), "application/json")
}

func TestAssertionHelper_AssertSliceContains(t *testing.T) {
	h := NewAssertionHelper(t)
	slice := []string{"a", "b", "c"}
	h.AssertSliceContains(slice, "b")
}

func TestAssertionHelper_AssertMapContainsKey(t *testing.T) {
	h := NewAssertionHelper(t)
	m := map[string]int{"a": 1, "b": 2}
	h.AssertMapContainsKey(m, "a")
}

func TestAssertionHelper_AssertMapValue(t *testing.T) {
	h := NewAssertionHelper(t)
	m := map[string]int{"a": 1, "b": 2}
	h.AssertMapValue(m, "a", 1)
}

func TestAssertionHelper_AssertStringNotEmpty(t *testing.T) {
	h := NewAssertionHelper(t)
	h.AssertStringNotEmpty("not empty")
}

func TestAssertionHelper_AssertStringLength(t *testing.T) {
	h := NewAssertionHelper(t)
	h.AssertStringLength("four", 4)
}
