# 📸 Compactify CLI

[![Go Report Card](https://goreportcard.com/badge/github.com/felipe-simis/go-compactify-cli)](https://goreportcard.com/report/github.com/felipe-simis/go-compactify-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://golang.org)
![Docker](https://img.shields.io/badge/Docker-supported-blue?logo=docker)

**Compactify CLI** is a high-performance, hardware-aware image optimization tool built in Go. It leverages `libvips` (via `bimg`) to provide lightning-fast image processing, including resizing, cropping, format conversion, and lossless compression.

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

### Prerequisites

- [Go](https://golang.org/doc/install) 1.21+
- [libvips](https://www.libvips.org/) (Required for image processing)
- [Docker](https://www.docker.com/products/docker-desktop/) (Optional, for containerized execution)

### Installation (Native)

1. Clone the repository:
   ```bash
   git clone https://github.com/felipe-simis/go-compactify-cli.git
   cd go-compactify-cli
   ```

2. Build the binary:
   ```bash
   go build -ldflags="-w -s" -trimpath -o compactify . 
   ```

### 🐳 Running with Docker (Alternative)
If you prefer not to install libvips or Go locally, use the pre-configured Docker environment. This ensures a consistent environment across all operating systems.

1. Build the image:
   ```bash
   docker build -t compactify-cli .
   ```
2. Execute via Docker:
You must map your local folder to the container's /workspace using a volume.
   ## Linux / macOS / WSL
   ```bash
   docker run --rm -v "$(pwd):/workspace" compactify-cli lossless -i /workspace/your-folder
   ```
   ## Windows (PowerShell)
   ```bash
   docker run --rm -v "${PWD}:/workspace" compactify-cli lossless -i /workspace/your-folder
   ```
   (> [!IMPORTANT])
Path Mapping: When using Docker, all input (-i) and output (-o) paths must be relative to the /workspace directory inside the container.

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
- [libvips](https://www.libvips.org/) - Fast image processing library.
- [Cobra](https://github.com/spf13/cobra) - CLI framework.
- [Testify](https://github.com/stretchr/testify) - Testing toolkit.
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling library.

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---
*Developed by [Felipe Simis](https://github.com/felipe-simis)*
