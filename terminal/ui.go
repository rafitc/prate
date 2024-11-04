package terminal

// Terminal UI related code. Built on top of bubble tea cli chat example

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Terminal struct {
	Message NewMsg
}

type NewMsg struct {
	User string
	Body string
}

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func InitTerminal() *tea.Program {
	terminalSession := tea.NewProgram(initialModel())
	return terminalSession
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		return m, nil

	case NewMsg:
		v := msg.Body
		if v == "" {
			// Don't send empty messages.
			return m, nil
		}
		// All things goes here
		// Send message to everyone from the user

		// Simulate sending a message. In your application you'll want to
		// also return a custom command to send the message off to
		// a server.
		m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("I'm %s: ", msg.User))+v)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		// m.textarea.Reset()
		m.viewport.GotoBottom()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			// Quit.
			fmt.Println(m.textarea.Value())
			return m, tea.Quit

		case "enter":
			v := m.textarea.Value()
			if v == "" {
				// Don't send empty messages.
				return m, nil
			}
			// All things goes here
			// Send message to everyone from the user

			// Simulate sending a message. In your application you'll want to
			// also return a custom command to send the message off to
			// a server.
			m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%s: ", "unknow"))+v)
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			return m, nil

		default:
			// Send all other keypresses to the textarea.
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}

	case cursor.BlinkMsg:
		// Textarea should also process cursor blinks.
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
