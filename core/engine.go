package core

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"sync"
)

// Constants
const (
	AlphaThreshold = 0.002
	MaxAlpha       = 0.99
	LogoValue      = 255.0
)

var (
	alphaMap48   []float64
	alphaMap96   []float64
	alphaMapOnce sync.Once
)

// WatermarkConfig holds the configuration for a given image size
type WatermarkConfig struct {
	LogoSize     int
	MarginRight  int
	MarginBottom int
}

// DetectWatermarkConfig determines the configuration based on image dimensions
func DetectWatermarkConfig(width, height int) WatermarkConfig {
	if width > 1024 && height > 1024 {
		return WatermarkConfig{
			LogoSize:     96,
			MarginRight:  64,
			MarginBottom: 64,
		}
	}
	return WatermarkConfig{
		LogoSize:     48,
		MarginRight:  32,
		MarginBottom: 32,
	}
}

// CalculateAlphaMap generates the alpha map from a reference background image
func CalculateAlphaMap(bg image.Image) []float64 {
	bounds := bg.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	alphaMap := make([]float64, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Convert to NRGBA to get non-premultiplied values
			c := color.NRGBAModel.Convert(bg.At(x, y)).(color.NRGBA)

			// Use max channel value
			maxChannel := math.Max(float64(c.R), math.Max(float64(c.G), float64(c.B)))
			alphaMap[y*width+x] = maxChannel / 255.0
		}
	}
	return alphaMap
}

func getAlphaMap(size int) []float64 {
	alphaMapOnce.Do(func() {
		alphaMap48 = CalculateAlphaMap(Bg48)
		alphaMap96 = CalculateAlphaMap(Bg96)
	})
	if size == 48 {
		return alphaMap48
	}
	return alphaMap96
}

// RemoveWatermark applies the reverse alpha blending to remove the watermark
// It returns a new image with the watermark removed.
func RemoveWatermark(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Create a mutable copy of the image (NRGBA for non-premultiplied manipulation)
	// We use NRGBA to ensure we're modifying raw pixel values
	out := image.NewNRGBA(bounds)
	draw.Draw(out, bounds, img, bounds.Min, draw.Src)

	config := DetectWatermarkConfig(width, height)

	// Determine which alpha map to use
	alphaMap := getAlphaMap(config.LogoSize)

	// Calculate position
	wmX := width - config.MarginRight - config.LogoSize
	wmY := height - config.MarginBottom - config.LogoSize

	// Process pixels in the watermark area
	for row := 0; row < config.LogoSize; row++ {
		for col := 0; col < config.LogoSize; col++ {
			// Coordinates in the image
			x := wmX + col
			y := wmY + row

			// Check bounds
			if x < 0 || x >= width || y < 0 || y >= height {
				continue
			}

			// Get alpha value from map
			alphaIdx := row*config.LogoSize + col
			if alphaIdx >= len(alphaMap) {
				continue
			}
			alpha := alphaMap[alphaIdx]

			if alpha < AlphaThreshold {
				continue
			}

			if alpha > MaxAlpha {
				alpha = MaxAlpha
			}

			oneMinusAlpha := 1.0 - alpha

			// Get current pixel
			// Since we are iterating over 'out' which is NRGBA, we can access Pix directly or use At
			offset := out.PixOffset(x, y)

			// Process channels
			processChannel := func(val uint8) uint8 {
				v := float64(val)
				original := (v - alpha*LogoValue) / oneMinusAlpha
				if original < 0 {
					original = 0
				}
				if original > 255 {
					original = 255
				}
				return uint8(math.Round(original))
			}

			out.Pix[offset+0] = processChannel(out.Pix[offset+0]) // R
			out.Pix[offset+1] = processChannel(out.Pix[offset+1]) // G
			out.Pix[offset+2] = processChannel(out.Pix[offset+2]) // B
			// Alpha (offset+3) remains unchanged
		}
	}

	return out, nil
}
