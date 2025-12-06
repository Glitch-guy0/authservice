package benchmark

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// BenchmarkHealthEndpoint benchmarks the health check endpoint
func BenchmarkHealthEndpoint(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": "2025-01-01T00:00:00Z",
			"version":   "1.0.0",
		})
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}

// BenchmarkMiddleware benchmarks the middleware stack
func BenchmarkMiddleware(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Next()
	})
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}

// BenchmarkJSONResponse benchmarks JSON response generation
func BenchmarkJSONResponse(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.GET("/json", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"user": gin.H{
				"id":     "123",
				"name":   "Test User",
				"email":  "test@example.com",
				"active": true,
				"metadata": map[string]interface{}{
					"last_login": "2025-01-01T00:00:00Z",
					"preferences": map[string]string{
						"theme": "dark",
						"lang":  "en",
					},
				},
			},
		})
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/json", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}
