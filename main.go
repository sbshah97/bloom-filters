package main

import (
	"log/slog"
	"os"

	"github.com/sbshah97/bloom-filters/bloom"
)

const (
	expectedElements  = 1000
	falsePositiveRate = 0.01
	bloomFilterFile   = "bloom_filter.gob"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func main() {
	bf := createAndPopulateFilter()
	checkFilter(bf)

	err := saveAndLoadFilter(bf)
	if err != nil {
		logger.Error("Error in save and load process", "error", err)
		return
	}

	cleanup()
}

func createAndPopulateFilter() *bloom.BloomFilter {
	size := bloom.OptimalSize(expectedElements, falsePositiveRate)
	numHashFuncs := bloom.OptimalHashFunctions(size, expectedElements)
	bf := bloom.NewBloomFilter(size, numHashFuncs, logger)

	bf.Add([]byte("hello"))
	bf.Add([]byte("world"))
	return bf
}

func checkFilter(bf *bloom.BloomFilter) {
	logger.Info("Checking Bloom filter",
		"contains_hello", bf.Contains([]byte("hello")),
		"contains_world", bf.Contains([]byte("world")),
		"contains_golang", bf.Contains([]byte("golang")))
}

func saveAndLoadFilter(bf *bloom.BloomFilter) error {
	if err := bloom.SaveFilterToFile(bf, bloomFilterFile); err != nil {
		return err
	}
	logger.Info("Bloom filter saved to file")

	loadedBF, err := bloom.LoadFilterFromFile(bloomFilterFile, logger)
	if err != nil {
		return err
	}
	logger.Info("Bloom filter loaded from file")

	checkFilter(loadedBF)
	return nil
}

func cleanup() {
	if err := os.Remove(bloomFilterFile); err != nil {
		logger.Error("Failed to delete file", "error", err)
		return
	}
	logger.Info("Bloom filter file deleted")
}
