package list

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertonstz/go-workflows/components/persist"
	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
)

type myItem struct {
	title, desc, command string
	dateAdded            time.Time
	dateUpdated          time.Time
}

const (
	addNewOff inputs = iota
	addNewOn
)

var docStyle = lipgloss.NewStyle().MarginTop(1)

func (i myItem) Title() string          { return i.title }
func (i myItem) Description() string    { return i.desc }
func (i myItem) Command() string        { return i.command }
func (i myItem) DateAdded() time.Time   { return i.dateAdded }
func (i myItem) DateUpdated() time.Time { return i.dateUpdated }
func (i myItem) FilterValue() string    { return i.title }

type Model struct {
	state  inputs
	inputs inputsModel
	list   list.Model
}

func (m Model) InputOn() bool {
	return m.state == addNewOn
}

func (m Model) CurentItem() myItem {
	return m.list.SelectedItem().(myItem)
}

func (m Model) AllItems() []myItem {
	var items []myItem
	for _, i := range m.list.Items() {
		items = append(items, i.(myItem))
	}
	return items
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) showAddNew() {
	m.state = addNewOn
}

func (m *Model) hideAddNew() {
	m.state = addNewOff
}

func (m *Model) SetSize(width, height int) {
	m.list.SetSize(width, height)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case addNewItemMsg:
		if m.inputs.Title.Value() == "" || m.inputs.Description.Value() == "" {
			return m, nil
		}
		newItem := myItem{
			title:       msg.Title,
			desc:        msg.Description,
			dateAdded:   time.Now(),
			dateUpdated: time.Now(),
		}
		m.list.InsertItem(len(m.list.Items()), newItem)
		m.hideAddNew()
		return m, nil
	case shared.DidUpdateItemMsg:
		for i, item := range m.list.Items() {
			if item.(myItem).title == msg.Item.Title {
				m.list.SetItem(i, myItem{
					title:       msg.Item.Title,
					desc:        msg.Item.Desc,
					command:     msg.Item.Command,
					dateAdded:   msg.Item.DateAdded,
					dateUpdated: msg.Item.DateUpdated,
				})
			}
		}

	case persist.LoadedDataFileMsg:
		var data []list.Item
		for _, i := range msg.Items.Items {
			data = append(data, list.Item(myItem{
				title:       i.Title,
				desc:        i.Desc,
				command:     i.Command,
				dateAdded:   i.DateAdded,
				dateUpdated: i.DateUpdated}))
		}
		m.list.SetItems(data)
	case tea.KeyMsg:
		if msg.String() == "esc" {
			if m.state == addNewOn {
				m.hideAddNew()
				return m, nil
			}
		}
		if msg.String() == "a" {
			if m.state == addNewOff {
				m.showAddNew()
				return m, nil
			}
		}
		if msg.String() == "y" {
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				if selected, ok := selectedItem.(myItem); ok {
					return m, shared.CopyToClipboardCmd(selected.command)
				}
			}
		}
		if msg.String() == "enter" {
			switch m.state {
			case addNewOff:
				return m, shared.SetCurrentItemCmd(models.Item{
					Title:       m.CurentItem().Title(),
					Desc:        m.CurentItem().Description(),
					Command:     m.CurentItem().Command(),
					DateAdded:   m.CurentItem().DateAdded(),
					DateUpdated: m.CurentItem().DateUpdated()},
				)
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
	// listView := docStyle.Render(m.list.View())
	inputView := docStyle.Render(m.inputs.View())
	inputViewHeight := strings.Count(inputView, "\n")
	switch m.state {
	case addNewOff:
		v = docStyle.Render(m.list.View())
	case addNewOn:
		v = lipgloss.JoinVertical(lipgloss.Top,
			docStyle.
				Height(m.list.Height()-inputViewHeight).
				Render(m.list.View()),
			inputView)
	}

	return v
}

func New() Model {
	m := Model{list: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0), inputs: newInputsModel()}
	m.list.Title = "Workflows"
	m.list.SetShowHelp(false)
	m.Init()
	return m
}
