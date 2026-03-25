package engine

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ProcessKey handles a single keypress and returns the updated state and
// whether the exercise is now complete.
func ProcessKey(s TypingState, key tea.KeyMsg) (TypingState, bool) {
	if s.Finished {
		return s, true
	}

	// Backspace: revert last character, keep error count
	if key.Type == tea.KeyBackspace || key.Type == tea.KeyDelete {
		if s.Cursor > 0 {
			s.Cursor--
			s.States[s.Cursor] = Untyped
		}
		return s, false
	}

	if s.Cursor >= len(s.Target) {
		return finish(s), true
	}

	typedRune := extractRune(key)
	if typedRune == 0 {
		return s, false
	}

	// Start timer on first meaningful keypress
	if !s.Started {
		s.Started = true
		s.StartedAt = time.Now()
	}

	expected := s.Target[s.Cursor]
	if typedRune == expected {
		s.States[s.Cursor] = Correct
	} else {
		s.States[s.Cursor] = Incorrect
		s.Errors++
	}
	s.Cursor++

	// Auto-indent: after a correct newline, skip leading whitespace automatically
	if typedRune == '\n' && typedRune == expected {
		for s.Cursor < len(s.Target) && (s.Target[s.Cursor] == ' ' || s.Target[s.Cursor] == '\t') {
			s.States[s.Cursor] = Correct
			s.Cursor++
		}
	}

	if s.Cursor >= len(s.Target) {
		return finish(s), true
	}
	return s, false
}

// ForceFinish marks the exercise as finished (used for timed mode timeout).
func ForceFinish(s TypingState) TypingState {
	return finish(s)
}

func finish(s TypingState) TypingState {
	if !s.Finished {
		s.Finished = true
		s.FinishedAt = time.Now()
		if !s.Started {
			s.StartedAt = s.FinishedAt
			s.Started = true
		}
	}
	return s
}

// extractRune extracts the typed rune from a key message.
// Returns 0 for keys we don't handle (ctrl combos, arrows, etc.).
func extractRune(key tea.KeyMsg) rune {
	switch key.Type {
	case tea.KeyEnter:
		return '\n'
	case tea.KeySpace:
		return ' '
	case tea.KeyTab:
		return '\t'
	case tea.KeyRunes:
		if len(key.Runes) > 0 {
			return key.Runes[0]
		}
	}
	return 0
}
