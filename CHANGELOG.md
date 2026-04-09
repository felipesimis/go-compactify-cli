# Changelog

All notable changes to this project will be documented here.

## [1.0.0] - 2026-04-09

### 🚀 Core Capabilities
- **High-Performance Processing**: Full image optimization pipeline powered by `libvips`.
- **Hardware-Aware Concurrency**: Intelligent resource management for heavy workloads.
- **Safe Execution**: Built-in `DryRun` mode to prevent accidental filesystem changes.

### 🛠 Engineering & Architecture
- **Decoupled Filesystem**: Implementation of the `OSOperations` interface, allowing for complete isolation of side effects and highly reliable unit testing.
- **Robust Error Handling**: Improved error aggregation and specific error types for better debugging.
- **Optimized Testing**: Transitioned to a sophisticated mocking architecture, ensuring core logic is thoroughly validated.
