# Changelog

## [1.2.0] - 2026-04-27

### 🎨 User Interface & Experience
- **Terminal UI overhaul**: Migrated the entire output system to Lipgloss, providing a modern, high-fidelity command-line experience.

- **Enhanced Visual Hierarchy**: Implemented a sophisticated error reporting system that distinguishes between file paths and error causes through color coding and bold typography.

- **Smart Result Dashboard**: Added a side-by-side impact analysis panel (Original vs. Processed) with automatic highlighting of performance gains.

- **Intelligent Visibility**: The UI now automatically hides redundant information, such as the "Skipped" count when no images were bypassed, ensuring a cleaner output.

### 🛠 Architecture & Refactoring
- **Domain-UI Decoupling**: Completely separated the ResultBuilder from the presentation layer. The core logic now yields a pure Data Transfer Object (DTO), while the cmd package handles orchestration.

- **Component-Based Rendering**: Introduced a reusable UI package with isolated components (Panel, Dashboard, ErrorList), significantly improving maintainability.

- **Advanced TDD Coverage**: Achieved high test coverage for visual components, ensuring that UI styles and data mappings are strictly validated.

## [1.1.0] - 2026-04-09

### 🐛 Bug Fixes
- **Concurrency Safety**: Prevented potential deadlocks when the concurrency flag is explicitly set to 0.

### 🛠 Engineering & Refactoring
- **Test Suite Modernization**: Migrated all core assertions to the testify/suite syntax for better lifecycle management.
- **Improved Embeds**: Modernized image processing tests using Go's embed package for more reliable test data handling.
- **Recursive Support Prep**: Refactored BuildOutputPath to support future recursive directory tree implementations.

## [1.0.0] - 2026-04-09

### 🚀 Core Capabilities
- **High-Performance Processing**: Full image optimization pipeline powered by `libvips`.
- **Hardware-Aware Concurrency**: Intelligent resource management for heavy workloads.
- **Safe Execution**: Built-in `DryRun` mode to prevent accidental filesystem changes.

### 🛠 Engineering & Architecture
- **Decoupled Filesystem**: Implementation of the `OSOperations` interface, allowing for complete isolation of side effects and highly reliable unit testing.
- **Robust Error Handling**: Improved error aggregation and specific error types for better debugging.
- **Optimized Testing**: Transitioned to a sophisticated mocking architecture, ensuring core logic is thoroughly validated.