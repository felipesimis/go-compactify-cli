# 📸 Compactify CLI

[![Go Report Card](https://goreportcard.com/badge/github.com/felipesimis/go-compactify-cli)](https://goreportcard.com/report/github.com/felipesimis/go-compactify-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://golang.org)
![Docker](https://img.shields.io/badge/Docker-supported-blue?logo=docker)
![GitHub Release](https://img.shields.io/github/v/release/felipesimis/go-compactify-cli)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/felipesimis/go-compactify-cli/release.yaml)

**Compactify CLI** is a high-performance image optimization tool. It uses the `bimg` library to leverage the extreme speed and low memory footprint of `libvips`, providing hardware-accelerated processing.

Designed with **software engineering excellence** in mind, the project follows strict architectural patterns to ensure testability, safety, and scalability.

---

## ✨ Key Features

- 🚀 **High Performance**: Uses `libvips` for low memory footprint and extreme speed.
- 🧠 **Hardware-Aware**: Intelligent concurrency management using a semaphore pattern to optimize CPU utilization.
- 🛡️ **Safety First**: Built-in `Dry Run` mode allows users to simulate filesystem operations before committing changes, preventing accidental data loss.
- ⚙️ **Multi-layer Configuration**: Support for Config Files, Environment Variables, and Flags with a strict precedence order.
- 🛠️ **Versatile Processing**:
    - Format conversion (JPEG, PNG, WebP, etc.)
    - Intelligent resizing and cropping.
    - Grayscale, flipping, and color palette optimization.
    - Lossless compression.
- 📊 **Detailed Analytics**: Execution summary with a side-by-side "Impact Dashboard" (Original vs. Processed).

---

## ⚙️ Configuration Hierarchy

Compactify follows a strict precedence order (from highest to lowest). This allows for flexible deployments in local, CI/CD, or Docker environments:

1. **Command Line Flags** (e.g., `--concurrency 10`)
2. **Environment Variables** (prefixed with `COMPACTIFY_`)
3. **Configuration File** (`config.yaml`)
4. **Hardware Defaults** (automatically calculated based on CPU cores)

### Environment Variables Mapping

| Environment Variable | Flag Equivalent |
| :--- | :--- |
| `COMPACTIFY_CONCURRENCY` | `-c, --concurrency` |
| `COMPACTIFY_INPUT` | `-i, --input` |
| `COMPACTIFY_OUTPUT` | `-o, --output` |
| `COMPACTIFY_DRY_RUN` | `--dry-run` |
| `COMPACTIFY_CONFIG` | `--config` |

---

## 🏗 Architecture & Engineering Decisions

### 🧩 Decoupled Architecture
The core logic is strictly isolated from external dependencies. By using the **Dependency Inversion Principle**, the internal packages interact with interfaces, allowing for nearly 100% test coverage of the command orchestration and configuration logic.

### 🌊 Concurrency Model
To handle thousands of images efficiently, Compactify uses a **Semaphore Pattern** (`chan struct{}`). This prevents goroutine explosion and ensures the tool respects the host machine's hardware limits.

### 🛡️ The Dry-Run Pattern
Implementing the `FileReaderWriter` interface, the tool supports a non-destructive simulation mode. This is critical for CLI tools that perform destructive operations, providing a "safety net" for the user.

---

## 📂 Project Structure

```text
.
├── cmd/                # CLI command implementations (Cobra)
├── internal/
│   ├── filesystem/     # Core filesystem abstraction & Dry Run logic
│   ├── image/          # bimg/libvips wrappers
│   ├── processing/     # Orchestration of the image processing pipeline
│   ├── templates/      # Configuration and UI templates
│   ├── ui/             # High-fidelity terminal UI components
│   ├── utils/          # Validation, path handling, and statistics
│   └── validation/     # Input validation logic
├── pkg/                # Publicly exportable packages
│   └── progress/       # Terminal progress bar implementation
└── main.go             # Application entrypoint
```
---

## 🚀 Getting Started

### 📥 1. Pre-compiled Binaries (Recommended)
The fastest way to use Compactify. No system dependencies (Go or libvips) are required.
1. Download the latest release for your OS from the [Releases Page](https://github.com/felipesimis/go-compactify-cli/releases).
2. Extract the archive and run the executable via terminal.
   * *Windows users: Keep the provided DLLs in the same folder as the executable.*

### 🐳 2. Running with Docker
Perfect for consistent environments without installing native dependencies. Requires [Docker](https://docs.docker.com/get-docker/) installed.

```bash
# Build the image locally
docker build -t compactify-cli .

# Execute via Docker (mapping your current directory)
docker run --rm -v "$(pwd):/workspace" compactify-cli lossless -i /workspace/images
```
> [!IMPORTANT]
> Path Mapping: When using Docker, all input (-i) and output (-o) paths must be relative to the /workspace directory inside the container.

### 🛠 3. Building from Source (Developers)
Building from source requires [Go](https://golang.org/doc/install) 1.21+ and [libvips](https://www.libvips.org/) headers installed in your system.

- **macOS**: `brew install vips`
- **Linux**: `sudo apt install libvips-dev`
- **Windows**: Follow the `vips` Windows installation guide.

#### Installation (Native)

Clone the repository:
   ```bash
   # Clone the repository
   git clone https://github.com/felipesimis/go-compactify-cli.git
   cd go-compactify-cli

   # Build with version injection
   go build -ldflags="-w -s -X 'github.com/felipesimis/go-compactify-cli/cmd.Version=$(git describe --tags --abbrev=0)'" -trimpath -o compactify .
   ```

#### Quick Start

```bash
# Initialize a default configuration file
./compactify init

# Batch resize all images with a high concurrency
./compactify resize -w 800 -H 600 -i ./images --concurrency 12

# Convert all images to WebP without actually touching the files (Preview)
./compactify convert --format webp -i ./assets --dry-run

# Run lossless optimization with concurrency set via environment variable
export COMPACTIFY_CONCURRENCY=10
./compactify lossless
```

---

## 🧪 Testing Standards

We aim for maximum reliability. Every feature is backed by a suite of unit tests using `testify/assert`.

**Run all tests with coverage:**
```bash
go test ./... -coverprofile="coverage.out"
```

**Run a specific test:**
```bash
go test -v ./internal/filesystem -run TestReadDir
```

---

## 🛠 Built With

- [Go](https://golang.org/) - The programming language.
- [bimg](https://github.com/h2non/bimg) - Go bindings for libvips.
- [Cobra](https://github.com/spf13/cobra) - CLI framework.
- [Testify](https://github.com/stretchr/testify) - Testing toolkit.
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling library.

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

*Developed by [Felipe Simis](https://github.com/felipesimis)*
