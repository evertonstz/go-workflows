package addnew

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	// localizer *i18n.Localizer // Removed package-level variable

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedTextAreaStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("205"))
	blurredTextAreaStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))

	focusedButton      string
	blurredButton      string
	focusedCloseButton string
	blurredCloseButton string

	mainStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))
)

func setButtons(loc *i18n.Localizer) { // Accept localizer as parameter
	saveLabel := loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "save_button_label"})
	cancelLabel := loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel_button_label"})

	focusedButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("[ " + saveLabel + " ]")
	blurredButton = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("[ " + saveLabel + " ]")
	focusedCloseButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("[ " + cancelLabel + " ]")
	blurredCloseButton = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("[ " + cancelLabel + " ]")
}

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
		localizer     *i18n.Localizer
	}
)

const (
	close inputs = iota
	title
	description
	textArea
	submit
)

func New(loc *i18n.Localizer) Model {
	setButtons(loc) // Call with the passed localizer

	titleModel := textinput.New()
	titleModel.Placeholder = loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "title_placeholder"})
	titleModel.Focus()
	descModel := textinput.New()
	descModel.Placeholder = loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "description_placeholder"})
	textareaModel := textarea.New()
	textareaModel.Placeholder = loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "command_placeholder"})
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
		localizer: loc, // Store the localizer in the model instance
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
