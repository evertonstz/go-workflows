package addnew

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedTextAreaStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("205"))
	blurredTextAreaStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))

	focusedButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredButton = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

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

	Notifications struct {
		fillAllFields string
	}

	Model struct {
		Title         textinput.Model
		Description   textinput.Model
		TextArea      textarea.Model
		selectedInput inputs
		styles        Styles
		notifications Notifications
		Keys          helpkeys.AddNewKeyMap
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
	i18n := di.GetService[*services.I18nService](di.I18nServiceKey)

	titleModel := textinput.New()
	titleModel.Placeholder = i18n.Translate("title_placeholder")
	titleModel.Focus()
	descModel := textinput.New()
	descModel.Placeholder = i18n.Translate("description_placeholder")
	textareaModel := textarea.New()
	textareaModel.Placeholder = i18n.Translate("command_placeholder")
	textareaModel.Prompt = ""
	textareaModel.ShowLineNumbers = false

	focusedSaveButton := focusedButton.Render(fmt.Sprintf("[ %s ]", i18n.Translate("save_button_label")))
	blurredSaveButton := blurredButton.Render(fmt.Sprintf("[ %s ]", i18n.Translate("save_button_label")))
	focusedCloseButton := focusedButton.Render(fmt.Sprintf("[ %s ]", i18n.Translate("cancel_button_label")))
	blurredCloseButton := blurredButton.Render(fmt.Sprintf("[ %s ]", i18n.Translate("cancel_button_label")))

	return Model{
		Title:         titleModel,
		Description:   descModel,
		TextArea:      textareaModel,
		selectedInput: title,
		Keys:          helpkeys.NewAddNewKeys(i18n),
		notifications: Notifications{
			fillAllFields: i18n.Translate("error_fill_all_fields"),
		},
		styles: Styles{
			focusedInput:       focusedStyle,
			blurredInput:       blurredStyle,
			focusedTextArea:    focusedTextAreaStyle,
			blurredTextArea:    blurredTextAreaStyle,
			focusedButton:      focusedSaveButton,
			blurredButton:      blurredSaveButton,
			blurredCloseButton: blurredCloseButton,
			focusedCloseButton: focusedCloseButton,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
