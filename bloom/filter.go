package bloom

import (
	"encoding/gob"
	"hash"
	"hash/fnv"
	"io"
	"log/slog"
	"math"
)

// BloomFilter represents a Bloom filter data structure
type BloomFilter struct {
	bitArray  []bool
	size      uint
	hashFuncs []hash.Hash64
	logger    *slog.Logger
}

// NewBloomFilter creates a new Bloom filter with the given size and number of hash functions
func NewBloomFilter(size uint, numHashFuncs uint, logger *slog.Logger) *BloomFilter {
	bf := &BloomFilter{
		bitArray:  make([]bool, size),
		size:      size,
		hashFuncs: make([]hash.Hash64, numHashFuncs),
		logger:    logger,
	}

	for i := uint(0); i < numHashFuncs; i++ {
		bf.hashFuncs[i] = fnv.New64()
	}

	bf.logger.Info("Created new Bloom filter", "size", size, "numHashFuncs", numHashFuncs)
	return bf
}

// Add adds an element to the Bloom filter
func (bf *BloomFilter) Add(element []byte) {
	// This loop iterates through all hash functions in the Bloom filter
	for i, h := range bf.hashFuncs {
		// Reset the hash function to clear any previous data
		h.Reset()
		bf.logger.Debug("Reset hash function", "hashFunc", i)

		// Write the element (as bytes) to the hash function
		h.Write(element)
		bf.logger.Debug("Wrote element to hash function", "hashFunc", i, "element", string(element))

		// Calculate the index in the bit array
		index := h.Sum64() % uint64(bf.size)
		bf.logger.Debug("Calculated index", "hashFunc", i, "index", index)

		// Set the bit at the calculated index to true
		bf.bitArray[index] = true
		bf.logger.Debug("Set bit in array", "hashFunc", i, "index", index)

		// Sample output for each step (assuming element is "hello" and bf.size is 10):
		// Step 1 (i=0): index might be 7, bf.bitArray becomes [0 0 0 0 0 0 0 1 0 0]
		// Step 2 (i=1): index might be 2, bf.bitArray becomes [0 0 1 0 0 0 0 1 0 0]
		// Step 3 (i=2): index might be 7 again, bf.bitArray stays [0 0 1 0 0 0 0 1 0 0]
		// ... and so on for each hash function
	}
	// After all hash functions, bf.bitArray might look like [0 0 1 0 1 0 0 1 1 0]
	// This means bits at indices 2, 4, 7, and 8 are set for the element "hello"
	bf.logger.Info("Added element to Bloom filter", "element", string(element))
}

// Contains checks if an element might be in the Bloom filter
func (bf *BloomFilter) Contains(element []byte) bool {
	for i, h := range bf.hashFuncs {
		h.Reset()
		h.Write(element)
		index := h.Sum64() % uint64(bf.size)
		if !bf.bitArray[index] {
			bf.logger.Debug("Element not found in Bloom filter", "element", string(element), "hashFunc", i)
			return false
		}
	}
	bf.logger.Info("Element possibly in Bloom filter", "element", string(element))
	return true
}

// FalsePositiveRate calculates the current false positive rate of the Bloom filter
func (bf *BloomFilter) FalsePositiveRate() float64 {
	setBits := 0
	for _, bit := range bf.bitArray {
		if bit {
			setBits++
		}
	}
	probability := float64(setBits) / float64(bf.size)
	return math.Pow(probability, float64(len(bf.hashFuncs)))
}

// Save serializes the Bloom filter to a writer
func (bf *BloomFilter) Save(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(struct {
		BitArray []bool
		Size     uint
		NumHash  uint
	}{
		BitArray: bf.bitArray,
		Size:     bf.size,
		NumHash:  uint(len(bf.hashFuncs)),
	})
}

// Load deserializes the Bloom filter from a reader
func (bf *BloomFilter) Load(r io.Reader, logger *slog.Logger) error {
	decoder := gob.NewDecoder(r)
	var data struct {
		BitArray []bool
		Size     uint
		NumHash  uint
	}
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	bf.bitArray = data.BitArray
	bf.size = data.Size
	bf.hashFuncs = make([]hash.Hash64, data.NumHash)
	for i := uint(0); i < data.NumHash; i++ {
		bf.hashFuncs[i] = fnv.New64()
	}
	bf.logger = logger
	return nil
}
