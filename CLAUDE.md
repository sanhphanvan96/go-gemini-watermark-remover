# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

- **Build**: `make build` or `go build -o gemini-wm-remove`
- **Install**: `make install` (installs binary to `~/.local/bin`)
- **Clean**: `make clean`
- **Run**: `./gemini-wm-remove -i <input_path> -o <output_path>`
  - Flags: `-i` (input), `-o` (output), `-w` (workers), `-v` (verbose)
- **Test**: `go test ./...`

## Architecture

- **Entry Point**: `main.go` handles CLI argument parsing using `flag`, file discovery, and manages a concurrent worker pool for processing images.
- **Core Logic**: `core/` package implements the watermark removal algorithm.
  - **Algorithm**: Reverse alpha blending (`engine.go`). Calculates alpha map from reference background images and reverses the blending operation to restore original pixels.
  - **Watermark Detection**: Automatically detects watermark size (48px or 96px) based on image dimensions (>1024x1024 uses 96px).
  - **Assets**: Reference watermark images (`bg_48.png`, `bg_96.png`) are embedded into the binary via `core/assets.go` using `//go:embed`.
