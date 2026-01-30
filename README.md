# Gemini Watermark Remover (Go)

A CLI tool to remove watermarks from images using the reverse alpha blending algorithm.

## Installation

You can install the tool globally using `make`:

```bash
make install
```

This will install `gemini-wm-remove` to `~/.local/bin`. Ensure this directory is in your `PATH`.

## Usage

You can run the tool from anywhere:

```bash
gemini-wm-remove -i <input> -o <output>
```

### Flags

- `-i, -input`: Input file or directory (required).
- `-o, -output`: Output directory (default: `output`).
- `-w, -workers`: Number of concurrent workers (default: NumCPU).
- `-v, -verbose`: Enable verbose logging.

### Examples

Process a directory:
```bash
gemini-wm-remove -i ./photos -o ./cleaned
```

Process a single file:
```bash
gemini-wm-remove -i image.jpg -o cleaned/
```

Process current directory:
```bash
gemini-wm-remove -i .
```

## Development

- `make build`: Build the binary locally.
- `make install`: Build and install globally.
- `make clean`: Remove built binary.
