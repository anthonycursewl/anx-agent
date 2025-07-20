package main

import (
	"fmt"
	"os"

	"github.com/anthonycursewl/anx-agent/internal/ai"
	"github.com/anthonycursewl/anx-agent/internal/cli"
	"github.com/anthonycursewl/anx-agent/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
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
	cli.Start(aiClient)
}
