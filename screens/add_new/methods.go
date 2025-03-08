package addnew

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) SetSize(width, height int) {
	m.Title.Width = width
	m.Description.Width = width
	m.TextArea.SetWidth(width)
	m.TextArea.SetHeight(height)
}

func (m *Model) SetValues(title, description, command string) {
	m.Title.SetValue(title)
	m.Description.SetValue(description)
	m.TextArea.SetValue(command)
}

func (m *Model) ResetForm() {
	m.Title.SetValue("")
	m.Description.SetValue("")
	m.TextArea.SetValue("")
	m.focusInput(title)
}

func (m *Model) focusInput(i inputs) (Model, tea.Cmd) {
	switch i {
	case title:
		m.Title.Focus()
		m.Description.Blur()
		m.TextArea.Blur()
		m.selectedInput = title
	case description:
		m.Title.Blur()
		m.Description.Focus()
		m.TextArea.Blur()
		m.selectedInput = description
	case textArea:
		m.Title.Blur()
		m.Description.Blur()
		m.TextArea.Focus()
		m.selectedInput = textArea
	case submit:
		m.Title.Blur()
		m.Description.Blur()
		m.TextArea.Blur()
		m.selectedInput = submit
	case close:
		m.Title.Blur()
		m.Description.Blur()
		m.TextArea.Blur()
		m.selectedInput = close
	}

	return *m, nil
}

func (m Model) isFormValid() bool {
	return m.Title.Value() != "" && m.Description.Value() != "" && m.TextArea.Value() != ""
}
