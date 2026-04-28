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
- 🧠 **Hardware-Aware**: Intelligent concurrency management using a semaphore pattern to optimize CPU utilization without exhausting resources.
- 🛡️ **Safety First**: Built-in `DryRun` mode allows users to simulate filesystem operations before committing changes, preventing accidental data loss.
- 🛠️ **Versatile Processing**:
    - Format conversion (JPEG, PNG, WebP, etc.)
    - Intelligent resizing and cropping.
    - Grayscale, flipping, and color palette optimization.
    - Lossless compression.
- 📊 **Detailed Analytics**: Comprehensive execution summary with color-coded statistics via Lipgloss.

---

## 🏗 Architecture & Engineering Decisions

### 🧩 Decoupled Architecture
The core logic is strictly isolated from external dependencies. By using the **Dependency Inversion Principle**, the `internal/filesystem` package interacts with an `OSOperations` interface. This allows for 100% unit test coverage by isolating side effects through sophisticated mocking.

### 🌊 Concurrency Model
To handle thousands of images efficiently, Compactify uses a **Semaphore Pattern** (`chan struct{}`). This prevents goroutine explosion and ensures the tool respects the host machine's hardware limits, maintaining stability under heavy load.

### 🛡️ The Dry-Run Pattern
Implementing the `FileReaderWriter` interface, the tool supports a non-destructive simulation mode. This is critical for CLI tools that perform destructive operations (like overwriting images), providing a "safety net" for the user.

---

## 📂 Project Structure

```text
.
├── cmd/                # CLI command implementations (Cobra)
├── internal/
│   ├── filesystem/     # Core filesystem abstraction & DryRun logic
│   ├── image/          # bimg/libvips wrappers
│   ├── processing/     # Orchestration of the image processing pipeline
│   └── utils/          # Validation, path handling, and statistics
├── pkg/                # Publicly exportable packages
└── main.go             # Application entry point
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

# Execute (Linux / macOS / WSL)
docker run --rm -v "$(pwd):/workspace" compactify-cli lossless -i /workspace/images

# Execute (Windows PowerShell)
docker run --rm -v "${PWD}:/workspace" compactify-cli lossless -i /workspace/images
```
> [!IMPORTANT]
> Path Mapping: When using Docker, all input (-i) and output (-o) paths must be relative to the /workspace directory inside the container.

### 🛠 3. Building from Source (Developers)
Requires [Go](https://golang.org/doc/install) 1.21+ and [libvips](https://www.libvips.org/) headers.

### Installation (Native)

1. Clone the repository:
   ```bash
   git clone https://github.com/felipesimis/go-compactify-cli.git
   cd go-compactify-cli
   ```

2. Build the binary (injecting version):
   ```bash
   # In the command below, replace 'v1.5.0' with the current version tag
   go build -ldflags="-w -s -X 'github.com/felipesimis/go-compactify-cli/cmd.Version=$(git describe --tags --abbrev=0)'" -trimpath -o compactify .
   ```

### Usage

Run the help command to see all available options:
```bash
./compactify --help
```

Example: Resize all images in a folder with a dry run:
```bash
./compactify resize -w 800 -H 600 --input ./images --dry-run
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
