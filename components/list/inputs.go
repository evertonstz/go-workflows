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

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type (
	inputs uint

	inputsModel struct {
		Title         textinput.Model
		Description   textinput.Model
		selectedInput inputs
	}

	addNewItemMsg struct {
		Title       string
		Description string
	}
)

const (
	title inputs = iota
	description
	submit
)

func newInputsModel() inputsModel {
	titleModel := textinput.New()
	titleModel.Placeholder = "Title"
	titleModel.Focus()
	descModel := textinput.New()
	descModel.Placeholder = "Description"

	return inputsModel{
		Title:         titleModel,
		Description:   descModel,
		selectedInput: title,
	}
}

func (m inputsModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputsModel) focusInput(i inputs) (inputsModel, tea.Cmd) {
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

func (m inputsModel) Update(msg tea.Msg) (inputsModel, tea.Cmd) {
	titleModel, titleCmd := m.Title.Update(msg)
	descModel, descCmd := m.Description.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "down" {
			switch m.selectedInput {
			case title:
				return m.focusInput(description)
			case description:
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
			case submit:
				return m.focusInput(description)
			}
		}

		if msg.String() == "enter" {
			if m.selectedInput == submit {
				return m, func() tea.Msg {
					return addNewItemMsg{
						Title:       m.Title.Value(),
						Description: m.Description.Value(),
					}
				}
			}
		}
	}

	return inputsModel{
		Title:         titleModel,
		Description:   descModel,
		selectedInput: m.selectedInput,
	}, tea.Batch(titleCmd, descCmd)
}

func (m inputsModel) View() string {
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
