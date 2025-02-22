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

	focusedButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("[ Submit ]")
	blurredButton = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("[ Submit ]")

	mainStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))
)

type (
	inputs uint

	Styles struct {
		focusedInput    lipgloss.Style
		blurredInput    lipgloss.Style
		focusedTextArea lipgloss.Style
		blurredTextArea lipgloss.Style
		focusedButton   string
		blurredButton   string
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
	title inputs = iota
	description
	textArea
	submit
)

func (m Model) SetSize(width, height int) {
	m.Title.Width = width
	m.Description.Width = width
	m.TextArea.SetWidth(width)
	m.TextArea.SetHeight(height)

	m.styles.focusedInput = m.styles.focusedInput.Width(width)
	m.styles.blurredInput = m.styles.blurredInput.Width(width)
	m.styles.focusedTextArea = m.styles.focusedTextArea.Width(width)
	m.styles.blurredTextArea = m.styles.blurredTextArea.Width(width)
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
			focusedInput:    focusedStyle,
			blurredInput:    blurredStyle,
			focusedTextArea: focusedTextAreaStyle,
			blurredTextArea: blurredTextAreaStyle,
			focusedButton:   focusedButton,
			blurredButton:   blurredButton,
		},
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
				return m, shared.AddNewItemCmd(m.Title.Value(), m.Description.Value(), m.TextArea.Value())
				// return m, func() tea.Msg {
				// 	return DidAddNewItemMsg{
				// 		Title:       m.Title.Value(),
				// 		Description: m.Description.Value(),
				// 		CommandText: m.TextArea.Value(),
				// 	}
				// }
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
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top, m.styles.focusedInput.Render(m.Title.View()),
			m.styles.blurredInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			m.styles.blurredButton))
	case description:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top, m.styles.blurredInput.Render(m.Title.View()),
			m.styles.focusedInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			m.styles.blurredButton))
	case submit:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top, m.styles.blurredInput.Render(m.Title.View()),
			m.styles.blurredInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			m.styles.focusedButton))
	default:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top, m.styles.blurredInput.Render(m.Title.View()),
			m.styles.blurredInput.Render(m.Description.View()),
			focusedTextAreaStyle.Render(m.TextArea.View()),
			m.styles.blurredButton))
	}
}
