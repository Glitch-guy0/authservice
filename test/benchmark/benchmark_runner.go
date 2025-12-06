package benchmark

import (
"fmt"
"os"
"os/exec"
"path/filepath"
)

// RunAllBenchmarks runs all benchmark tests
func RunAllBenchmarks() error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	benchmarkDir := filepath.Join(projectRoot, "test", "benchmark")
	
	// Change to benchmark directory
	if err := os.Chdir(benchmarkDir); err != nil {
		return fmt.Errorf("failed to change to benchmark directory: %w", err)
	}

	// Run benchmarks
	cmd := exec.Command("go", "test", "-bench=.", "-benchmem", "-count=3")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run benchmarks: %w", err)
	}

	return nil
}

// RunSpecificBenchmark runs a specific benchmark test
func RunSpecificBenchmark(benchmarkName string) error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	benchmarkDir := filepath.Join(projectRoot, "test", "benchmark")
	
	// Change to benchmark directory
	if err := os.Chdir(benchmarkDir); err != nil {
		return fmt.Errorf("failed to change to benchmark directory: %w", err)
	}

	// Run specific benchmark
	cmd := exec.Command("go", "test", "-bench="+benchmarkName, "-benchmem", "-count=3")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run benchmark %s: %w", benchmarkName, err)
	}

	return nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("project root not found")
}
