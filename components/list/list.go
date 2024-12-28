package list

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertonstz/go-workflows/components/persist"
	"github.com/evertonstz/go-workflows/shared"
)

const (
	addNewOff inputs = iota
	addNewOn
)

type item struct {
	title, desc, command string
	dateAdded            time.Time
	dateUpdated          time.Time
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (i item) Title() string          { return i.title }
func (i item) Description() string    { return i.desc }
func (i item) Command() string        { return i.command }
func (i item) DateAdded() time.Time   { return i.dateAdded }
func (i item) DateUpdated() time.Time { return i.dateUpdated }
func (i item) FilterValue() string    { return i.title }

type Model struct {
	state  inputs
	inputs InputsModel
	list   list.Model
}

func (m Model) InputOn() bool {
	return m.state == addNewOn
}

func (m Model) CurentItem() item {
	return m.list.SelectedItem().(item)
}

func (m Model) AllItems() []item {
	var items []item
	for _, i := range m.list.Items() {
		items = append(items, i.(item))
	}
	return items
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) changeState(v inputs) inputs {
	m.state = v
	return m.state
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case AddNewItemMsg:
		if m.inputs.Title.Value() == "" || m.inputs.Description.Value() == "" {
			return m, nil
		}
		newItem := item{
			title:       msg.Title,
			desc:        msg.Description,
			dateAdded:   time.Now(),
			dateUpdated: time.Now(),
		}
		m.list.InsertItem(len(m.list.Items()), newItem)
		m.changeState(addNewOff)
		return m, func() tea.Msg {
			return shared.SaveCommandMsg{Command: ""}
		}

	case shared.SaveCommandMsg:
		selectedItem := m.list.SelectedItem()
		if selectedItem != nil {
			if selected, ok := selectedItem.(item); ok {
				selected.command = msg.Command
				selected.dateUpdated = time.Now()
				m.list.SetItem(m.list.Index(), item{
					title:       selected.title,
					desc:        selected.desc,
					command:     selected.command,
					dateAdded:   selected.dateAdded,
					dateUpdated: selected.dateUpdated})
			}
		}
	case persist.PersistionFileLoadedMsg:
		var data []list.Item
		for _, i := range msg.Items.Items {
			data = append(data, list.Item(item{
				title:       i.Title,
				desc:        i.Desc,
				command:     i.Command,
				dateAdded:   i.DateAdded,
				dateUpdated: i.DateUpdated}))
		}
		m.list.SetItems(data)
	case tea.KeyMsg:
		if msg.String() == "ctrl+a" {
			if m.state == addNewOff {
				m.changeState(addNewOn)
			} else {
				m.changeState(addNewOff)
			}
		}
		if msg.String() == "y" {
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				if selected, ok := selectedItem.(item); ok {
					return m, func() tea.Msg {
						return shared.CopyToClipboardMsg{Command: selected.command}
					}

				}
			}
		}
		if msg.String() == "enter" {
			switch m.state {
			case addNewOn:
				var c tea.Cmd
				m.inputs, c = m.inputs.Update(msg)
				return m, c
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-6)
	}

	switch m.state {
	case addNewOff:
		var c tea.Cmd
		m.list, c = m.list.Update(msg)
		cmds = append(cmds, c)
	case addNewOn:
		var c tea.Cmd
		m.inputs, c = m.inputs.Update(msg)
		cmds = append(cmds, c)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var v string
	listView := docStyle.Render(m.list.View())
	inputView := docStyle.Render(m.inputs.View())
	switch m.state {
	case addNewOff:
		v = listView
	case addNewOn:
		v = lipgloss.JoinVertical(lipgloss.Top, listView, inputView)
	}

	return v
}

func New() Model {
	m := Model{list: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0), inputs: NewInputsModel()}
	m.list.Title = "Workflows"
	m.Init()
	return m
}
