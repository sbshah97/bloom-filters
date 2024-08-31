package bloom

import "math"

// OptimalSize calculates the optimal size of the Bloom filter
func OptimalSize(expectedElements int, falsePositiveRate float64) uint {
	// This formula is the standard and widely accepted optimal size calculation for Bloom filters.
	// It's derived from probability theory and optimization principles. Here's why it's considered optimal:
	//
	// 1. It minimizes the probability of false positives given the number of elements and desired error rate.
	// 2. It balances the trade-off between the filter's size and its accuracy.
	// 3. It's based on the assumption that the hash functions are independent and uniformly distributed.
	//
	// While this is the most common formula, there are alternatives for specific use cases:
	//
	// - For very large datasets, some implementations use approximations to avoid potential integer overflow.
	// - In memory-constrained environments, a suboptimal smaller size might be chosen deliberately.
	// - For dynamic Bloom filters that can grow, different sizing strategies might be employed.
	//
	// However, for most standard Bloom filter implementations, this formula provides the best balance
	// of space efficiency and false positive rate, making it the de facto standard in the field.
	size := uint(math.Ceil(-float64(expectedElements) * math.Log(falsePositiveRate) / math.Pow(math.Log(2), 2)))
	return size
}

// OptimalHashFunctions calculates the optimal number of hash functions for a Bloom filter
func OptimalHashFunctions(size uint, expectedElements int) uint {
	// This formula calculates the optimal number of hash functions to minimize the false positive rate.
	// It's derived from probability theory and optimization principles:
	//
	// 1. It balances the trade-off between false positive rate and computational cost.
	// 2. Too few hash functions increase false positives, while too many slow down the filter.
	// 3. This formula minimizes the probability of false positives for a given filter size and number of elements.
	//
	// The formula (m/n * ln(2)) is optimal because:
	// - 'm' is the filter size (in bits)
	// - 'n' is the number of expected elements
	// - ln(2) ≈ 0.693 optimizes the bit setting probability to about 50%, which is ideal
	//
	// This ensures each bit has an equal probability of being set, maximizing the filter's efficiency.
	// While other formulas exist for specific use cases, this is generally considered the most optimal
	// for standard Bloom filter implementations, providing the best balance of accuracy and performance.
	numHash := uint(math.Ceil(float64(size) / float64(expectedElements) * math.Log(2)))
	return numHash
}
