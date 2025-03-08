package list

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	helpkeys "github.com/evertonstz/go-workflows/components/keys"
	"github.com/evertonstz/go-workflows/components/persist"
	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
)

type MyItem struct {
	title, desc, command string
	dateAdded            time.Time
	dateUpdated          time.Time
}

var docStyle = lipgloss.NewStyle()

func (i MyItem) Title() string          { return i.title }
func (i MyItem) Description() string    { return i.desc }
func (i MyItem) Command() string        { return i.command }
func (i MyItem) DateAdded() time.Time   { return i.dateAdded }
func (i MyItem) DateUpdated() time.Time { return i.dateUpdated }
func (i MyItem) FilterValue() string    { return i.title }

type Model struct {
	inputs          inputsModel
	list            list.Model
	lastSelectedIdx int
}

func (m Model) CurentItem() MyItem {
	return m.list.SelectedItem().(MyItem)
}

func (m Model) CurrentItemIndex() int {
	return m.list.Index()
}

func (m Model) AllItems() []MyItem {
	var items []MyItem
	for _, i := range m.list.Items() {
		items = append(items, i.(MyItem))
	}
	return items
}

func (m Model) Init() tea.Cmd {
	return nil
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
		newItem := MyItem{
			title:       msg.Title,
			desc:        msg.Description,
			command:     msg.CommandText,
			dateAdded:   time.Now(),
			dateUpdated: time.Now(),
		}
		m.list.InsertItem(len(m.list.Items()), newItem)
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
			if item.(MyItem).title == msg.Item.Title {
				m.list.SetItem(i, MyItem{
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
			data = append(data, list.Item(MyItem{
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
		case key.Matches(msg, helpkeys.LisKeys.CopyWorkflow):
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				if selected, ok := selectedItem.(MyItem); ok {
					return m, shared.CopyToClipboardCmd(selected.command)
				}
			}
		}
	}

	var c tea.Cmd
	m.list, c = m.list.Update(msg)
	cmds = append(cmds, c)
	if m.list.Index() != m.lastSelectedIdx {
		cmds = m.setCurrentItemCmd(cmds)
		m.lastSelectedIdx = m.list.Index()
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}

func New() Model {
	m := Model{list: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0), inputs: newInputsModel()}
	m.list.SetShowTitle(false)
	m.list.SetShowHelp(false)
	m.Init()
	return m
}
