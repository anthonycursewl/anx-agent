// cmd/agentcli/main.go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/anthonycursewl/anx-agent/internal/ai"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file (e.g., config.yaml)")
	inputPath := flag.String("input", ".", "Path to the file or directory to analyze")

	flag.Parse()

	if *inputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Input path is required.")
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Verifica que la API Key de Gemini esté disponible
	if cfg.GEMINI_API_KEY == "" {
		fmt.Fprintln(os.Stderr, "Error: Gemini API Key is not set. Please check configuration.")
		os.Exit(1)
	}

	aiClient, err := ai.NewClient(cfg.GEMINI_API_KEY)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing AI client: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := aiClient.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing AI client: %v\n", err)
		}
	}()

	promt := "wola como estás"
	err, response := aiClient.CallModel(promt)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling AI model: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", response)
	fmt.Println("Analysis completed successfully.")
}
