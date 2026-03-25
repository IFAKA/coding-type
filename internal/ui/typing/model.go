package typing

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/IFAKA/coding-type/internal/engine"
	"github.com/IFAKA/coding-type/internal/snippets"
	"github.com/IFAKA/coding-type/internal/sound"
	"github.com/IFAKA/coding-type/internal/ui/msgs"
)

const timedDuration = 60 * time.Second

type tickMsg struct{}

// Model is the BubbleTea model for the typing exercise screen.
type Model struct {
	state         engine.TypingState
	snippet       snippets.Snippet
	config        snippets.Config
	bestWPM       int
	avgWPM        int
	width         int
	height        int
	cursorVisible bool
	errorFlash    int // counts down from 4 → 0; cursor shows red while > 0
	tickCount     int
}

// New creates a typing model from a StartTypingMsg.
func New(msg msgs.StartTypingMsg, width, height int) Model {
	chromaLang := snippets.ChromaLang[msg.Config.Language]
	colors := engine.SyntaxColors(msg.Snippet.Code, chromaLang)
	state := engine.NewTypingState(msg.Snippet.Code, colors)
	return Model{
		state:         state,
		snippet:       msg.Snippet,
		config:        msg.Config,
		bestWPM:       msg.BestWPM,
		avgWPM:        msg.AvgWPM,
		width:         width,
		height:        height,
		cursorVisible: true,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tickMsg:
		m.tickCount++
		// Blink cursor every 5 ticks (500ms)
		if m.tickCount%5 == 0 {
			m.cursorVisible = !m.cursorVisible
		}
		// Decay error flash
		if m.errorFlash > 0 {
			m.errorFlash--
		}
		if m.state.Finished {
			return m, nil
		}
		// Timed mode: force finish when time runs out
		if m.config.Mode == "timed" && m.state.Started {
			if time.Since(m.state.StartedAt) >= timedDuration {
				m.state = engine.ForceFinish(m.state)
				return m, m.doneCmd()
			}
		}
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return msgs.NavigateMsg{To: msgs.ScreenMenu} }
		case "ctrl+r":
			return m, m.retryCmd()
		default:
			if m.state.Finished {
				return m, nil
			}
			prevErrors := m.state.Errors
			prevCursor := m.state.Cursor
			var done bool
			m.state, done = engine.ProcessKey(m.state, msg)
			if done {
				isPersonalBest := m.bestWPM == 0 || m.state.WPM() > m.bestWPM
				if isPersonalBest {
					sound.PlayPersonalBest()
				} else {
					sound.PlayComplete()
				}
				return m, m.doneCmd()
			}
			if m.state.Cursor > prevCursor {
				if m.state.Errors > prevErrors {
					m.errorFlash = 4
					sound.PlayError()
				} else if msg.Type == tea.KeyEnter {
					sound.PlayNewline()
				} else {
					sound.PlayCorrect()
				}
			}
		}
	}
	return m, nil
}

func (m Model) doneCmd() tea.Cmd {
	return func() tea.Msg {
		isPersonalBest := m.bestWPM == 0 || m.state.WPM() > m.bestWPM
		diff := m.state.WPM() - m.avgWPM
		return msgs.TypingDoneMsg{
			Snippet:        m.snippet,
			Config:         m.config,
			WPM:            m.state.WPM(),
			Accuracy:       m.state.Accuracy(),
			Errors:         m.state.Errors,
			Duration:       m.state.FinishedAt.Sub(m.state.StartedAt),
			IsPersonalBest: isPersonalBest,
			DiffFromAvg:    diff,
		}
	}
}

func (m Model) retryCmd() tea.Cmd {
	return func() tea.Msg {
		return msgs.RetryMsg{
			Snippet: m.snippet,
			Config:  m.config,
			BestWPM: m.bestWPM,
			AvgWPM:  m.avgWPM,
		}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

// Accessors for the view
func (m Model) State() engine.TypingState  { return m.state }
func (m Model) Snippet() snippets.Snippet  { return m.snippet }
func (m Model) Config() snippets.Config    { return m.config }
func (m Model) Width() int                 { return m.width }
func (m Model) Height() int                { return m.height }
