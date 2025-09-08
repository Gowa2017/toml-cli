# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a CLI tool for editing and querying TOML files, written in Go. It provides multiple commands for managing TOML configuration files:
- `get`: Extract values from TOML files using dot notation queries
- `set`: Update values in TOML files with optional output file specification
- `merge`: Merge two TOML files with recursive object merging
- Additional commands for listing, deleting, clearing, dumping, importing, and scanning TOML files

## Development Commands

### Build and Run
```bash
go build -o cm
./cm get -c ./sample/get-set/app.toml app.name
./cm set -c ./sample/get-set/app.toml 192.168.11.11 title test [-o output_file]
./cm merge ./sample/merge/base.toml ./sample/merge/override.toml [-o output_file]
```

### Testing
```bash
go test -v ./...           # Run all tests
go test -v ./toml          # Run specific package tests
```

### Standard Go Commands
```bash
go run main.go <args>      # Run without building
go build                  # Build binary
go mod tidy               # Clean up dependencies
go fmt ./...              # Format code
```

## Architecture

### Core Components

1. **Command Structure** (`cmd/`):
   - `root.go`: Main cobra command setup and execution
   - `get.go`: Implementation of get command for querying TOML values
   - `set.go`: Implementation of set command for updating TOML values
   - `merge.go`: Implementation of merge command for combining TOML files
   - Additional command files for extended functionality (list, delete, clear, dump, import, scan, etc.)

2. **TOML Processing** (`toml/`):
   - `toml.go`: Core TOML wrapper around pelletier/go-toml library with additional functionality
   - `tomlw.go`: File I/O operations for reading/writing TOML files

3. **Entry Point**:
   - `main.go`: Simple entry point that calls cmd.Execute()

### Key Dependencies
- `github.com/pelletier/go-toml`: Core TOML parsing library
- `github.com/spf13/cobra`: CLI framework
- `github.com/stretchr/testify`: Testing framework
- `github.com/fatih/color`: Terminal color output

### Data Flow
1. Commands parse arguments using cobra
2. Create `Toml` struct instance from file path
3. For `get`: Query value using dot notation and print result
4. For `set`: Update value, optionally specify output file, then write changes
5. For `merge`: Load two TOML files, merge recursively, then write combined result

### Value Type Handling
The `set` command automatically detects and converts input types:
- Boolean values (true/false)
- Integer values (int64)
- Float values (float64)
- TOML date/time formats (LocalDate, LocalDateTime, LocalTime)
- Strings (fallback)

## Testing Strategy

The codebase includes comprehensive tests in the `toml/` package covering:
- Core TOML operations (get/set/merge)
- Error cases and edge conditions
- Value type handling
- File I/O operations