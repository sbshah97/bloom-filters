package main

import (
	"log/slog"
	"os"

	"github.com/sbshah97/bloom-filters/bloom"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	size := bloom.OptimalSize(1000, 0.01)
	numHashFuncs := bloom.OptimalHashFunctions(size, 1000)

	bf := bloom.NewBloomFilter(size, numHashFuncs, logger)

	bf.Add([]byte("hello"))
	bf.Add([]byte("world"))

	logger.Info("Checking Bloom filter", "contains_hello", bf.Contains([]byte("hello")))
	logger.Info("Checking Bloom filter", "contains_world", bf.Contains([]byte("world")))
	logger.Info("Checking Bloom filter", "contains_golang", bf.Contains([]byte("golang")))
}
