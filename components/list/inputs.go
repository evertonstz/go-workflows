package list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	// cursorStyle         = focusedStyle
	// noStyle             = lipgloss.NewStyle()
	// helpStyle           = blurredStyle
	// cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type inputs uint

const (
	title inputs = iota
	description
	submit
)

type InputsModel struct {
	Title         textinput.Model
	Description   textinput.Model
	selectedInput inputs
}

type AddNewItemMsg struct {
	Title       string
	Description string
}

func NewInputsModel() InputsModel {
	titleModel := textinput.New()
	titleModel.Placeholder = "Title"
	titleModel.Focus()
	descModel := textinput.New()
	descModel.Placeholder = "Description"

	return InputsModel{
		Title:         titleModel,
		Description:   descModel,
		selectedInput: title,
	}
}

func (m InputsModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputsModel) Focus(i inputs) (InputsModel, tea.Cmd) {
	switch i {
	case title:
		m.Title.Focus()
		m.Description.Blur()
		m.selectedInput = title
	case description:
		m.Title.Blur()
		m.Description.Focus()
		m.selectedInput = description
	case submit:
		m.Title.Blur()
		m.Description.Blur()
		m.selectedInput = submit
	}

	return m, nil
}

func (m InputsModel) Update(msg tea.Msg) (InputsModel, tea.Cmd) {
	titleModel, titleCmd := m.Title.Update(msg)
	descModel, descCmd := m.Description.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "down" {
			switch m.selectedInput {
			case title:
				return m.Focus(description)
			case description:
				return m.Focus(submit)
			case submit:
				return m, nil
			}

		}
		if msg.String() == "up" {
			switch m.selectedInput {
			case title:
				return m, nil
			case description:
				return m.Focus(title)
			case submit:
				return m.Focus(description)
			}
		}

		if msg.String() == "enter" {
			if m.selectedInput == submit {
				return m, func() tea.Msg {
					return AddNewItemMsg{
						Title:       m.Title.Value(),
						Description: m.Description.Value(),
					}
				}
			}
		}
	}

	return InputsModel{
		Title:         titleModel,
		Description:   descModel,
		selectedInput: m.selectedInput,
	}, tea.Batch(titleCmd, descCmd)
}

func (m InputsModel) View() string {
	switch m.selectedInput {
	case title:
		return lipgloss.JoinVertical(lipgloss.Top, focusedStyle.Render(m.Title.View()),
			blurredStyle.Render(m.Description.View()),
			blurredButton)
	case description:
		return lipgloss.JoinVertical(lipgloss.Top, blurredStyle.Render(m.Title.View()),
			focusedStyle.Render(m.Description.View()),
			blurredButton)
	case submit:
		return lipgloss.JoinVertical(lipgloss.Top, blurredStyle.Render(m.Title.View()),
			blurredStyle.Render(m.Description.View()),
			focusedButton)
	default:
		return lipgloss.JoinVertical(lipgloss.Top, blurredStyle.Render(m.Title.View()),
			blurredStyle.Render(m.Description.View()),
			blurredButton)
	}
}
