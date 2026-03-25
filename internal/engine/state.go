package engine

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

// CharState represents the typing state of a single character.
type CharState int

const (
	Untyped   CharState = iota
	Correct             // typed correctly
	Incorrect           // typed but wrong
)

// TypingState holds the full state of a typing exercise.
type TypingState struct {
	Target       []rune
	States       []CharState
	SyntaxColors []lipgloss.Color // per-rune syntax highlight color
	Cursor       int
	Errors       int
	StartedAt    time.Time
	FinishedAt   time.Time
	Started      bool
	Finished     bool
}

// NewTypingState creates a TypingState for the given code and syntax colors.
func NewTypingState(code string, colors []lipgloss.Color) TypingState {
	runes := []rune(code)
	return TypingState{
		Target:       runes,
		States:       make([]CharState, len(runes)),
		SyntaxColors: colors,
	}
}

// WPM returns net words per minute (correct chars / 5 / elapsed minutes).
func (s TypingState) WPM() int {
	if !s.Started {
		return 0
	}
	correct := s.CorrectCount()
	end := s.FinishedAt
	if !s.Finished {
		end = time.Now()
	}
	elapsed := end.Sub(s.StartedAt).Minutes()
	if elapsed <= 0 {
		return 0
	}
	return int(float64(correct) / 5.0 / elapsed)
}

// Accuracy returns the percentage of correctly typed characters so far.
func (s TypingState) Accuracy() float64 {
	if s.Cursor == 0 {
		return 100.0
	}
	return float64(s.CorrectCount()) / float64(s.Cursor) * 100.0
}

// CorrectCount returns the number of correctly typed characters.
func (s TypingState) CorrectCount() int {
	n := 0
	for _, st := range s.States {
		if st == Correct {
			n++
		}
	}
	return n
}

// ElapsedSeconds returns elapsed seconds since typing started.
func (s TypingState) ElapsedSeconds() int {
	if !s.Started {
		return 0
	}
	end := s.FinishedAt
	if !s.Finished {
		end = time.Now()
	}
	return int(end.Sub(s.StartedAt).Seconds())
}

// Progress returns chars typed / total chars (0.0–1.0).
func (s TypingState) Progress() float64 {
	if len(s.Target) == 0 {
		return 0
	}
	return float64(s.Cursor) / float64(len(s.Target))
}
