package bloom

import (
	"testing"
	"log/slog"
	"os"
)

func BenchmarkBloomFilterAdd(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	bf := NewBloomFilter(1000, 3, logger)
	element := []byte("benchmark")

	for i := 0; i < b.N; i++ {
		bf.Add(element)
	}
}

func BenchmarkBloomFilterContains(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	bf := NewBloomFilter(1000, 3, logger)
	element := []byte("benchmark")
	bf.Add(element)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Contains(element)
	}
}

func BenchmarkOptimalSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OptimalSize(1000, 0.01)
	}
}

func BenchmarkOptimalHashFunctions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OptimalHashFunctions(1000, 100)
	}
}