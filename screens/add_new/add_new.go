package addnew

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type (
	inputs uint

	Model struct {
		Title         textinput.Model
		Description   textinput.Model
		TextArea      textarea.Model
		selectedInput inputs
	}

	DidAddNewItemMsg struct {
		Title       string
		Description string
		CommandText string
	}
)

const (
	title inputs = iota
	description
	textArea
	submit
)

func New() Model {
	titleModel := textinput.New()
	titleModel.Placeholder = "Title"
	titleModel.Focus()
	descModel := textinput.New()
	descModel.Placeholder = "Description"
	textareaModel := textarea.New()

	return Model{
		Title:         titleModel,
		Description:   descModel,
		TextArea:      textareaModel,
		selectedInput: title,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) focusInput(i inputs) (Model, tea.Cmd) {
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
	}

	return m, nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	titleModel, titleCmd := m.Title.Update(msg)
	descModel, descCmd := m.Description.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "down" {
			switch m.selectedInput {
			case title:
				return m.focusInput(description)
			case description:
				return m.focusInput(textArea)
			case textArea:
				return m.focusInput(submit)
			case submit:
				return m, nil
			}

		}
		if msg.String() == "up" {
			switch m.selectedInput {
			case title:
				return m, nil
			case description:
				return m.focusInput(title)
			case textArea:
				return m.focusInput(description)
			case submit:
				return m.focusInput(textArea)
			}
		}

		if msg.String() == "enter" {
			if m.selectedInput == submit {
				return m, func() tea.Msg {
					return DidAddNewItemMsg{
						Title:       m.Title.Value(),
						Description: m.Description.Value(),
						CommandText: m.TextArea.Value(),
					}
				}
			}
		}
	}
	textModel, textCmd := m.TextArea.Update(msg)
	return Model{
		Title:         titleModel,
		Description:   descModel,
		TextArea:      textModel,
		selectedInput: m.selectedInput,
	}, tea.Batch(titleCmd, descCmd, textCmd)
}

func (m Model) View() string {
	switch m.selectedInput {
	case title:
		return lipgloss.JoinVertical(lipgloss.Top, focusedStyle.Render(m.Title.View()),
			blurredStyle.Render(m.Description.View()),
			m.TextArea.View(),
			blurredButton)
	case description:
		return lipgloss.JoinVertical(lipgloss.Top, blurredStyle.Render(m.Title.View()),
			focusedStyle.Render(m.Description.View()),
			m.TextArea.View(),
			blurredButton)
	case submit:
		return lipgloss.JoinVertical(lipgloss.Top, blurredStyle.Render(m.Title.View()),
			blurredStyle.Render(m.Description.View()),
			m.TextArea.View(),
			focusedButton)
	default:
		return lipgloss.JoinVertical(lipgloss.Top, blurredStyle.Render(m.Title.View()),
			blurredStyle.Render(m.Description.View()),
			m.TextArea.View(),
			blurredButton)
	}
}
