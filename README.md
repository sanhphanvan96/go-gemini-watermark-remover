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

## Technical Details

This tool uses **Reverse Alpha Blending** to restore original pixel values.

### How it works

1.  **Watermark Detection**:
    - If image is > 1024x1024, it assumes a 96px watermark with 64px margin.
    - Otherwise, it assumes a 48px watermark with 32px margin.
    - Position is always fixed to the bottom-right corner.

2.  **Alpha Map**:
    - The tool calculates an "alpha map" from a reference watermark image on a black background.
    - Alpha value `α` for each pixel is calculated as `max(R, G, B) / 255.0`.

3.  **Restoration**:
    - The tool applies the reverse blending formula to recover the original pixel color `C_orig`:
        ```
        C_orig = (C_mixed - α * C_watermark) / (1 - α)
        ```
    - Where `C_mixed` is the current pixel color and `C_watermark` is white (255).
    - Pixels with very low alpha (< 0.002) are ignored to preserve noise/texture.
    - Alpha is clamped to 0.99 to prevent division by zero.
