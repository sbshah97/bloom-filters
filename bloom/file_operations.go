package bloom

import (
	"log/slog"
	"os"
)

// SaveFilterToFile saves a Bloom filter to a file
func SaveFilterToFile(bf *BloomFilter, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return bf.Save(file)
}

// LoadFilterFromFile loads a Bloom filter from a file
func LoadFilterFromFile(filename string, logger *slog.Logger) (*BloomFilter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	loadedBF := &BloomFilter{}
	err = loadedBF.Load(file, logger)
	if err != nil {
		return nil, err
	}

	return loadedBF, nil
}
