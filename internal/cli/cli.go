package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonycursewl/anx-agent/internal/ai"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type styles struct {
	app           lipgloss.Style
	header        lipgloss.Style
	userMsg       lipgloss.Style
	aiMsg         lipgloss.Style
	errorMsg      lipgloss.Style
	infoMsg       lipgloss.Style
	statusBar     lipgloss.Style
	statusText    lipgloss.Style
	statusSpinner lipgloss.Style
}

func defaultStyles() styles {
	return styles{
		app:           lipgloss.NewStyle().Margin(1, 2),
		header:        lipgloss.NewStyle().Foreground(lipgloss.Color("#76D7C4")).Bold(true).MarginBottom(1),
		userMsg:       lipgloss.NewStyle().Foreground(lipgloss.Color("#5DADE2")).MarginLeft(2),
		aiMsg:         lipgloss.NewStyle().Foreground(lipgloss.Color("#F7DC6F")).MarginLeft(2),
		errorMsg:      lipgloss.NewStyle().Foreground(lipgloss.Color("#E74C3C")).MarginLeft(2),
		infoMsg:       lipgloss.NewStyle().Foreground(lipgloss.Color("#AAB7B8")).MarginLeft(2),
		statusBar:     lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false, false).BorderForeground(lipgloss.Color("#566573")).Padding(0, 1),
		statusText:    lipgloss.NewStyle().Foreground(lipgloss.Color("#85929E")),
		statusSpinner: lipgloss.NewStyle().Foreground(lipgloss.Color("#FAD02E")),
	}
}

const (
	modeChat int = iota
	modeExplorer
	modeCreateFileInput
	modeAIFilenameInput
	modeAIPromptInput
)

type aiResponseMsg string
type aiFileContentMsg struct {
	fileName string
	content  string
}
type fileWrittenMsg string
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type item struct {
	path  string
	isDir bool
}

func (i item) Title() string {
	if i.isDir {
		return "📁 " + filepath.Base(i.path)
	}
	return "📄 " + filepath.Base(i.path)
}
func (i item) Description() string { return i.path }
func (i item) FilterValue() string { return i.path }

type model struct {
	aiClient         *ai.Client
	commands         map[string]Command
	list             list.Model
	textInput        textinput.Model
	messages         []string
	spinner          spinner.Model
	loading          bool
	width            int
	height           int
	mode             int
	currentPath      string
	styles           styles
	fileCreationName string
}

type Command struct {
	Name        string
	Description string
	Execute     func(m *model, args []string) tea.Cmd
}

func initialModel(aiClient *ai.Client) model {
	ti := textinput.New()
	ti.Placeholder = "Escribe un mensaje o comando ('help' para ayuda)..."
	ti.Focus()
	ti.CharLimit = 1024
	ti.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAD02E"))

	items := []list.Item{}
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("#C472DA")).Foreground(lipgloss.Color("#C472DA")).Padding(0, 0, 0, 1)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle.Copy()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Explorador de Archivos"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "volver")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "crear vacío")),
			key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "crear con IA")),
		}
	}

	m := model{
		aiClient:    aiClient,
		textInput:   ti,
		messages:    []string{"info:Bienvenido al Agente ANX. Escribe 'help' para ver los comandos disponibles."},
		spinner:     s,
		currentPath: ".",
		list:        l,
		mode:        modeChat,
		styles:      defaultStyles(),
	}

	m.registerCommands()
	return m
}

func (m *model) registerCommands() {
	m.commands = map[string]Command{
		"help": {
			Name: "help", Description: "Muestra este mensaje de ayuda",
			Execute: helpCommand,
		},
		"ls": {
			Name: "ls", Description: "Muestra el explorador de archivos",
			Execute: listCommand,
		},
		"exit": {
			Name: "exit", Description: "Sale de la aplicación",
			Execute: exitCommand,
		},
		"quit": {
			Name: "quit", Description: "Alias para 'exit'",
			Execute: exitCommand,
		},
	}
}

func helpCommand(m *model, args []string) tea.Cmd {
	helpText := "Comandos Disponibles:\n"
	for name, cmd := range m.commands {
		helpText += fmt.Sprintf("  %-15s %s", name, cmd.Description)
		m.messages = append(m.messages, "info:"+helpText)
		helpText = ""
	}
	m.messages = append(m.messages, "info:\nEn el explorador ('ls'):\n  'c' para crear archivo vacío\n  'a' para crear archivo con IA")
	return nil
}

func listCommand(m *model, args []string) tea.Cmd {
	m.mode = modeExplorer
	return m.listDirectory(m.currentPath)
}

func exitCommand(m *model, args []string) tea.Cmd {
	m.messages = append(m.messages, "info:¡Hasta luego!")
	return tea.Quit
}

func (m *model) listDirectory(path string) tea.Cmd {
	m.currentPath = path
	return func() tea.Msg {
		entries, err := os.ReadDir(path)
		if err != nil {
			return errMsg{err}
		}

		items := []list.Item{
			item{path: "..", isDir: true},
		}
		for _, entry := range entries {
			items = append(items, item{path: entry.Name(), isDir: entry.IsDir()})
		}
		return m.list.SetItems(items)
	}
}

func (m *model) createFileWithContent(filePath string, content string) tea.Cmd {
	m.loading = true
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			err := os.WriteFile(filePath, []byte(content), 0644)
			if err != nil {
				return errMsg{err}
			}
			return fileWrittenMsg(filePath)
		},
	)
}

func (m *model) createFile(filePath string) tea.Cmd {
	m.loading = true
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return func() tea.Msg {
			return errMsg{fmt.Errorf("el archivo '%s' ya existe", filePath)}
		}
	}

	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			f, err := os.Create(filePath)
			if err != nil {
				return errMsg{err}
			}
			_ = f.Close()
			return fileWrittenMsg(filePath)
		},
	)
}

func (m *model) parseCommand(input string) (string, []string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.listDirectory(m.currentPath))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-4, msg.Height-6)
		m.textInput.Width = msg.Width - 10
		return m, nil

	case errMsg:
		m.loading = false
		m.messages = append(m.messages, "error:"+msg.Error())
		m.mode = modeChat
		return m, nil

	case aiResponseMsg:
		m.loading = false
		m.messages = append(m.messages, "ai:"+string(msg))
		return m, nil

	case aiFileContentMsg:
		m.messages = append(m.messages, "info:Contenido generado por IA. Escribiendo en archivo...")
		return m, m.createFileWithContent(msg.fileName, msg.content)

	case fileWrittenMsg:
		m.loading = false
		m.messages = append(m.messages, "info:✅ Archivo creado/escrito: "+string(msg))
		m.mode = modeChat
		m.textInput.Placeholder = "Escribe un mensaje o comando..."
		return m, m.listDirectory(m.currentPath)

	case tea.KeyMsg:
		if m.mode == modeExplorer {
			return m.updateExplorer(msg)
		}
		if m.mode == modeChat || m.mode == modeCreateFileInput || m.mode == modeAIFilenameInput || m.mode == modeAIPromptInput {
			return m.updateTextInputModes(msg)
		}
	}

	if m.loading {
		spinner, cmd := m.spinner.Update(msg)
		m.spinner = spinner
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) updateTextInputModes(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.mode = modeExplorer
		m.textInput.Reset()
		m.textInput.Placeholder = "Write a message or command..."
		return m, nil

	case tea.KeyEnter:
		input := strings.TrimSpace(m.textInput.Value())
		m.textInput.Reset()

		switch m.mode {
		case modeChat:
			if input == "" {
				return m, nil
			}
			m.messages = append(m.messages, "user: "+input)
			cmdName, args := m.parseCommand(input)
			if command, exists := m.commands[cmdName]; exists {
				return m, command.Execute(m, args)
			}
			m.loading = true
			return m, tea.Batch(
				m.spinner.Tick,
				func() tea.Msg {
					res, err := m.aiClient.GetResponse(input)
					if err != nil {
						return errMsg{err}
					}
					return aiResponseMsg(res)
				},
			)

		case modeCreateFileInput:
			filePath := filepath.Join(m.currentPath, input)
			return m, m.createFile(filePath)

		case modeAIFilenameInput:
			m.fileCreationName = filepath.Join(m.currentPath, input)
			m.mode = modeAIPromptInput
			m.textInput.Placeholder = "Describe what the file should accomplish..."
			m.messages = append(m.messages, "info:File to create: "+m.fileCreationName)
			return m, textinput.Blink

		case modeAIPromptInput:
			m.loading = true
			prompt := input
			fileName := m.fileCreationName
			m.messages = append(m.messages, "user: "+prompt)
			m.mode = modeChat

			finalPrompt := fmt.Sprintf("Generate the complete file content for a file named `%s`. The file should accomplish the following: %s. Only output the raw file content, without any explanation or markdown formatting.", filepath.Base(fileName), prompt)

			return m, tea.Batch(
				m.spinner.Tick,
				func() tea.Msg {
					res, err := m.aiClient.GetResponse(finalPrompt)
					if err != nil {
						return errMsg{err}
					}
					return aiFileContentMsg{fileName: fileName, content: res}
				},
			)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *model) updateExplorer(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.mode = modeChat
		return m, nil
	case "c":
		m.mode = modeCreateFileInput
		m.textInput.Placeholder = "Name of the new empty file..."
		m.textInput.Focus()
		return m, textinput.Blink
	case "a":
		m.mode = modeAIFilenameInput
		m.textInput.Placeholder = "Name of the file to generate by AI..."
		m.textInput.Focus()
		return m, textinput.Blink
	case "enter":
		selectedItem, ok := m.list.SelectedItem().(item)
		if !ok {
			return m, nil
		}

		targetPath := filepath.Join(m.currentPath, selectedItem.path)
		if selectedItem.path == ".." {
			targetPath = filepath.Dir(m.currentPath)
		}

		if selectedItem.isDir {
			return m, m.listDirectory(targetPath)
		} else {
			m.mode = modeChat
			m.messages = append(m.messages, "info:Selected: "+targetPath)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var view string

	switch m.mode {
	case modeExplorer:
		view = m.styles.app.Render(m.list.View())
	default:
		var messages strings.Builder
		for _, msg := range m.messages {
			parts := strings.SplitN(msg, ":", 2)
			prefix, content := parts[0], parts[1]

			switch prefix {
			case "user":
				messages.WriteString(m.styles.userMsg.Render("You: "+content) + "\n\n")
			case "ai":
				messages.WriteString(m.styles.aiMsg.Render("🤖: "+content) + "\n\n")
			case "error":
				messages.WriteString(m.styles.errorMsg.Render("❌ Error: "+content) + "\n\n")
			case "info":
				messages.WriteString(m.styles.infoMsg.Render(content) + "\n\n")
			}
		}

		mainContent := lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.header.Render("ANX Agent"),
			messages.String(),
		)

		var status string
		if m.loading {
			status = m.spinner.View() + " Procesando..."
		} else {
			switch m.mode {
			case modeChat:
				status = "MODO: Chat | 'ls' to explore | 'exit' to stop the program."
			case modeCreateFileInput:
				status = "MODO: Create File | 'Enter' to confirm | 'Esc' to cancel"
			case modeAIFilenameInput:
				status = "MODO: AI Filename | 'Enter' to continue | 'Esc' to cancel"
			case modeAIPromptInput:
				status = "MODO: AI Prompt | 'Enter' to generate | 'Esc' to cancel"
			}
		}

		inputView := m.textInput.View()
		statusTextView := m.styles.statusText.Render(status)
		availableWidth := m.width - lipgloss.Width(statusTextView) - 5
		m.textInput.Width = availableWidth

		statusBar := m.styles.statusBar.Width(m.width - 4).Render(
			lipgloss.JoinHorizontal(lipgloss.Left,
				inputView,
				"  ",
				statusTextView,
			),
		)

		view = m.styles.app.Render(lipgloss.JoinVertical(lipgloss.Left, mainContent, statusBar))
	}

	return view
}

func Start(aiClient *ai.Client) {
	p := tea.NewProgram(initialModel(aiClient), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("Error al iniciar la aplicación: ", err)
	}
}
