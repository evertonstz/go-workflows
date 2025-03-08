package commandlist

import (
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertonstz/go-workflows/components/list"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
)

var (
	leftPanelWidthPercentage = 0.5
	leftPanelStyle           = lipgloss.NewStyle().
	Background(lipgloss.Color("20")).
					AlignHorizontal(lipgloss.Left)
	rightPanelStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("20")).
			AlignHorizontal(lipgloss.Left).
			PaddingTop(2).
			Width(15).
			Height(5)
)

type panelsStyle struct {
	leftPanelStyle         lipgloss.Style
	rightPanelStyle        lipgloss.Style
}

type Model struct {
	list              list.Model
	textArea          textarea.Model
	panelsStyle	      panelsStyle
}

func (m Model) GetAllItems() []list.MyItem {
	return m.list.AllItems()
}

func (m *Model) SetSize(width, height int) {
	m.panelsStyle.leftPanelStyle = m.panelsStyle.leftPanelStyle.
		Width(int(math.Floor(float64(width) * leftPanelWidthPercentage))).
		Height(height)
	
	m.panelsStyle.rightPanelStyle = m.panelsStyle.rightPanelStyle.
		Width(width - m.panelsStyle.leftPanelStyle.GetWidth()).
		Height(height)

	leftWidthFrameSize, leftHeightFrameSize := m.panelsStyle.leftPanelStyle.GetFrameSize()
	rightWidthFrameSize, rightHeightFrameSize := m.panelsStyle.rightPanelStyle.GetFrameSize()

	leftPanelWidth := m.panelsStyle.leftPanelStyle.GetWidth() - leftWidthFrameSize
	rightPanelWidth := m.panelsStyle.rightPanelStyle.GetWidth() - rightWidthFrameSize

	m.list.SetSize(leftPanelWidth, m.panelsStyle.leftPanelStyle.GetHeight()-leftHeightFrameSize)
	m.textArea.SetSize(rightPanelWidth, m.panelsStyle.rightPanelStyle.GetHeight()-rightHeightFrameSize)
}
	
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	var listModel tea.Model
	listModel, cmd = m.list.Update(msg)
	m.list = listModel.(list.Model)
	cmds = append(cmds, cmd)
	taModel, cmd := m.textArea.Update(msg)
	m.textArea = taModel.(textarea.Model)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Bottom,
		m.panelsStyle.leftPanelStyle.Render(m.list.View()),
		m.panelsStyle.rightPanelStyle.Render(m.textArea.View()))
}

func New() Model {
	listModel := list.New()
	textAreaModel := textarea.New()

	return Model{
		list:     listModel,
		textArea: textAreaModel,
		panelsStyle: panelsStyle{
			leftPanelStyle:  leftPanelStyle,
			rightPanelStyle: rightPanelStyle,
		},
	}
}