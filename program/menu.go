package program

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Menu struct {
	Options       []MenuItem
	SelectedIndex int
	Exit          bool
}

type MenuItem struct {
	Action func() error
	Text   string
}

func (m Menu) Init() tea.Cmd { return nil }

func (m Menu) View() string {
	var options []string
	for i, o := range m.Options {
		if i == m.SelectedIndex {
			options = append(options, fmt.Sprintf("-> %s", o.Text))
		} else {
			options = append(options, fmt.Sprintf("   %s", o.Text))
		}
	}
	return fmt.Sprintf(`%s

Press enter/return to select a list item, arrow keys to move, or Ctrl+C to exit.`,
		strings.Join(options, "\n"))
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		switch v.String() {
		case "ctrl+c", "q":
			println()
			os.Exit(1)
		case "down", "right", "up", "left":
			return m.moveCursor(v), nil
		case "enter", "return", " ":
			return m, tea.Quit
		}
	case int:

	}
	return m, nil
}

func (m Menu) moveCursor(msg tea.KeyMsg) Menu {
	switch msg.String() {
	case "up", "left":
		m.SelectedIndex--
	case "down", "right":
		m.SelectedIndex++
	default:
		// do nothing
	}

	optCount := len(m.Options)
	m.SelectedIndex = (m.SelectedIndex + optCount) % optCount
	return m
}
