package addnew

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
