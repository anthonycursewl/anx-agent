package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/anthonycursewl/anx-agent/internal/ai"
	"github.com/anthonycursewl/anx-agent/internal/config"
	"github.com/anthonycursewl/anx-agent/internal/utils"
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

	reader := bufio.NewReader(os.Stdin)

	utils.AgentName.Println("\n  ╔══════════════════════════════╗")
	utils.AgentName.Println("  ║        ANX AGENT v1.0       ║")
	utils.AgentName.Println("  ╚══════════════════════════════╝")
	fmt.Println("  Escribe 'ayuda' para ver los comandos disponibles")
	fmt.Println("  Escribe 'salir' para terminar")

	systemPrompt := `Eres ANX, un asistente de programación especializado en Go. 
Tus respuestas deben ser concisas y técnicas. Si el usuario te pide leer un archivo, 
muestra su contenido formateado de manera legible.`

	_, err = aiClient.CallModel("SYSTEM: " + systemPrompt)
	if err != nil {
		utils.ErrorMsg.Printf("Error al configurar el prompt del sistema: %v\n", err)
	}

	for {
		utils.UserInput.Print(`
┌─[🦋 ANX]
└─▪ `)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		lowerInput := strings.ToLower(input)
		switch {
		case input == "":
			continue
		case lowerInput == "salir" || lowerInput == "exit":
			utils.AgentName.Println("\n¡Hasta luego! 👋")
			return
		case lowerInput == "ayuda" || lowerInput == "help":
			showHelp()
			continue
		case strings.HasPrefix(lowerInput, "leer "):
			filePath := strings.TrimSpace(input[5:])
			if filePath == "" {
				utils.ErrorMsg.Println("Por favor especifica la ruta del archivo a leer")
				continue
			}
			content, err := utils.ReadFile(filePath)
			if err != nil {
				utils.ErrorMsg.Printf("Error leyendo archivo: %v\n", err)
				continue
			}
			utils.AgentResponse.Println("\nContenido del archivo:")
			utils.FileContent.Println("\n" + content + "\n")
			continue
		}

		stopLoading := utils.StartLoading("Procesando tu solicitud...")
		response, err := aiClient.CallModel(input)
		stopLoading()

		if err != nil {
			utils.ErrorMsg.Printf("Error: %v\n", err)
			continue
		}

		utils.AgentResponse.Println("\n" + formatResponse(response))
	}
}

func showHelp() {
	helpText := `
Comandos disponibles:
  leer <ruta>    - Muestra el contenido de un archivo
  ayuda          - Muestra esta ayuda
  salir          - Termina la sesión

Puedes hacer preguntas directamente al agente o pedirle que analice código.`

	utils.AgentResponse.Println(helpText)
}

func formatResponse(response string) string {
	return "  " + strings.ReplaceAll(response, "\n", "\n  ")
}
