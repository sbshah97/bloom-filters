package bloom

import "math"

// OptimalSize calculates the optimal size of the Bloom filter
func OptimalSize(expectedElements int, falsePositiveRate float64) uint {
	size := uint(math.Ceil(-float64(expectedElements) * math.Log(falsePositiveRate) / math.Pow(math.Log(2), 2)))
	return size
}

// OptimalHashFunctions calculates the optimal number of hash functions
func OptimalHashFunctions(size uint, expectedElements int) uint {
	numHash := uint(math.Ceil(float64(size) / float64(expectedElements) * math.Log(2)))
	return numHash
}
