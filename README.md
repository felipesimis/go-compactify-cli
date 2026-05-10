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
- 🚀 **High Performance**: Uses `libvips` for low memory footprint and extreme speed, featuring instant directory traversal.
- 🧠 **Hardware-Aware**: Intelligent concurrency management that automatically scales to your host's CPU cores.
- 📁 **Deep Recursive Processing**: Seamlessly processes nested directory trees, mirroring the exact folder hierarchy.
- 🛡️ **Safety First**: Built-in `DryRun` mode to simulate filesystem operations before committing changes.
- ⚙️ **Multi-layer Configuration**: Precedence order: Flags > Env Vars > Config File > Defaults.
- 🛠️ **Versatile Processing**: Format conversion, resizing, cropping, grayscale, flipping, and lossless compression.
- 📊 **Detailed Analytics**: Execution summary with a side-by-side "Impact Dashboard".

---

## 🚀 Common Use Cases (How-To Guide)

This section provides direct answers to common developer questions. Each example is copy-paste ready.

### How do I process a massive library recursively while maintaining folder structure?
Use the `--recursive` (or `-r`) flag to mirror the input hierarchy in the output destination:
```bash
./compactify convert --format webp -i ./massive_library -o ./optimized_library --recursive
```

### How do I batch resize images with custom concurrency for high-end hardware?
Use the `--concurrency` flag to override the automatic CPU detection:
```bash
./compactify resize -w 800 -H 600 -i ./images --concurrency 16
```

### How can I preview changes without actually modifying any files?
Use the `--dry-run` flag to simulate the operation and verify output paths:
```bash
./compactify convert --format webp -i ./assets --dry-run
```

### How do I optimize images and strip sensitive EXIF/GPS metadata for privacy?
Add the `--strip-metadata` flag to remove camera and location info while preserving color profiles:
```bash
./compactify convert --format jpeg -i ./input -o ./output --strip-metadata
```

### How do I balance file size and visual quality?
Use the --quality (or -q) flag on any transformation command (1-100). Default is 75:
```bash
./compactify crop -w 800 -H 600 -q 85 -i ./assets
```

### How do I optimize images without losing any quality?
Use the `lossless` command to apply format-specific lossless optimizations:
```bash
./compactify lossless -i ./photos -o ./optimized_photos
```

### How do I specify a custom configuration file?
Use the `init` command or the `--config` flag to point to your YAML configuration (see `config.yaml.example` for reference):
```bash
./compactify init
```

### How do I set global settings using environment variables?
Prefix settings with `COMPACTIFY_`. For example, to set concurrency to 10:
```bash
export COMPACTIFY_CONCURRENCY=10
./compactify lossless
```

---

## ⚙️ Configuration Hierarchy
Compactify follows a strict precedence order (from highest to lowest):
1. **Command Line Flags** (e.g., `--concurrency 10`)
2. **Environment Variables** (prefixed with `COMPACTIFY_`)
3. **Configuration File** (`config.yaml`)
4. **Hardware Defaults** (automatically calculated)

### Environment Variables Mapping
| Environment Variable         | Flag Equivalent         |
| :---                         | :---                    |
| `COMPACTIFY_CONCURRENCY`     | `-c, --concurrency`     |
| `COMPACTIFY_INPUT`           | `-i, --input`           |
| `COMPACTIFY_OUTPUT`          | `-o, --output`          |
| `COMPACTIFY_QUALITY`         | `-q, --quality`         |
| `COMPACTIFY_RECURSIVE`       | `-r, --recursive`       |
| `COMPACTIFY_DRY_RUN`         | `--dry-run`             |
| `COMPACTIFY_STRIP_METADATA`  | `--strip-metadata`      |
| `COMPACTIFY_CONFIG`          | `--config`              |

---

## 🏗 Architecture & Engineering Decisions

#### 🧩 Intent-Based Architecture
By using the **Functional Options Pattern**, the CLI layer remains "intent-based". The underlying `bimg` engine translates these into a single `libvips` operation, ensuring **single-pass CGO execution** and reducing memory allocations.

### 🌊 Concurrency & Memory
Uses a **Semaphore Pattern** and a **Collector Goroutine**. The memory footprint remains statically flat ($O(C)$) regardless of the number of images, as the orchestration balances discovery and execution in parallel.

### 🛡️ Segregated Filesystem Abstraction
Instead of raw OS calls, the project uses a specialized interface for workers. This enables:
- **Dry-Run Mode**: Safely simulate operations without touching the disk.
- **Lazy Creation**: Directories are created "just-in-time", only when a processed file is ready to be saved, preventing empty "ghost" folders.
- **Infinite Loop Prevention**: Intelligent path detection that prevents the engine from recursively processing its own output when the destination resides within the source tree.

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
1. Download the latest release from the [Releases Page](https://github.com/felipesimis/go-compactify-cli/releases).
2. Extract and run. (*Windows: Keep DLLs in the same folder*).

### 🐳 2. Running with Docker
```bash
# Build the image locally
docker build -t compactify-cli .

# Execute via Docker (mapping your current directory)
docker run --rm -v "$(pwd):/workspace" compactify-cli lossless -i /workspace/images
```
> [!IMPORTANT]
> Paths must be relative to the `/workspace` directory inside the container.

### 🛠 3. Building from Source
Requires [Go](https://golang.org/doc/install) 1.26+ and [libvips](https://www.libvips.org/) headers.
- **macOS**: `brew install vips`
- **Linux**: `sudo apt install libvips-dev`

```bash
git clone https://github.com/felipesimis/go-compactify-cli.git
cd go-compactify-cli
go build -ldflags="-w -s -X 'github.com/felipesimis/go-compactify-cli/cmd.Version=$(git describe --tags --abbrev=0)'" -trimpath -o compactify .
```

---

## 🧪 Testing Standards
- **Unit & Integration**: High coverage via Dependency Injection and `testify/suite`.
- **Functional Mocking**: `FakeImageProcessor` stubs for fast orchestration validation.
- **E2E Tests**: Binary validation against real images for `CGO/libvips` stability.

**Commands:**
- All tests: `make test`
- Coverage: `make coverage`
- E2E tests: `make test-e2e`

---

## 🛠 Built With
- [Go](https://golang.org/) | [bimg](https://github.com/h2non/bimg) | [Cobra](https://github.com/spf13/cobra) | [Testify](https://github.com/stretchr/testify) | [Lipgloss](https://github.com/charmbracelet/lipgloss)

---

## 📄 License
Distributed under the MIT License. See `LICENSE` for more information.

---

*Developed by [Felipe Simis](https://github.com/felipesimis)*
