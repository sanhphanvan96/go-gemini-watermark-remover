package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"gemini-wm-remover/core"
)

var (
	inputPath  string
	outputPath string
	workers    int
	verbose    bool
)

func init() {
	flag.StringVar(&inputPath, "input", "", "Input file or directory")
	flag.StringVar(&inputPath, "i", "", "Input file or directory (shorthand)")

	flag.StringVar(&outputPath, "output", "output", "Output directory")
	flag.StringVar(&outputPath, "o", "output", "Output directory (shorthand)")

	flag.IntVar(&workers, "workers", runtime.NumCPU(), "Number of concurrent workers")
	flag.IntVar(&workers, "w", runtime.NumCPU(), "Number of concurrent workers (shorthand)")

	flag.BoolVar(&verbose, "verbose", false, "Verbose logging")
	flag.BoolVar(&verbose, "v", false, "Verbose logging (shorthand)")

	// Register formats
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func main() {
	flag.Parse()

	if inputPath == "" {
		fmt.Println("Usage: gemini-wm-remove -i <path> [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create output directory
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Collect files
	files, err := collectFiles(inputPath)
	if err != nil {
		fmt.Printf("Error collecting files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No supported images found (png, jpg, jpeg, webp).")
		return
	}

	fmt.Printf("Found %d images. Processing with %d workers...\n", len(files), workers)

	// Worker pool
	jobs := make(chan string, len(files))
	results := make(chan string, len(files))
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// Wait for completion in background
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results
	successCount := 0
	start := time.Now()
	for res := range results {
		if res != "" {
			if verbose {
				fmt.Println(res)
			}
		} else {
			// Error case, handled in worker
		}
		successCount++
	}

	duration := time.Since(start)
	fmt.Printf("Processed %d images in %.2fs\n", successCount, duration.Seconds())
	fmt.Printf("Output saved to: %s\n", outputPath)
}

func collectFiles(path string) ([]string, error) {
	var files []string
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return []string{path}, nil
	}

	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(p))
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".webp" {
				files = append(files, p)
			}
		}
		return nil
	})
	return files, err
}

func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range jobs {
		err := processFile(file)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", filepath.Base(file), err)
			results <- ""
		} else {
			results <- fmt.Sprintf("Processed: %s", filepath.Base(file))
		}
	}
}

func processFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	cleaned, err := core.RemoveWatermark(img)
	if err != nil {
		return fmt.Errorf("processing error: %w", err)
	}

	outName := filepath.Join(outputPath, filepath.Base(path))
	outF, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer outF.Close()

	if format == "png" {
		return png.Encode(outF, cleaned)
	} else {
		return jpeg.Encode(outF, cleaned, &jpeg.Options{Quality: 95})
	}
}
