package bloom

import (
	"log/slog"
	"os"
)

// SaveFilterToFile saves a Bloom filter to a file
func SaveFilterToFile(bf *Filter, filename string, logger *slog.Logger) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("Failed to close file", "error", err)
		}
	}()

	return bf.Save(file)
}

// LoadFilterFromFile loads a Bloom filter from a file
func LoadFilterFromFile(filename string, logger *slog.Logger) (*Filter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("Failed to close file", "error", err)
		}
	}()

	loadedBF := &Filter{}
	err = loadedBF.Load(file, logger)
	if err != nil {
		return nil, err
	}

	return loadedBF, nil
}
