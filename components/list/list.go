package list

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
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

var docStyle = lipgloss.NewStyle()

func (i myItem) Title() string          { return i.title }
func (i myItem) Description() string    { return i.desc }
func (i myItem) Command() string        { return i.command }
func (i myItem) DateAdded() time.Time   { return i.dateAdded }
func (i myItem) DateUpdated() time.Time { return i.dateUpdated }
func (i myItem) FilterValue() string    { return i.title }

type Model struct {
	state           inputs
	inputs          inputsModel
	list            list.Model
	lastSelectedIdx int
}

func (m Model) InputOn() bool {
	return m.state == addNewOn
}

func (m Model) CurentItem() myItem {
	return m.list.SelectedItem().(myItem)
}

func (m Model) CurrentItemIndex() int {
	return m.list.Index()
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
	h, v := docStyle.GetFrameSize()
	m.list.SetSize(width-h, height-v)
}

func (m Model) setCurrentItemCmd(cmds []tea.Cmd) []tea.Cmd {
	cmds = append(cmds, shared.SetCurrentItemCmd(models.Item{
		Title:       m.CurentItem().Title(),
		Desc:        m.CurentItem().Description(),
		Command:     m.CurentItem().Command(),
		DateAdded:   m.CurentItem().DateAdded(),
		DateUpdated: m.CurentItem().DateUpdated()},
	))
	return cmds
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case shared.DidAddNewItemMsg:
		// if m.inputs.Title.Value() == "" || m.inputs.Description.Value() == "" {
		// 	return m, nil
		// }
		newItem := myItem{
			title:       msg.Title,
			desc:        msg.Description,
			command:     msg.CommandText,
			dateAdded:   time.Now(),
			dateUpdated: time.Now(),
		}
		m.list.InsertItem(len(m.list.Items()), newItem)
		// m.hideAddNew()
		return m, nil
	case shared.DidDeleteItemMsg:
		m.list.RemoveItem(msg.Index) 
		if m.list.Index() >= len(m.list.Items()) {
			newIndex := len(m.list.Items()) - 1
			if newIndex < 0 {
				newIndex = 0
			}
			m.list.Select(newIndex)
		}

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
		cmds = m.setCurrentItemCmd(cmds)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Esc):
			if m.state == addNewOn {
				m.hideAddNew()
				return m, nil
			}
		// case key.Matches(msg, Keys.AddNewWorkflow):
		// 	if m.state == addNewOff {
		// 		m.showAddNew()
		// 		return m, nil
		// 	}
		case key.Matches(msg, Keys.CopyWorkflow):
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				if selected, ok := selectedItem.(myItem); ok {
					return m, shared.CopyToClipboardCmd(selected.command)
				}
			}
		case key.Matches(msg, Keys.Enter):
			if m.state == addNewOff {
				return m, shared.SetCurrentItemCmd(models.Item{
					Title:       m.CurentItem().Title(),
					Desc:        m.CurentItem().Description(),
					Command:     m.CurentItem().Command(),
					DateAdded:   m.CurentItem().DateAdded(),
					DateUpdated: m.CurentItem().DateUpdated()},
				)
			}
			var c tea.Cmd
			m.inputs, c = m.inputs.Update(msg)
			return m, c
		}
	}

	switch m.state {
	case addNewOff:
		var c tea.Cmd
		m.list, c = m.list.Update(msg)
		cmds = append(cmds, c)
		if m.list.Index() != m.lastSelectedIdx {
			cmds = m.setCurrentItemCmd(cmds)
			m.lastSelectedIdx = m.list.Index()
		}
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
	inputView := m.inputs.View()
	inputViewHeight := strings.Count(inputView, "\n")
	switch m.state {
	case addNewOff:
		v = docStyle.Render(m.list.View())
		// return  m.list.View()
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
	m := Model{list: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0), inputs: newInputsModel(), state: addNewOff}
	m.list.SetShowTitle(false)
	m.list.SetShowHelp(false)
	m.Init()
	return m
}
