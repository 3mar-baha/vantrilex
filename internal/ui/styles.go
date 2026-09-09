package ui

import "github.com/charmbracelet/lipgloss"

// Cosmic Cyber-Void palette.
var (
	Cyan    = lipgloss.Color("#00FFFF")
	Sky     = lipgloss.Color("#38BDF8")
	Navy    = lipgloss.Color("#030712")
	White   = lipgloss.Color("#FFFFFF")
	Violet  = lipgloss.Color("#8B5CF6")
	DimCyan = lipgloss.Color("#0E7490")
	Muted   = lipgloss.Color("#64748B")
	Green   = lipgloss.Color("#22C55E")
	Red     = lipgloss.Color("#EF4444")
	Amber   = lipgloss.Color("#F59E0B")
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	subStyle = lipgloss.NewStyle().
			Foreground(Sky)

	mutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	whiteStyle = lipgloss.NewStyle().
			Foreground(White)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Cyan).
			Padding(1, 2)

	okPill = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#03130A")).
		Background(Green).
		Bold(true).
		Padding(0, 1)

	missPill = lipgloss.NewStyle().
		Foreground(White).
		Background(Red).
		Bold(true).
		Padding(0, 1)

	optPill = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1A1000")).
		Background(Amber).
		Bold(true).
		Padding(0, 1)

	selStyle = lipgloss.NewStyle().
			Foreground(Navy).
			Background(Cyan).
			Bold(true).
			Padding(0, 1)

	rowStyle = lipgloss.NewStyle().
			Foreground(White)

	keyStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	haloStyle = lipgloss.NewStyle().
			Foreground(Cyan)
)

func pill(found, optional bool) string {
	switch {
	case found:
		return okPill.Render("OK")
	case optional:
		return optPill.Render("OPTIONAL")
	default:
		return missPill.Render("MISSING")
	}
}
