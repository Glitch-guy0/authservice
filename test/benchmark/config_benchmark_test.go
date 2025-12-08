package benchmark

import (
	"testing"

	"github.com/Glitch-guy0/authService/modules/core/config"
)

// BenchmarkConfigLoad benchmarks configuration loading
func BenchmarkConfigLoad(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := config.Init("")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConfigGet benchmarks configuration value retrieval
func BenchmarkConfigGet(b *testing.B) {
	err := config.Init("")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = config.Config.GetString("server.port")
			_ = config.Config.GetString("log.level")
			_ = config.Config.GetBool("server.debug")
		}
	})
}

// BenchmarkConfigValidation benchmarks configuration validation
func BenchmarkConfigValidation(b *testing.B) {
	err := config.Init("")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Since there's no explicit Validate function, we'll just access config to trigger validation
		_ = config.Config.GetString("server.port")
	}
}
