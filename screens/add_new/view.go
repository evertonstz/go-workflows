package addnew

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	switch m.selectedInput {
	case title:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.focusedInput.Render(m.Title.View()),
			m.styles.blurredInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			lipgloss.NewStyle().Align(lipgloss.Right).Width(m.Title.Width).Render(
				lipgloss.JoinHorizontal(lipgloss.Top, m.styles.blurredButton, m.styles.blurredCloseButton))))
	case description:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.blurredInput.Render(m.Title.View()),
			m.styles.focusedInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			lipgloss.NewStyle().Align(lipgloss.Right).Width(m.Title.Width).Render(
				lipgloss.JoinHorizontal(lipgloss.Top, m.styles.blurredButton, m.styles.blurredCloseButton))))
	case submit:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.blurredInput.Render(m.Title.View()),
			m.styles.blurredInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			lipgloss.NewStyle().Align(lipgloss.Right).Width(m.Title.Width).Render(
				lipgloss.JoinHorizontal(lipgloss.Top, m.styles.focusedButton, m.styles.blurredCloseButton))))
	case close:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.blurredInput.Render(m.Title.View()),
			m.styles.blurredInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			lipgloss.NewStyle().Align(lipgloss.Right).Width(m.Title.Width).Render(
				lipgloss.JoinHorizontal(lipgloss.Top, m.styles.blurredButton, m.styles.focusedCloseButton))))
	default:
		return mainStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
			m.styles.blurredInput.Render(m.Title.View()),
			m.styles.blurredInput.Render(m.Description.View()),
			m.styles.blurredTextArea.Render(m.TextArea.View()),
			lipgloss.NewStyle().Align(lipgloss.Right).Width(m.Title.Width).Render(
				lipgloss.JoinHorizontal(lipgloss.Top, m.styles.blurredButton, m.styles.blurredCloseButton))))
	}
}
