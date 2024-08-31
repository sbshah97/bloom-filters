package bloom

import (
	"log/slog"
	"os"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tests := []struct {
		name             string
		size             uint
		numHashFuncs     uint
		addElements      []string
		checkElements    []string
		expectedResults  []bool
	}{
		{
			name:             "Basic functionality",
			size:             1000,
			numHashFuncs:     3,
			addElements:      []string{"hello", "world"},
			checkElements:    []string{"hello", "world", "golang"},
			expectedResults:  []bool{true, true, false},
		},
		{
			name:             "Empty filter",
			size:             100,
			numHashFuncs:     2,
			addElements:      []string{},
			checkElements:    []string{"test"},
			expectedResults:  []bool{false},
		},
		{
			name:             "Single element",
			size:             50,
			numHashFuncs:     1,
			addElements:      []string{"unique"},
			checkElements:    []string{"unique", "not_unique"},
			expectedResults:  []bool{true, false},
		},
		{
			name:             "Large filter",
			size:             10000,
			numHashFuncs:     5,
			addElements:      []string{"a", "b", "c", "d", "e"},
			checkElements:    []string{"a", "c", "e", "f", "g"},
			expectedResults:  []bool{true, true, true, false, false},
		},
		{
			name:             "Small filter (higher false positive rate)",
			size:             10,
			numHashFuncs:     2,
			addElements:      []string{"red", "green", "blue"},
			checkElements:    []string{"red", "green", "blue", "yellow", "purple"},
			expectedResults:  []bool{true, true, true, true, true}, // Note: May have false positives
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := NewBloomFilter(tt.size, tt.numHashFuncs, logger)

			// Add elements
			for _, elem := range tt.addElements {
				bf.Add([]byte(elem))
			}

			// Check elements
			for i, elem := range tt.checkElements {
				result := bf.Contains([]byte(elem))
				if result != tt.expectedResults[i] {
					t.Errorf("Expected Contains(%s) to be %v, but got %v", elem, tt.expectedResults[i], result)
				}
			}
		})
	}
}

func TestOptimalSize(t *testing.T) {
	tests := []struct {
		name               string
		expectedElements   int
		falsePositiveRate  float64
		expectedSize       uint
	}{
		{"Small set", 100, 0.01, 959},
		{"Medium set", 1000, 0.001, 14378},
		{"Large set", 1000000, 0.0001, 1917011},
		{"Edge case: very small set", 1, 0.1, 5},
		{"Edge case: very low false positive rate", 1000, 0.000001, 28755},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := OptimalSize(tt.expectedElements, tt.falsePositiveRate)
			if size != tt.expectedSize {
				t.Errorf("Expected optimal size to be %d, but got %d", tt.expectedSize, size)
			}
		})
	}
}

func TestOptimalHashFunctions(t *testing.T) {
	tests := []struct {
		name               string
		size               uint
		expectedElements   int
		expectedNumHash    uint
	}{
		{"Small filter", 1000, 100, 7},
		{"Medium filter", 10000, 1000, 7},
		{"Large filter", 100000, 10000, 7},
		{"Edge case: very small filter", 10, 1, 7},
		{"Edge case: size equals expected elements", 1000, 1000, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			numHash := OptimalHashFunctions(tt.size, tt.expectedElements)
			if numHash != tt.expectedNumHash {
				t.Errorf("Expected optimal number of hash functions to be %d, but got %d", tt.expectedNumHash, numHash)
			}
		})
	}
}
