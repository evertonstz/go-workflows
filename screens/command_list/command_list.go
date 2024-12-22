package cmdlist

import (

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Item struct {
	title string
    desc string
    code string
}

func (i Item) ItemTitle() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) Code() string         { return i.code }
func (i Item) FilterValue() string { return i.title }

type Model struct {
	List list.Model
}
func (m Model) SetTitle(t string) {m.List.Title = t}

func (m Model) Init() tea.Cmd {
	return nil
}

func NewCommandList() list.Model {
    items := []list.Item{
        Item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
        Item{title: "Nutella", desc: "It's good on toast"},
        Item{title: "Bitter melon", desc: "It cools you down"},
        Item{title: "Nice socks", desc: "And by that I mean socks without holes"},
        Item{title: "Eight hours of sleep", desc: "I had this once"},
        Item{title: "Cats", desc: "Usually"},
        Item{title: "Plantasia, the album", desc: "My plants love it too"},
        Item{title: "Pour over coffee", desc: "It takes forever to make though"},
        Item{title: "VR", desc: "Virtual reality...what is there to say?"},
        Item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
        Item{title: "Linux", desc: "Pretty much the best OS"},
        Item{title: "Business school", desc: "Just kidding"},
        Item{title: "Pottery", desc: "Wet clay is a great feeling"},
        Item{title: "Shampoo", desc: "Nothing like clean hair"},
        Item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
        Item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
        Item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
        Item{title: "Stickers", desc: "The thicker the vinyl the better"},
    }

    cmdlist := list.New(items, list.NewDefaultDelegate(), 0, 0)
    cmdlist.Title = "GO Workflows"
    return cmdlist
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.List.View())
}
