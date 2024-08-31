package main

import (
	"log/slog"
	"os"

	"github.com/sbshah97/bloom-filters/bloom"
)

const expectedElements = 1000
const falsePositiveRate = 0.01

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	size := bloom.OptimalSize(expectedElements, falsePositiveRate)
	numHashFuncs := bloom.OptimalHashFunctions(size, expectedElements)

	bf := bloom.NewBloomFilter(size, numHashFuncs, logger)

	bf.Add([]byte("hello"))
	bf.Add([]byte("world"))

	logger.Info("Checking Bloom filter", "contains_hello", bf.Contains([]byte("hello")))
	logger.Info("Checking Bloom filter", "contains_world", bf.Contains([]byte("world")))
	logger.Info("Checking Bloom filter", "contains_golang", bf.Contains([]byte("golang")))
}
