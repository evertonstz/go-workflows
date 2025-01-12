package textarea

import (
	"regexp"
	"strings"

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
		{`\b(cd|ls|echo|cat|grep|awk|sed|export|sudo|mkdir)\b`, lipgloss.NewStyle().Foreground(lipgloss.Color("34"))},
		// Flags (e.g., -l, --help)
		{`\s(-{1,2}\w+)`, lipgloss.NewStyle().Foreground(lipgloss.Color("32"))},
		// Strings (e.g., "text" or 'text')
		{`"[^"]*"|'[^']*'`, lipgloss.NewStyle().Foreground(lipgloss.Color("33"))},
		// Environment variables (e.g., $HOME, $PATH)
		{`\$[a-zA-Z_][a-zA-Z0-9_]*`, lipgloss.NewStyle().Foreground(lipgloss.Color("35"))},
		// Numbers
		{`\b\d+\b`, lipgloss.NewStyle().Foreground(lipgloss.Color("36"))},
		// Operators (e.g., &&, ||, >, <, >>, ;)
		{`(\|\||&&|;|>|>>|<)`, lipgloss.NewStyle().Foreground(lipgloss.Color("31"))},
		// Special variables like "$terminfo[kcud1]"
		{`\$\[[^\]]*\]`, lipgloss.NewStyle().Foreground(lipgloss.Color("37"))},
	}

	for _, rule := range styles {
		re := regexp.MustCompile(rule.pattern)
		command = re.ReplaceAllStringFunc(command, func(match string) string {
			return rule.style.Render(match)
		})
	}

	command = highlightFirstNonFlagWords(command)

	return command
}

func highlightFirstNonFlagWords(command string) string {
	re := regexp.MustCompile(`(&&|\|\||;)`)
	parts := re.Split(command, -1)
	delimiters := re.FindAllString(command, -1)

	for i, part := range parts {
		words := regexp.MustCompile(`\S+`).FindAllString(part, -1)
		if len(words) == 0 {
			continue
		}

		for _, word := range words {
			if isKeyword(word) || isFlag(word) {
				continue
			}

			highlighted := lipgloss.NewStyle().Foreground(lipgloss.Color("32")).Render(word)
			parts[i] = strings.Replace(part, word, highlighted, 1)
			break
		}
	}

	var result strings.Builder
	for i, part := range parts {
		result.WriteString(part)
		if i < len(delimiters) {
			result.WriteString(delimiters[i])
		}
	}
	return result.String()
}

func isKeyword(word string) bool {
	keywords := []string{"cd", "ls", "echo", "cat", "grep", "awk", "sed", "export", "sudo", "mkdir"}
	for _, keyword := range keywords {
		if word == keyword {
			return true
		}
	}
	return false
}

func isFlag(word string) bool {
	return regexp.MustCompile(`^-{1,2}\w+`).MatchString(word)
}
