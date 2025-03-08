package commandlist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertonstz/go-workflows/components/list"
	textarea "github.com/evertonstz/go-workflows/components/text_area"
)

var (
	leftPanelWidthPercentage = 0.5
	leftPanelStyle           = lipgloss.NewStyle().
					AlignHorizontal(lipgloss.Left)
	rightPanelStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Left).
			PaddingTop(2).
			Width(15).
			Height(5)
)

type (
	panelsStyle struct {
	leftPanelStyle         lipgloss.Style
	rightPanelStyle        lipgloss.Style
}

 Model struct {
	list              list.Model
	textArea          textarea.Model
	panelsStyle	      panelsStyle
})
	
func (m Model) Init() tea.Cmd {
	return nil
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