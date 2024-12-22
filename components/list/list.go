package list

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertonstz/go-workflows/shared"
	// "github.com/evertonstz/go-workflows/shared"
)

type addNewItemState uint

const (
	viewOff addNewItemState = iota
	viewOn
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
	dateAdded   time.Time
	dateUpdated time.Time
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Model struct {
	state addNewItemState
	input textinput.Model
	list  list.Model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+a" {
			if m.state == viewOff {
				m.state = viewOn
			} else {
				m.state = viewOff
			}
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			switch m.state {
			case viewOff:
				selectedItem := m.list.SelectedItem()
				if selectedItem != nil {
					if selected, ok := selectedItem.(item); ok {
						return m, func() tea.Msg {
							return shared.SelectedItemMsg{Item: selected.title}
						}
					}
				}
			case viewOn:
				var c tea.Cmd
				m.list.InsertItem(0, item{title: m.input.Value(), desc: "", dateAdded: time.Now(), dateUpdated: time.Now()})
				m.input.Reset()
				m.list, c = m.list.Update(msg)
				return m, c
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-6)
	}

	switch m.state {
	case viewOff:
		var c tea.Cmd
		m.list, c = m.list.Update(msg)
		cmds = append(cmds, c)
	case viewOn:
		var c tea.Cmd
		m.input, c = m.input.Update(msg)
		cmds = append(cmds, c)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var v string
	listView := docStyle.Render(m.list.View())
	inputView := docStyle.Render(m.input.View())
	switch m.state {
	case viewOff:
		v = listView
	case viewOn:
		v = lipgloss.JoinVertical(lipgloss.Top, listView, inputView)
	}

	return v
}

func New() Model {
	items := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house", dateAdded: time.Now(), dateUpdated: time.Now()},
		item{title: "Nutella", desc: "It's good on toast", dateAdded: time.Now(), dateUpdated: time.Now()},
		item{title: "Bitter melon", desc: "It cools you down", dateAdded: time.Now(), dateUpdated: time.Now()},
		// item{title: "Nice socks", desc: "And by that I mean socks without holes"},
		// item{title: "Eight hours of sleep", desc: "I had this once"},
		// item{title: "Cats", desc: "Usually"},
		// item{title: "Plantasia, the album", desc: "My plants love it too"},
		// item{title: "Pour over coffee", desc: "It takes forever to make though"},
		// item{title: "VR", desc: "Virtual reality...what is there to say?"},
		// item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		// item{title: "Linux", desc: "Pretty much the best OS"},
		// item{title: "Business school", desc: "Just kidding"},
		// item{title: "Pottery", desc: "Wet clay is a great feeling"},
		// item{title: "Shampoo", desc: "Nothing like clean hair"},
		// item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
		// item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
		// item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		// item{title: "Stickers", desc: "The thicker the vinyl the better"},
		// item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		// item{title: "Warm light", desc: "Like around 2700 Kelvin"},
		// item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		// item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		// item{title: "Terrycloth", desc: "In other words, towel fabric"},
	}
	ti := textinput.New()
	ti.Placeholder = "New command name..."
	ti.Focus()
	m := Model{list: list.New(items, list.NewDefaultDelegate(), 0, 0), input: ti}
	m.list.Title = "Workflows"
	m.Init()
	return m
}
