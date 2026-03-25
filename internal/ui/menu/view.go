package menu

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/IFAKA/coding-type/internal/snippets"
	"github.com/IFAKA/coding-type/internal/theme"
)

const logo = `
  ██████╗ ██████╗ ██████╗ ███████╗
 ██╔════╝██╔═══██╗██╔══██╗██╔════╝
 ██║     ██║   ██║██║  ██║█████╗
 ██║     ██║   ██║██║  ██║██╔══╝
 ╚██████╗╚██████╔╝██████╔╝███████╗
  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝`

func (m Model) View() string {
	logoStyled := theme.Title.Render(logo)
	subtitle := theme.Muted.Render("  practice real interview code in your terminal")

	rows := []string{
		renderRow("language", m.langIdx, snippets.Languages, snippets.LangDisplay, m.activeRow == 0),
		renderRow("difficulty", m.diffIdx, snippets.Difficulties, nil, m.activeRow == 1),
		renderRow("mode", m.modeIdx, snippets.Modes, modeDisplay, m.activeRow == 2),
	}

	options := strings.Join(rows, "\n")

	help := renderHelp()

	content := strings.Join([]string{
		logoStyled,
		"",
		subtitle,
		"",
		"",
		options,
		"",
		theme.Separator.Render(strings.Repeat("─", 46)),
		"",
		help,
	}, "\n")

	box := theme.BoxBorder.
		Width(50).
		Padding(1, 3).
		Render(content)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, box)
}

func renderRow(label string, activeIdx int, options []string, display map[string]string, focused bool) string {
	labelStyle := theme.StatLabel
	if focused {
		labelStyle = theme.ActiveOption
	}

	var parts []string
	for i, opt := range options {
		name := opt
		if display != nil {
			if d, ok := display[opt]; ok {
				name = d
			}
		}
		if i == activeIdx {
			parts = append(parts, theme.SelectedOption.Render("[ "+name+" ]"))
		} else {
			parts = append(parts, theme.InactiveOption.Render(name))
		}
	}

	var arrow string
	if focused {
		arrow = theme.ActiveOption.Render(" ›")
	} else {
		arrow = "  "
	}

	return fmt.Sprintf("%s  %-12s  %s",
		arrow,
		labelStyle.Render(label),
		strings.Join(parts, theme.Muted.Render("  ·  ")))
}

func renderHelp() string {
	entries := []struct{ key, desc string }{
		{"enter", "start"},
		{"s", "stats"},
		{"←/→", "change"},
		{"↑/↓", "row"},
		{"q", "quit"},
	}
	var parts []string
	for _, e := range entries {
		parts = append(parts, theme.HelpKey.Render(e.key)+" "+theme.HelpDesc.Render(e.desc))
	}
	return "  " + strings.Join(parts, theme.Muted.Render("   "))
}

var modeDisplay = map[string]string{
	"practice": "practice",
	"timed":    "timed 60s",
}
