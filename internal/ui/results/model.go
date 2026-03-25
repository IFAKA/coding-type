package results

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/IFAKA/coding-type/internal/ui/msgs"
)

// Model is the BubbleTea model for the results screen.
type Model struct {
	done   msgs.TypingDoneMsg
	width  int
	height int
}

// New creates a results model from a TypingDoneMsg.
func New(done msgs.TypingDoneMsg, width, height int) Model {
	return Model{done: done, width: width, height: height}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "r":
			return m, func() tea.Msg {
				return msgs.RetryMsg{
					Snippet: m.done.Snippet,
					Config:  m.done.Config,
					BestWPM: m.currentBest(),
					AvgWPM:  m.done.WPM,
				}
			}

		case "n":
			return m, func() tea.Msg {
				return msgs.NextSnippetMsg{
					Config:  m.done.Config,
					BestWPM: m.currentBest(),
					AvgWPM:  m.done.WPM,
				}
			}

		case "m", "esc":
			return m, func() tea.Msg { return msgs.NavigateMsg{To: msgs.ScreenMenu} }
		}
	}
	return m, nil
}

// currentBest returns the current personal best WPM after this session.
func (m Model) currentBest() int {
	if m.done.IsPersonalBest {
		return m.done.WPM
	}
	// Approximate: prior best was avgWPM + (WPM - avgWPM) normalized
	// Just return 0 to let app.go recalculate from history
	return 0
}

// Done returns the underlying TypingDoneMsg.
func (m Model) Done() msgs.TypingDoneMsg { return m.done }
func (m Model) Width() int               { return m.width }
func (m Model) Height() int              { return m.height }
