package theme

import "github.com/charmbracelet/lipgloss"

// Catppuccin Mocha palette
const (
	ColorBase     = "#1E1E2E"
	ColorText     = "#CDD6F4"
	ColorSubtext  = "#A6ADC8"
	ColorOverlay  = "#585B70"
	ColorSurface0 = "#313244"
	ColorSurface1 = "#45475A"
	ColorBlue     = "#89B4FA"
	ColorGreen    = "#A6E3A1"
	ColorRed      = "#F38BA8"
	ColorMauve    = "#CBA6F7"
	ColorPeach    = "#FAB387"
	ColorYellow   = "#F9E2AF"
	ColorSky      = "#89DCEB"
	ColorGray     = "#6C7086"
	ColorTeal     = "#94E2D5"
)

var (
	Base     = lipgloss.Color(ColorBase)
	Text     = lipgloss.Color(ColorText)
	Subtext  = lipgloss.Color(ColorSubtext)
	Overlay  = lipgloss.Color(ColorOverlay)
	Surface0 = lipgloss.Color(ColorSurface0)
	Surface1 = lipgloss.Color(ColorSurface1)
	Blue     = lipgloss.Color(ColorBlue)
	Green    = lipgloss.Color(ColorGreen)
	Red      = lipgloss.Color(ColorRed)
	Mauve    = lipgloss.Color(ColorMauve)
	Peach    = lipgloss.Color(ColorPeach)
	Yellow   = lipgloss.Color(ColorYellow)
	Sky      = lipgloss.Color(ColorSky)
	Gray     = lipgloss.Color(ColorGray)
	Teal     = lipgloss.Color(ColorTeal)
)

var (
	BoxBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Surface0)

	Title = lipgloss.NewStyle().
		Foreground(Blue).
		Bold(true)

	Accent = lipgloss.NewStyle().
		Foreground(Blue)

	Muted = lipgloss.NewStyle().
		Foreground(Overlay)

	ActiveOption = lipgloss.NewStyle().
			Foreground(Mauve).
			Bold(true)

	InactiveOption = lipgloss.NewStyle().
			Foreground(Overlay)

	SelectedOption = lipgloss.NewStyle().
			Foreground(Blue).
			Bold(true)

	StatLabel = lipgloss.NewStyle().
			Foreground(Gray)

	StatValue = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	CorrectChar = lipgloss.NewStyle().
			Foreground(Green)

	IncorrectChar = lipgloss.NewStyle().
			Background(Red).
			Foreground(Base)

	CursorChar = lipgloss.NewStyle().
			Background(Mauve).
			Foreground(Base)

	UntypedChar = lipgloss.NewStyle().
			Foreground(Overlay)

	HelpKey = lipgloss.NewStyle().
		Foreground(Blue)

	HelpDesc = lipgloss.NewStyle().
			Foreground(Gray)

	Success = lipgloss.NewStyle().
		Foreground(Green).
		Bold(true)

	PersonalBest = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	Separator = lipgloss.NewStyle().
			Foreground(Surface1)

	HeaderBadge = lipgloss.NewStyle().
			Foreground(Base).
			Background(Blue).
			Padding(0, 1)

	DiffEasy = lipgloss.NewStyle().
			Foreground(Green)

	DiffMedium = lipgloss.NewStyle().
			Foreground(Yellow)

	DiffHard = lipgloss.NewStyle().
			Foreground(Red)
)
