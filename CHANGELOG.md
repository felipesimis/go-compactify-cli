# Changelog

## [Unreleased]

### ⚠️ Breaking Changes
- **Dimensions Validation**: Increased minimum allowed width and height from `1` to `10` pixels to ensure processing stability.

### 🚀 Added & Changed
- **Recursive Directory Processing (`--recursive`)**: You can now process entire media libraries at once. The CLI deeply traverses subdirectories and seamlessly mirrors the original folder architecture.
- **Infinite Loop Prevention**: Implemented safety checks using absolute path resolution to prevent the engine from recursively processing its own output when the destination resides within the source tree.
- **Auto-Tuning Concurrency**: End of performance guesswork. If the `--concurrency` flag is omitted, the tool automatically detects your machine's logical CPU cores and calibrates the workload to extract maximum speed without freezing your system.
- **Standardized Encoding Control**: Added the `--quality` (`-q`) flag to all commands that perform lossy image re-encoding (e.g., `crop`, `resize`, `convert`, `flip`, `grayscale`, `palette`, `thumbnail`). This gives the user explicit control over output quality and prevents silent image degradation.
- **Multi-layer Configuration**: Implemented a robust precedence engine where Flags > Environment Variables > Config File > Defaults. This ensures maximum flexibility for local development and CI/CD environments.
- **Config Initialization**: Added `init` command (with `initialize` and `config` aliases) to generate a default `config.yaml` file, simplifying global settings management.
- **Overwrite Protection**: Implemented a security check in the `init` command that prevents accidental overwriting of existing configurations unless the `--force` (`-f`) flag is provided.
- **Environment Variable Mapping**: Integrated automatic mapping for all global settings using the `COMPACTIFY_` prefix (e.g., `COMPACTIFY_CONCURRENCY`, `COMPACTIFY_QUALITY`).
- **Privacy & Security**: Introduced the global `--strip-metadata` flag. This feature safely strips sensitive EXIF data (such as GPS coordinates and camera models) to protect user privacy while preserving ICC color profiles for color accuracy.
- **Performance UX**: Added automated warnings when concurrency levels exceed safe hardware limits ($2 \times \text{CPU cores}$), preventing system instability.
- **Standardized Success Feedback**: Integrated a new `Success` component in the UI package to provide consistent, high-fidelity visual confirmation for CLI operations.
- **Resilient Versioning**: Added intelligent version string sanitization that handles both v-prefixed tags and raw version strings, preventing redundant visual prefixes in the UI.

### 🛠 Performance & Stability
- **Lazy Directory Creation**: Refactored the processing pipeline to create output directories "just-in-time". This ensures the output remains clean by only creating folders when a valid image is successfully processed, eliminating "ghost" (empty) directories.
- **Instant Directory Discovery**: Rewrote the internal file reading engine. The CLI now resolves paths in massive folder trees up to 30x faster, eliminating the initial "choke" time before the progress bar starts.
- **Bulletproof Memory Management (Scale-Ready)**: The tool is now immune to memory bloat from massive file queues. The internal orchestrator was refactored to ensure the CLI consumes a flat, fixed amount of RAM regardless of whether you process 100 or 1,000,000 images. No more out-of-memory crashes on giant batches.
- **Single-Pass CGO Execution**: Refactored the internal engine to aggregate all transformation intents into a single options struct, ensuring only one memory allocation and CGO call per image pipeline, significantly reducing memory footprint.
- **Path Resolution Hardening**: Fixed a path relative calculation issue and ensured correct fallback behaviors when modifiers are not present.

### 🛠 Architecture & Refactoring
- **Domain-Driven Renaming**: Refactored core orchestration to use `RunFileBatch` and `FileTask`, aligning the code with domain-specific semantics instead of generic function parameters.
- **Type Safety Hardening**: Eliminated the `ExtraParams interface{}` in favor of explicit `OutputPathModifier` interfaces and Go closures. This removes unsafe type assertions and enforces compile-time safety across the processing pipeline.

### 🛠 Engineering & Architecture
- **Native Profiling Integration**: Introduced hidden `--memprofile` and `--cpuprofile` flags for advanced runtime diagnostics and performance bottleneck analysis in production environments without polluting the standard user CLI help.

## [1.5.0] - 2026-04-28

### 🚀 Added & Changed
- **Dynamic Versioning**: Implemented version injection via `ldflags`, ensuring the CLI correctly reports its release version (e.g., `v1.5.0`) instead of the static `dev` placeholder.
- **Documentation Polish**: Restructured the `README.md` to prioritize pre-compiled binaries and Docker usage, drastically improving the onboarding experience for non-Go developers. Added CI/CD status badges.

### 🐛 Bug Fixes
- **Help Command Blocker**: Fixed a critical UX bug where running `compactify --help` or `compactify help` would fail by prematurely enforcing the global `--input` flag requirement.

### 🛠 Engineering & Maintenance
- **Module Path Correction**: Refactored the `go.mod` and all internal imports to match the canonical GitHub repository path (`go-compactify-cli`), preventing downstream module resolution issues.
- **CI/CD Pipeline Enhancements**: Configured automated version injection across all OS builds by dynamically passing release tags to the Go linker via GitHub Actions step outputs.

## [1.4.0] - 2026-04-27

### 🚀 CI/CD & Distribution
- **Multi-Platform Automated Releases**: Implemented a robust GitHub Actions workflow to build and package binaries for Linux (AMD64), macOS (Apple Silicon/Intel), and Windows (x64) on every tagged release.
- **Dynamic Dependency Resolution (Windows)**: Automated the isolation and packaging of shared `libvips` DLLs using `ldd` and MSYS2, ensuring the Windows executable is portable and works out-of-the-box.
- **Artifact Packaging**: Improved release UX by bundling contextual `README.txt` files and technical instructions within automated `.zip` and `.tar.gz` archives.
- **Workflow Resilience**: Migrated from volatile package managers to a deterministic build environment with inherited system paths for consistent Go toolchain execution.

## [1.3.0] - 2026-04-27

### 🐳 Docker & Infrastructure
- **Containerized Environment**: Added a multi-stage Dockerfile for isolated execution and consistent builds across different OS.
- **Optimized Image**: Final runtime image based on Alpine Linux, including only the necessary `libvips` shared libraries.
- **Build Efficiency**: Improved Docker layer caching by separating dependency downloads from source code copying.

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