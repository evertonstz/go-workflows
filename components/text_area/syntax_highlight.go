package textarea

import (
	"regexp"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func SyntaxHighlight(command string) string {
	lipgloss.SetColorProfile(termenv.ANSI)

	styles := []struct {
		pattern string
		style   lipgloss.Style
	}{
		// Keywords (e.g., shell builtins or commands)
		{`\b(cd|ls|echo|cat|grep|awk|sed|export|sudo)\b`, lipgloss.NewStyle().Foreground(lipgloss.Color("34"))},
		// Flags (e.g., -l, --help)
		{`\s(-{1,2}\w+)`, lipgloss.NewStyle().Foreground(lipgloss.Color("32"))},
		// Strings (e.g., "text" or 'text')
		{`"[^"]*"|'[^']*'`, lipgloss.NewStyle().Foreground(lipgloss.Color("33"))},
		// Environment variables (e.g., $HOME, $PATH)
		{`\$[a-zA-Z_][a-zA-Z0-9_]*`, lipgloss.NewStyle().Foreground(lipgloss.Color("35"))},
		// Numbers
		{`\b\d+\b`, lipgloss.NewStyle().Foreground(lipgloss.Color("36"))},
	}

	for _, rule := range styles {
		re := regexp.MustCompile(rule.pattern)
		command = re.ReplaceAllStringFunc(command, func(match string) string {
			return rule.style.Render(match)
		})
	}

	return command
}
