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
	modeAIModifyInput
)

type aiResponseMsg string
type aiFileContentMsg struct {
	fileName string
	content  string
}

type aiModifiedContentMsg struct {
	path    string
	content string
}

type fileReadMsg struct {
	path    string
	content []byte
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
	aiClient                *ai.Client
	commands                map[string]Command
	list                    list.Model
	textInput               textinput.Model
	messages                []string
	spinner                 spinner.Model
	loading                 bool
	width                   int
	height                  int
	mode                    int
	currentPath             string
	styles                  styles
	fileCreationName        string
	fileModificationPath    string
	fileModificationContent string
}

type Command struct {
	Name        string
	Description string
	Execute     func(m *model, args []string) tea.Cmd
}

func initialModel(aiClient *ai.Client) model {
	ti := textinput.New()
	ti.Placeholder = "Write a message or command ('help' to show help)..."
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
	l.Title = "File Explorer"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "back")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open/modify")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "create empty")),
			key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "create with AI")),
		}
	}

	m := model{
		aiClient:    aiClient,
		textInput:   ti,
		messages:    []string{"info:Welcome to ANX Agent. Write 'help' to show help."},
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
			Name: "help", Description: "Show this help message",
			Execute: helpCommand,
		},
		"ls": {
			Name: "ls", Description: "Show the file explorer",
			Execute: listCommand,
		},
		"exit": {
			Name: "exit", Description: "Exit the application",
			Execute: exitCommand,
		},
		"quit": {
			Name: "quit", Description: "Alias for 'exit'",
			Execute: exitCommand,
		},
	}
}

func helpCommand(m *model, args []string) tea.Cmd {
	helpText := "Available commands:\n"
	for name, cmd := range m.commands {
		helpText += fmt.Sprintf("  %-15s %s", name, cmd.Description)
		m.messages = append(m.messages, "info:"+helpText)
		helpText = ""
	}
	m.messages = append(m.messages, "info:\nIn the explorer ('ls'):\n  'enter' to open dir or modify file\n  'c' to create empty file\n  'a' to create file with AI")
	return nil
}

func listCommand(m *model, args []string) tea.Cmd {
	m.mode = modeExplorer
	return m.listDirectory(m.currentPath)
}

func exitCommand(m *model, args []string) tea.Cmd {
	m.messages = append(m.messages, "info:Goodbye!")
	return tea.Quit
}

func (m *model) readFileContent(filePath string) tea.Cmd {
	return func() tea.Msg {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return errMsg{err}
		}
		return fileReadMsg{path: filePath, content: content}
	}
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
		m.messages = append(m.messages, "info:Content generated by AI. Writing to file...")
		return m, m.createFileWithContent(msg.fileName, msg.content)

	case fileReadMsg:
		m.loading = false
		m.mode = modeAIModifyInput
		m.fileModificationPath = msg.path
		m.fileModificationContent = string(msg.content)
		m.messages = append(m.messages, "info:File '"+filepath.Base(msg.path)+"' read. How do you want to modify it?")
		m.textInput.Placeholder = "Ej: 'Add a comment to the main function'..."
		m.textInput.Focus()
		return m, textinput.Blink

	case aiModifiedContentMsg:
		m.messages = append(m.messages, "info:Content modified by AI. Writing changes...")
		return m, m.createFileWithContent(msg.path, msg.content)

	case fileWrittenMsg:
		m.loading = false
		m.messages = append(m.messages, "info:✅ File created/modified: "+string(msg))
		m.mode = modeExplorer
		m.textInput.Placeholder = "Write a message or command..."
		return m, m.listDirectory(m.currentPath)

	case tea.KeyMsg:
		if m.mode == modeExplorer {
			return m.updateExplorer(msg)
		}
		if m.mode == modeChat || m.mode == modeCreateFileInput || m.mode == modeAIFilenameInput || m.mode == modeAIPromptInput || m.mode == modeAIModifyInput {
			return m.updateTextInputModes(msg)
		}
	}

	if m.loading {
		spinner, cmd := m.spinner.Update(msg)
		m.spinner = spinner
		return m, cmd
	}

	return m, nil
}

func (m *model) updateTextInputModes(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.mode = modeExplorer
		m.textInput.Reset()
		m.textInput.Placeholder = "Write a message or command..."
		m.fileCreationName = ""
		m.fileModificationPath = ""
		m.fileModificationContent = ""
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
			m.textInput.Placeholder = "Describe what sould do the file..."
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

		case modeAIModifyInput:
			m.loading = true
			instructions := input
			originalContent := m.fileModificationContent
			filePath := m.fileModificationPath
			m.messages = append(m.messages, "user: "+instructions)
			m.mode = modeChat

			finalPrompt := fmt.Sprintf(
				"You are an expert file editor. The user wants to modify a file. Below is the original content of the file and the user's instructions. Your task is to return the *entire*, *new* content of the file with the modifications applied. \n\nIMPORTANT: Only output the raw, complete, modified file content. Do not include any explanations, greetings, or markdown code fences like ```go ... ```.\n\n--- ORIGINAL FILE CONTENT ---\n%s\n\n--- USER INSTRUCTIONS ---\n%s",
				originalContent,
				instructions,
			)

			return m, tea.Batch(
				m.spinner.Tick,
				func() tea.Msg {
					res, err := m.aiClient.GetResponse(finalPrompt)
					if err != nil {
						return errMsg{err}
					}
					return aiModifiedContentMsg{path: filePath, content: res}
				},
			)
		}
	}

	var cmd tea.Cmd
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
		m.textInput.Placeholder = "New file name (empty)..."
		m.textInput.Focus()
		return m, textinput.Blink
	case "a":
		m.mode = modeAIFilenameInput
		m.textInput.Placeholder = "File name to generate by AI..."
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
			m.loading = true
			m.messages = append(m.messages, "info:Reading file "+targetPath+" to modify...")
			return m, m.readFileContent(targetPath)
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
			status = m.spinner.View() + " Processing..."
		} else {
			switch m.mode {
			case modeChat:
				status = "MODO: Chat | 'ls' to explore | 'exit' to exit"
			case modeCreateFileInput:
				status = "MODO: Create File | 'Enter' to confirm | 'Esc' to cancel"
			case modeAIFilenameInput:
				status = "MODO: File Name (AI) | 'Enter' to continue | 'Esc' to cancel"
			case modeAIPromptInput:
				status = "MODO: Description (AI) | 'Enter' to generate | 'Esc' to cancel"
			case modeAIModifyInput:
				status = "MODO: Modify with AI | 'Enter' to send | 'Esc' to cancel"
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
		log.Fatal("Error starting the application: ", err)
	}
}
