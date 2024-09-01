# Bloom Filters

[![Go Report Card](https://goreportcard.com/badge/github.com/sbshah97/bloom-filters)](https://goreportcard.com/report/github.com/sbshah97/bloom-filters)
[![GoDoc](https://godoc.org/github.com/sbshah97/bloom-filters?status.svg)](https://godoc.org/github.com/sbshah97/bloom-filters)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go implementation of Bloom Filters, a space-efficient probabilistic data structure used to test whether an element is a member of a set.

## Features

- Create Bloom filters with customizable size and number of hash functions
- Add elements to the filter
- Check for element membership
- Calculate false positive rate
- Save and load Bloom filters to/from files

## Installation

To use this Bloom Filter implementation in your Go project:

```bash
go get github.com/sbshah97/bloom-filters
```

```go
package main
import (
    "fmt"
    "log/slog"
    "os"
    "github.com/sbshah97/bloom-filters/bloom"
)
func main() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    bf := bloom.NewBloomFilter(1000, 3, logger)
    bf.Add([]byte("hello"))
    bf.Add([]byte("world"))
    fmt.Println(bf.Contains([]byte("hello"))) // true
    fmt.Println(bf.Contains([]byte("world"))) // true
    fmt.Println(bf.Contains([]byte("golang"))) // false (probably)
}
```


For more detailed usage examples, please refer to the `main.go` file in the project root.

## Project Structure

- `bloom/filter.go`: Core implementation of the Bloom Filter
- `bloom/optimal.go`: Functions for calculating optimal Bloom Filter parameters
- `bloom/file_operations.go`: Functions for saving and loading Bloom Filters
- `bloom/*_test.go`: Unit tests for the Bloom Filter implementation
- `main.go`: Example usage of the Bloom Filter

## Running Tests

To run the tests for this project:

## Setup Instructions

1. Ensure you have Go installed on your system. If not, download and install it from [golang.org](https://golang.org/).

2. Clone this repository:
   ```
   git clone https://github.com/sbshah97/bloom-filters.git
   cd bloom-filters
   ```

3. Run the tests:
   ```
   go test
   ```

4. Run the main program:
   ```
   go run main.go
   ```

## Step-by-Step Implementation Guide

For those new to Bloom Filters, this project is structured to guide you through the implementation process with progressive test cases. Follow these steps to understand and implement a Bloom Filter:

- [x] **Basic Structure**: 
   - Implement the basic structure of a Bloom Filter (array of bits, hash functions).
   - Test: Create a Bloom Filter with a given size and number of hash functions.

- [x] **Add Operation**: 
   - Implement the method to add an element to the Bloom Filter.
   - Test: Add an element and verify that the corresponding bits are set.

- [x] **Query Operation**: 
   - Implement the method to check if an element might be in the set.
   - Test: Query for added elements (should return true) and non-added elements (should mostly return false).

4. **False Positive Rate**: 
   - Calculate and verify the false positive rate.
   - Test: Add a set of elements, then query for a different set and measure the false positive rate.

5. **Optimization**: 
   - Implement methods to optimize the size of the filter and number of hash functions based on the expected number of elements and desired false positive rate.
   - Test: Create filters with different parameters and compare their effectiveness.

6. **Serialization** (Optional):
   - Implement methods to save and load the Bloom Filter state.
   - Test: Save a filter's state, load it, and verify that queries still work correctly.

Each step builds upon the previous one, allowing you to gradually understand and implement the Bloom Filter concept. The accompanying test file (`bloom_filter_test.go`) will contain test cases for each step, helping you verify your implementation as you progress.

By following this guide and running the tests at each step, you'll gain a comprehensive understanding of how Bloom Filters work and how to implement them efficiently in Go.


## Contributing

Contributions are welcome! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines on how to contribute to this project.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## TODO

- Fix failing tests in `file_operation_test.go`
- Update mock usage in tests to match new function signatures
- Improve error handling and logging
- Add more comprehensive examples and documentation

## Acknowledgments

- [Bloom Filter concept](https://en.wikipedia.org/wiki/Bloom_filter)
- Go community for providing excellent tools and libraries

## Contact

For any questions or concerns, please open an issue on the GitHub repository.
