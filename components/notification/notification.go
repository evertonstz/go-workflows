package notification

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var localizer *i18n.Localizer

type Msg struct {
	Text    string
	IsRaw   bool // Indicates if Text is a raw string or a message ID
	DataMap map[string]interface{}
}

type Model struct {
	Text         string
	defaultText  string // This could also be a message ID
	visible      bool
	timerDone    chan struct{}
	isDefaultRaw bool
}

func New(defaultTextOrID string, isRaw bool, loc *i18n.Localizer) Model {
	localizer = loc
	actualDefaultText := defaultTextOrID
	if !isRaw && defaultTextOrID != "" {
		actualDefaultText = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: defaultTextOrID})
	}
	return Model{
		Text:         "",
		defaultText:  actualDefaultText,
		visible:      false,
		timerDone:    nil,
		isDefaultRaw: isRaw,
	}
}

var style = lipgloss.NewStyle().
	PaddingLeft(1).
	PaddingRight(1).
	Background(lipgloss.Color("62")).
	Foreground(lipgloss.Color("230"))

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Msg:
		if msg.IsRaw {
			m.Text = msg.Text
		} else {
			m.Text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: msg.Text, TemplateData: msg.DataMap})
		}
		m.visible = true
		m.timerDone = make(chan struct{})
		return m, startTimer(2*time.Second, m.timerDone)

	case timerMsg:
		if m.timerDone != nil {
			close(m.timerDone)
		}
		m.Text = ""
		m.visible = false
		m.timerDone = nil
	}

	return m, nil
}

func (m Model) View() string {
	if !m.visible {
		if m.defaultText == "" {
			return ""
		}
		return style.SetString(m.defaultText).Render()
	}
	return style.SetString(m.Text).Render()
}

type timerMsg struct{}

func startTimer(duration time.Duration, done chan struct{}) tea.Cmd {
	return func() tea.Msg {
		select {
		case <-time.After(duration):
			return timerMsg{}
		case <-done:
			return nil
		}
	}
}

func ShowNotificationCmd(textOrID string, isRaw bool, dataMap map[string]interface{}) tea.Cmd {
	return func() tea.Msg {
		return Msg{Text: textOrID, IsRaw: isRaw, DataMap: dataMap}
	}
}
