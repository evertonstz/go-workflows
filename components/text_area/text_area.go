package textarea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/shared"
)

type errMsg error

type Model struct {
	TextArea textarea.Model
	err      error
}

func New() Model {
	ti := textarea.New()
	ti.Placeholder = "Paste or type your command here..."
	ti.Focus()

	return Model{
		TextArea: ti,
		err:      nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case shared.ItemMsg:
		m.TextArea.SetValue(msg.Command)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.TextArea.Focused() {
				m.TextArea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlS:
			return m, func() tea.Msg {
				return shared.SaveItem{Command: m.TextArea.Value()}
			}
		default:
			if !m.TextArea.Focused() {
				cmd = m.TextArea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.TextArea, cmd = m.TextArea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"Tell me a story.\n\n%s\n\n%s",
		m.TextArea.View(),
		"(ctrl+c to quit)",
	) + "\n\n"
}