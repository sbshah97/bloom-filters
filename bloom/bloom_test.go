package bloom

import (
	"bytes"
	"log/slog"
	"math"
	"math/rand"
	"os"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tests := []struct {
		name            string
		size            uint
		numHashFuncs    uint
		addElements     []string
		checkElements   []string
		expectedResults []bool
	}{
		{
			name:            "Basic functionality",
			size:            1000,
			numHashFuncs:    3,
			addElements:     []string{"hello", "world"},
			checkElements:   []string{"hello", "world", "golang"},
			expectedResults: []bool{true, true, false},
		},
		{
			name:            "Empty filter",
			size:            100,
			numHashFuncs:    2,
			addElements:     []string{},
			checkElements:   []string{"test"},
			expectedResults: []bool{false},
		},
		{
			name:            "Single element",
			size:            50,
			numHashFuncs:    1,
			addElements:     []string{"unique"},
			checkElements:   []string{"unique", "not_unique"},
			expectedResults: []bool{true, false},
		},
		{
			name:            "Large filter",
			size:            10000,
			numHashFuncs:    5,
			addElements:     []string{"a", "b", "c", "d", "e"},
			checkElements:   []string{"a", "c", "e", "f", "g"},
			expectedResults: []bool{true, true, true, false, false},
		},
		{
			name:            "Small filter (higher false positive rate)",
			size:            10,
			numHashFuncs:    2,
			addElements:     []string{"red", "green", "blue"},
			checkElements:   []string{"red", "green", "blue", "yellow", "purple"},
			expectedResults: []bool{true, true, true, true, false}, // Note: May have false positives
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
		name              string
		expectedElements  int
		falsePositiveRate float64
		expectedSize      uint
	}{
		{"Small set", 100, 0.01, 959},
		{"Medium set", 1000, 0.001, 14378},
		{"Large set", 1000000, 0.0001, 19170117},
		{"Edge case: very small set", 1, 0.1, 5},
		{"Edge case: very low false positive rate", 1000, 0.000001, 28756},
		{"Very large set", 1000000000, 0.00001, 23962645944},
		{"Extremely low false positive rate", 10000, 0.0000001, 335478},
		{"High false positive rate", 1000, 0.5, 1443},
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
		name             string
		size             uint
		expectedElements int
		expectedNumHash  uint
	}{
		{"Small filter", 1000, 100, 7},
		{"Medium filter", 10000, 1000, 7},
		{"Large filter", 100000, 10000, 7},
		{"Edge case: very small filter", 10, 1, 7},
		{"Edge case: size equals expected elements", 1000, 1000, 1},
		{"Very large filter", 1000000000, 100000000, 7},
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

func TestFalsePositiveRate(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tests := []struct {
		name          string
		size          uint
		numHashFuncs  uint
		addElements   []string
		checkElements []string
		expectedFPR   float64
		tolerance     float64
	}{
		{
			name:          "Small filter",
			size:          100,
			numHashFuncs:  3,
			addElements:   []string{"a", "b", "c", "d", "e"},
			checkElements: []string{"f", "g", "h", "i", "j"},
			expectedFPR:   0.0,
			tolerance:     0.5,
		},
		{
			name:          "Medium filter",
			size:          1000,
			numHashFuncs:  5,
			addElements:   generateRandomStrings(100, 10),
			checkElements: generateRandomStrings(10000, 10),
			expectedFPR:   0.03,
			tolerance:     0.02,
		},
		{
			name:          "Large filter",
			size:          10000,
			numHashFuncs:  7,
			addElements:   generateRandomStrings(1000, 20),
			checkElements: generateRandomStrings(10000, 20),
			expectedFPR:   0.001,
			tolerance:     0.0005,
		},
		{
			name:          "Very large filter",
			size:          1000000,
			numHashFuncs:  7,
			addElements:   generateRandomStrings(100000, 20),
			checkElements: generateRandomStrings(1000000, 20),
			expectedFPR:   0.0001,
			tolerance:     0.00005,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := NewBloomFilter(tt.size, tt.numHashFuncs, logger)

			// Add elements
			for _, elem := range tt.addElements {
				bf.Add([]byte(elem))
			}

			// Check elements and count false positives
			falsePositives := 0
			for _, elem := range tt.checkElements {
				if bf.Contains([]byte(elem)) {
					falsePositives++
				}
			}

			actualFPR := float64(falsePositives) / float64(len(tt.checkElements))
			calculatedFPR := bf.FalsePositiveRate()

			t.Logf("Actual FPR: %f, Calculated FPR: %f, Expected FPR: %f", actualFPR, calculatedFPR, tt.expectedFPR)

			if math.Abs(actualFPR-tt.expectedFPR) > tt.tolerance {
				t.Errorf("Actual false positive rate (%f) differs from expected (%f) by more than tolerance (%f)", actualFPR, tt.expectedFPR, tt.tolerance)
			}

			if math.Abs(calculatedFPR-actualFPR) > tt.tolerance {
				t.Errorf("Calculated false positive rate (%f) differs from actual (%f) by more than tolerance (%f)", calculatedFPR, actualFPR, tt.tolerance)
			}
		})
	}
}

// Helper function to generate random strings
func generateRandomStrings(count, length int) []string {
	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = randomString(length)
	}
	return result
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func TestSaveAndLoad(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tests := []struct {
		name          string
		size          uint
		numHashFuncs  uint
		addElements   []string
		checkElements []string
	}{
		{
			name:          "Small filter",
			size:          100,
			numHashFuncs:  3,
			addElements:   []string{"apple", "banana", "cherry"},
			checkElements: []string{"apple", "banana", "cherry", "date", "elderberry"},
		},
		{
			name:          "Medium filter",
			size:          1000,
			numHashFuncs:  5,
			addElements:   generateRandomStrings(50, 10),
			checkElements: generateRandomStrings(100, 10),
		},
		{
			name:          "Large filter",
			size:          10000,
			numHashFuncs:  7,
			addElements:   generateRandomStrings(1000, 20),
			checkElements: generateRandomStrings(2000, 20),
		},
		{
			name:          "Empty filter",
			size:          100,
			numHashFuncs:  3,
			addElements:   []string{},
			checkElements: []string{"test1", "test2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create and populate original filter
			originalBF := NewBloomFilter(tt.size, tt.numHashFuncs, logger)
			for _, elem := range tt.addElements {
				originalBF.Add([]byte(elem))
			}

			// Save the filter
			var buf bytes.Buffer
			err := originalBF.Save(&buf)
			if err != nil {
				t.Fatalf("Failed to save Bloom filter: %v", err)
			}

			// Load the filter
			loadedBF := &BloomFilter{}
			err = loadedBF.Load(&buf, logger)
			if err != nil {
				t.Fatalf("Failed to load Bloom filter: %v", err)
			}

			// Verify the loaded filter
			if loadedBF.size != originalBF.size {
				t.Errorf("Loaded filter size (%d) doesn't match original (%d)", loadedBF.size, originalBF.size)
			}
			if len(loadedBF.hashFuncs) != len(originalBF.hashFuncs) {
				t.Errorf("Loaded filter hash functions count (%d) doesn't match original (%d)", len(loadedBF.hashFuncs), len(originalBF.hashFuncs))
			}

			// Check elements
			for _, elem := range tt.checkElements {
				originalResult := originalBF.Contains([]byte(elem))
				loadedResult := loadedBF.Contains([]byte(elem))
				if originalResult != loadedResult {
					t.Errorf("Mismatch for element %s: original %v, loaded %v", elem, originalResult, loadedResult)
				}
			}

			// Verify false positive rates are similar
			originalFPR := originalBF.FalsePositiveRate()
			loadedFPR := loadedBF.FalsePositiveRate()
			if math.Abs(originalFPR-loadedFPR) > 1e-6 {
				t.Errorf("False positive rates differ: original %v, loaded %v", originalFPR, loadedFPR)
			}
		})
	}
}
