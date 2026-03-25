package typing

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/IFAKA/coding-type/internal/theme"
)

type finger int

const (
	lp finger = iota // left pinky
	lr               // left ring
	lm               // left middle
	li               // left index
	ri               // right index
	rm               // right middle
	rr               // right ring
	rp               // right pinky
)

var fingerColor = [8]lipgloss.Color{
	theme.Mauve,  // lp
	theme.Blue,   // lr
	theme.Sky,    // lm
	theme.Teal,   // li
	theme.Green,  // ri
	theme.Yellow, // rm
	theme.Peach,  // rr
	theme.Red,    // rp
}

type keyDef struct {
	ch      rune
	display rune // 0 means use ch
	f       finger
}

func (k keyDef) label() string {
	if k.display != 0 {
		return string(k.display)
	}
	return string(k.ch)
}

var (
	row0 = []keyDef{
		{'`', 0, lp}, {'1', 0, lp}, {'2', 0, lr}, {'3', 0, lm}, {'4', 0, li}, {'5', 0, li},
		{'6', 0, ri}, {'7', 0, ri}, {'8', 0, rm}, {'9', 0, rr}, {'0', 0, rp}, {'-', 0, rp}, {'=', 0, rp},
	}
	row1 = []keyDef{
		{'q', 0, lp}, {'w', 0, lr}, {'e', 0, lm}, {'r', 0, li}, {'t', 0, li},
		{'y', 0, ri}, {'u', 0, ri}, {'i', 0, rm}, {'o', 0, rr}, {'p', 0, rp},
		{'[', 0, rp}, {']', 0, rp}, {'\\', 0, rp},
	}
	row2 = []keyDef{
		{'a', 0, lp}, {'s', 0, lr}, {'d', 0, lm}, {'f', 0, li}, {'g', 0, li},
		{'h', 0, ri}, {'j', 0, ri}, {'k', 0, rm}, {'l', 0, rr}, {';', 0, rp}, {'\'', 0, rp},
		{'\n', '↵', rp},
	}
	row3 = []keyDef{
		{'z', 0, lp}, {'x', 0, lr}, {'c', 0, lm}, {'v', 0, li}, {'b', 0, li},
		{'n', 0, ri}, {'m', 0, ri}, {',', 0, rm}, {'.', 0, rr}, {'/', 0, rp},
	}

	kbRows    = [][]keyDef{row0, row1, row2, row3}
	gapAfter  = []int{5, 4, 4, 4}
	rowIndent = []string{"", " ", "  ", ""}
)

var shiftMap = map[rune]rune{
	'!': '1', '@': '2', '#': '3', '$': '4', '%': '5',
	'^': '6', '&': '7', '*': '8', '(': '9', ')': '0',
	'_': '-', '+': '=',
	'{': '[', '}': ']', '|': '\\',
	':': ';', '"': '\'',
	'<': ',', '>': '.', '?': '/',
	'~': '`',
}

func resolveKey(ch rune) (base rune, needsShift bool) {
	switch ch {
	case 0, '\t':
		return 0, false
	case '\n':
		return '\n', false
	case ' ':
		return ' ', false
	}
	if unicode.IsUpper(ch) {
		return unicode.ToLower(ch), true
	}
	if b, ok := shiftMap[ch]; ok {
		return b, true
	}
	return ch, false
}

// activeFinger returns the finger responsible for the given base key, or -1.
func activeFinger(base rune) finger {
	for _, row := range kbRows {
		for _, k := range row {
			if k.ch == base {
				return k.f
			}
		}
	}
	return -1
}

// fingerLabels defines the 10 finger columns in left→right order for the label strip.
// The strip is designed to align with the home row (row2, 2-space indent).
var fingerLabels = []struct {
	label string
	f     finger
}{
	{"P", lp}, {"R", lr}, {"M", lm}, {"I", li}, {"I", li},
	{"I", ri}, {"I", ri}, {"M", rm}, {"R", rr}, {"P", rp},
}

func renderKeyboard(currentChar rune) string {
	base, needsShift := resolveKey(currentChar)
	af := activeFinger(base) // finger to highlight (-1 if none)

	activeStyle := lipgloss.NewStyle().Background(theme.Mauve).Foreground(theme.Base).Bold(true)
	shiftOnStyle := lipgloss.NewStyle().Background(theme.Teal).Foreground(theme.Base).Bold(true)
	shiftOffStyle := lipgloss.NewStyle().Foreground(theme.Surface1)
	spaceOnStyle := lipgloss.NewStyle().Background(theme.Teal).Foreground(theme.Base).Bold(true)
	spaceOffStyle := lipgloss.NewStyle().Foreground(theme.Surface1)

	renderKey := func(k keyDef) string {
		label := k.label()
		isHit := base != 0 && k.ch == base
		isSameFinger := af >= 0 && k.f == af && !isHit

		switch {
		case isHit:
			return activeStyle.Render(label)
		case isSameFinger:
			return lipgloss.NewStyle().Foreground(fingerColor[k.f]).Render(label)
		default:
			return lipgloss.NewStyle().Foreground(fingerColor[k.f]).Faint(true).Render(label)
		}
	}

	var lines []string

	// Hand label row — "LEFT" left-aligned, "RIGHT" right-aligned over their halves
	mutedStyle := lipgloss.NewStyle().Foreground(theme.Surface1)
	lines = append(lines, "  "+mutedStyle.Render("LEFT")+"               "+mutedStyle.Render("RIGHT"))

	// Finger label strip — aligned with home row (2-space indent, same gap position)
	var labelSB strings.Builder
	labelSB.WriteString("  ")
	for i, fl := range fingerLabels {
		if i == 5 { // hand gap — matches gapAfter[2]+1
			labelSB.WriteString(" ")
		}
		isActive := af >= 0 && fl.f == af
		if isActive {
			labelSB.WriteString(lipgloss.NewStyle().Foreground(fingerColor[fl.f]).Bold(true).Render(fl.label))
		} else {
			labelSB.WriteString(lipgloss.NewStyle().Foreground(fingerColor[fl.f]).Faint(true).Render(fl.label))
		}
		if i < len(fingerLabels)-1 {
			labelSB.WriteString(" ")
		}
	}
	lines = append(lines, labelSB.String())

	for rowIdx, row := range kbRows {
		var sb strings.Builder

		if rowIdx == 3 {
			sb.WriteString("   ")
			if needsShift {
				sb.WriteString(shiftOnStyle.Render("⇧"))
			} else {
				sb.WriteString(shiftOffStyle.Render("⇧"))
			}
			sb.WriteString(" ")
		} else {
			sb.WriteString(rowIndent[rowIdx])
		}

		for i, k := range row {
			if i == gapAfter[rowIdx]+1 {
				sb.WriteString(" ")
			}
			sb.WriteString(renderKey(k))
			sb.WriteString(" ")
		}

		if rowIdx == 3 {
			if needsShift {
				sb.WriteString(shiftOnStyle.Render("⇧"))
			} else {
				sb.WriteString(shiftOffStyle.Render("⇧"))
			}
		}

		lines = append(lines, sb.String())
	}

	// Space bar
	const spaceBarLabel = "___________"
	var spaceSB strings.Builder
	spaceSB.WriteString("        ")
	if base == ' ' {
		spaceSB.WriteString(spaceOnStyle.Render(spaceBarLabel))
	} else {
		spaceSB.WriteString(spaceOffStyle.Render(spaceBarLabel))
	}
	lines = append(lines, spaceSB.String())

	return strings.Join(lines, "\n")
}
