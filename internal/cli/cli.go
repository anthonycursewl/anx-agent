package cli

import (
	"fmt"
	"strings"

	"github.com/anthonycursewl/anx-agent/internal/ai"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CLI struct {
	aiClient *ai.Client
}

func NewCLI(aiClient *ai.Client) *CLI {
	return &CLI{
		aiClient: aiClient,
	}
}

func (c *CLI) Start() error {
	p := tea.NewProgram(initialModel(c.aiClient))
	return p.Start()
}

var (
	agentNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("13")).
			Background(lipgloss.Color("6"))

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)
)

type model struct {
	aiClient  *ai.Client
	textInput textinput.Model
	messages  []string
	err       error
	width     int
	height    int
}

func initialModel(aiClient *ai.Client) model {
	ti := textinput.New()
	ti.Placeholder = "Escribe tu mensaje aquí..."
	ti.Prompt = "❯ "
	ti.PromptStyle = promptStyle
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

	ti.Focus()
	ti.CharLimit = 1000
	ti.Width = 50

	return model{
		aiClient:  aiClient,
		textInput: ti,
		messages:  []string{"Bienvenido al Agente de Análisis de Código. Escribe 'ayuda' para ver los comandos disponibles."},
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textInput.Width = msg.Width - 10
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			input := strings.TrimSpace(m.textInput.Value())
			if input == "" {
				return m, nil
			}

			m.messages = append(m.messages, "\n  Tú: "+input)

			if strings.ToLower(input) == "salir" || strings.ToLower(input) == "exit" {
				m.messages = append(m.messages, "\n  Agente: ¡Hasta luego!")
				return m, tea.Quit
			}

			m.textInput.Reset()

			return m, tea.Batch(
				cmd,
				func() tea.Msg {
					response, err := m.aiClient.CallModel(input)
					if err != nil {
						return err
					}
					return response
				},
			)

		case tea.KeyCtrlC, tea.KeyEsc:
			m.messages = append(m.messages, "\n  Agente: ¡Hasta luego!")
			return m, tea.Quit
		}

	case string:
		m.messages = append(m.messages, "\n  Agente: "+msg)
		return m, nil

	case error:
		errMsg := fmt.Sprintf("\n  Error: %v", msg)
		m.messages = append(m.messages, errorStyle.Render(errMsg))
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	header := agentNameStyle.Padding(0, 1).Render(" ANX AGENT ")
	header = lipgloss.PlaceHorizontal(m.width, lipgloss.Left, header)

	var messages string
	for _, msg := range m.messages {
		messages += fmt.Sprintf("\n%s", msg)
	}

	input := fmt.Sprintf(
		"\n%s\n%s",
		promptStyle.Render("Escribe tu mensaje (presiona Ctrl+C para salir):"),
		m.textInput.View(),
	)

	return fmt.Sprintf(
		"%s\n%s\n%s",
		header,
		messages,
		input,
	)
}

func RunCLI(aiClient *ai.Client) error {
	cli := NewCLI(aiClient)
	return cli.Start()
}
