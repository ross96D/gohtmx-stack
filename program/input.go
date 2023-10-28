package program

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type TextInput struct {
	TextInput textinput.Model
	Err       error
}

func NewTextInput(placeholder string) TextInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.PromptStyle.PaddingRight(3)

	return TextInput{
		TextInput: ti,
		Err:       nil,
	}
}

func (m TextInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEscape:
			os.Exit(1)
		}

	// We handle errors just like any other message
	case errMsg:
		m.Err = msg
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m TextInput) View() string {
	return fmt.Sprintf(
		"What’s the name of the go module?\n\n\n\n%s\n\n",
		m.TextInput.View(),
	)
}
