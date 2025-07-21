package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func generateFileChunks(filepath string, chunkSize int) ([]string, error) {
	chunks := []string{}
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file '%s': %w", filepath, err)
	}
	defer file.Close()

	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading file '%s': %w", filepath, err)
		}
		if n == 0 {
			break
		}

		chunk := buffer[:n]
		hash := sha256.Sum256(chunk)
		chunks = append(chunks, hex.EncodeToString(hash[:]))
	}

	return chunks, nil
}

func AnalyzeFiles() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: go run main.go FILE [FILE...] ")
		fmt.Println("Analyzes files and generates chunks to detect changes.")
		fmt.Println("Paths to files to analyze (maximum 10).")
		os.Exit(1)
	}

	if len(args) > 10 {
		fmt.Println("Error: More than 10 files provided. Please limit to 10 files.")
		os.Exit(1)
	}

	fmt.Println("Analyzing the following files:")
	for _, filepath := range args {
		fmt.Printf("- %s\n", filepath)
	}

	fmt.Println("\nChunk analysis results:")
	for _, filepath := range args {
		_, err := os.Stat(filepath)
		if os.IsNotExist(err) {
			fmt.Printf("\nWarning: File '%s' does not exist and will be skipped.\n", filepath)
			continue
		} else if err != nil {
			fmt.Printf("\nError accessing '%s': %v. It will be skipped.\n", filepath, err)
			continue
		}

		fmt.Printf("\n--- %s ---\n", filepath)
		chunks, err := generateFileChunks(filepath, 4096)
		if err != nil {
			fmt.Printf("  Could not generate chunks for this file due to an error: %v\n", err)
		} else {
			if len(chunks) > 0 {
				fmt.Printf("Number of chunks: %d\n", len(chunks))
				for i, chunkHash := range chunks {
					fmt.Printf("  Chunk %d: %s\n", i+1, chunkHash)
				}
			} else {
				fmt.Println("  Could not generate chunks for this file (possibly empty or error).")
			}
		}
	}
}
