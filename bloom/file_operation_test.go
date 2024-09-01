package bloom

import (
	"errors"
	"hash"
	"log/slog"
	"os"
	"testing"
)

func TestSaveFilterToFile(t *testing.T) {
	defer cleanup(t)
	bf := &Filter{
		bitArray:  []bool{true, false, true},
		size:      3,
		hashFuncs: make([]hash.Hash64, 2),
		logger:    slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	tests := []struct {
		name          string
		filename      string
		expectedError error
	}{
		{
			name:          "Successful save",
			filename:      "test.gob",
			expectedError: nil,
		},
		{
			name:          "File creation error",
			filename:      "/nonexistent/error.gob",
			expectedError: os.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

			err := SaveFilterToFile(bf, tt.filename, logger)

			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && !errors.Is(err, tt.expectedError)) {
				t.Errorf("SaveFilterToFile() error = %v, expectedError %v", err, tt.expectedError)
			}

			// Clean up the file if it was created
			if tt.expectedError == nil {
				if err := os.Remove(tt.filename); err != nil {
					t.Errorf("Failed to remove file %s: %v", tt.filename, err)
				}
			}
		})
	}
}

func TestLoadFilterFromFile(t *testing.T) {
	defer cleanup(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a valid Bloom filter file for testing
	validFilter := &Filter{
		bitArray:  []bool{true, false, true},
		size:      3,
		hashFuncs: make([]hash.Hash64, 2),
		logger:    logger,
	}
	validFilename := "valid_test.gob"
	err := SaveFilterToFile(validFilter, validFilename, logger)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() {
		if err := os.Remove(validFilename); err != nil {
			t.Errorf("Failed to remove file %s: %v", validFilename, err)
		}
	}()

	tests := []struct {
		name          string
		filename      string
		expectedError error
	}{
		{
			name:          "Successful load",
			filename:      validFilename,
			expectedError: nil,
		},
		{
			name:          "File open error",
			filename:      "nonexistent.gob",
			expectedError: os.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := LoadFilterFromFile(tt.filename, logger)

			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && !errors.Is(err, tt.expectedError)) {
				t.Errorf("LoadFilterFromFile() error = %v, expectedError %v", err, tt.expectedError)
			}

			if err == nil && filter == nil {
				t.Errorf("LoadFilterFromFile() returned nil filter")
			}
		})
	}
}

func cleanup(t *testing.T) {
	files := []string{"test.gob", "error.gob", "write_error.gob", "read_error.gob", "valid_test.gob"}
	for _, file := range files {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			t.Logf("Failed to remove file %s: %v", file, err)
		}
	}
}
