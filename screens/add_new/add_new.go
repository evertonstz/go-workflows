package addnew

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertonstz/go-workflows/shared"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedTextAreaStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("205"))
	blurredTextAreaStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))

	focusedButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("[ Save ]")
	blurredButton = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("[ Save ]")

	focusedCloseButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("[ Cancel ]")
	blurredCloseButton = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("[ Cancel ]")

	mainStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))
)

type (
	inputs uint

	Styles struct {
		focusedInput       lipgloss.Style
		blurredInput       lipgloss.Style
		focusedTextArea    lipgloss.Style
		blurredTextArea    lipgloss.Style
		focusedButton      string
		blurredButton      string
		blurredCloseButton string
		focusedCloseButton string
	}

	Model struct {
		Title         textinput.Model
		Description   textinput.Model
		TextArea      textarea.Model
		selectedInput inputs
		styles        Styles
	}
)

const (
	close inputs = iota
	title
	description
	textArea
	submit
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

func New() Model {
	titleModel := textinput.New()
	titleModel.Placeholder = "Title"
	titleModel.Focus()
	descModel := textinput.New()
	descModel.Placeholder = "Description"
	textareaModel := textarea.New()
	textareaModel.Placeholder = "Paste or type your command here..."
	textareaModel.Prompt = ""
	textareaModel.ShowLineNumbers = false

	return Model{
		Title:         titleModel,
		Description:   descModel,
		TextArea:      textareaModel,
		selectedInput: title,
		styles: Styles{
			focusedInput:       focusedStyle,
			blurredInput:       blurredStyle,
			focusedTextArea:    focusedTextAreaStyle,
			blurredTextArea:    blurredTextAreaStyle,
			focusedButton:      focusedButton,
			blurredButton:      blurredButton,
			blurredCloseButton: blurredCloseButton,
			focusedCloseButton: focusedCloseButton,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	titleModel, titleCmd := m.Title.Update(msg)
	descModel, descCmd := m.Description.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			switch m.selectedInput {
			case title:
				return m.focusInput(description)
			case description:
				return m.focusInput(textArea)
			case textArea:
				return m.focusInput(submit)
			case submit, close:
				return m, nil
			}
		case "up":
			switch m.selectedInput {
			case title:
				return m, nil
			case description:
				return m.focusInput(title)
			case textArea:
				return m.focusInput(description)
			case submit, close:
				return m.focusInput(textArea)
			}
		case "right":
			if m.selectedInput == submit {
				return m.focusInput(close)
			}
		case "left":
			if m.selectedInput == close {
				return m.focusInput(submit)
			}
		case "esc":
			m.ResetForm()
			return m, shared.CloseAddNewScreenCmd()
		case "enter":
			switch m.selectedInput {
			case submit:
				if m.isFormValid() {
					var title, description, command string
					title = m.Title.Value()
					description = m.Description.Value()
					command = m.TextArea.Value()

					m.ResetForm()
					return m, shared.AddNewItemCmd(title, description, command)
				}
			case close:
				m.ResetForm()
				return m, shared.CloseAddNewScreenCmd()
			}
		}
	}
	textModel, textCmd := m.TextArea.Update(msg)
	return Model{
		Title:         titleModel,
		Description:   descModel,
		TextArea:      textModel,
		selectedInput: m.selectedInput,
		styles:        m.styles,
	}, tea.Batch(titleCmd, descCmd, textCmd)
}

func (m Model) View() string {
	switch m.selectedInput {
	case title:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.focusedInput.Render(m.Title.View()),
			m.styles.blurredInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			lipgloss.NewStyle().Align(lipgloss.Right).Width(m.Title.Width).Render(
				lipgloss.JoinHorizontal(lipgloss.Top, m.styles.blurredButton, m.styles.blurredCloseButton))))
	case description:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.blurredInput.Render(m.Title.View()),
			m.styles.focusedInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			lipgloss.NewStyle().Align(lipgloss.Right).Width(m.Title.Width).Render(
				lipgloss.JoinHorizontal(lipgloss.Top, m.styles.blurredButton, m.styles.blurredCloseButton))))
	case submit:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.blurredInput.Render(m.Title.View()),
			m.styles.blurredInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			lipgloss.NewStyle().Align(lipgloss.Right).Width(m.Title.Width).Render(
				lipgloss.JoinHorizontal(lipgloss.Top, m.styles.focusedButton, m.styles.blurredCloseButton))))
	case close:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.blurredInput.Render(m.Title.View()),
			m.styles.blurredInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			lipgloss.NewStyle().Align(lipgloss.Right).Width(m.Title.Width).Render(
				lipgloss.JoinHorizontal(lipgloss.Top, m.styles.blurredButton, m.styles.focusedCloseButton))))
	default:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.blurredInput.Render(m.Title.View()),
			m.styles.blurredInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			lipgloss.NewStyle().Align(lipgloss.Right).Width(m.Title.Width).Render(
				lipgloss.JoinHorizontal(lipgloss.Top, m.styles.blurredButton, m.styles.blurredCloseButton))))
	}
}
