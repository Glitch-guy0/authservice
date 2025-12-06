package benchmark

import (
	"context"
	"testing"

	"github.com/Glitch-guy0/authService/modules/logger"
)

// BenchmarkLoggerInfo benchmarks the Info method of the logger
func BenchmarkLoggerInfo(b *testing.B) {
	log := logger.New()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("Test message for benchmarking")
		}
	})
}

// BenchmarkLoggerError benchmarks the Error method of the logger
func BenchmarkLoggerError(b *testing.B) {
	log := logger.New()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("Test error message for benchmarking")
		}
	})
}

// BenchmarkLoggerWithFields benchmarks logging with fields
func BenchmarkLoggerWithFields(b *testing.B) {
	log := logger.New()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.WithField("user_id", "12345").
				WithField("action", "login").
				Info("User action logged")
		}
	})
}

// BenchmarkLoggerWithContext benchmarks logging with context
func BenchmarkLoggerWithContext(b *testing.B) {
	log := logger.NewContextual()
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.WithContext(ctx).Info("Contextual log message")
		}
	})
}
