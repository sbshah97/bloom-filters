# Contributing to Bloom Filters

We welcome contributions to the Bloom Filters project! This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository on GitHub.
2. Clone your fork locally.
3. Create a new branch for your feature or bug fix.
4. Make your changes and commit them with clear, descriptive commit messages.
5. Push your changes to your fork on GitHub.
6. Submit a pull request to the main repository.

## Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) for style guidelines.
- Use `gofmt` to format your code before submitting.
- Write clear, concise comments for functions and complex logic.

## Testing

- Add or update tests for any new features or bug fixes.
- Ensure all tests pass before submitting a pull request.
- Aim for high test coverage for new code.

## Pull Request Process

1. Ensure your code adheres to the style guidelines and passes all tests.
2. Update the README.md with details of changes to the interface, if applicable.
3. Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent.
4. Your pull request will be reviewed by maintainers, who may request changes or ask questions.

## Reporting Issues

- Use the GitHub issue tracker to report bugs or suggest features.
- Provide as much context as possible when reporting bugs, including your Go version, operating system, and steps to reproduce the issue.

Thank you for contributing to Bloom Filters!

- Aim for high test coverage for new code.
- The CI pipeline will automatically run on your pull request, checking for build errors, running tests, and performing linting. Ensure all checks pass before requesting a review.