package ui

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"vantrilex/internal/doctor"
)

func toolkitStatusSafe() (string, map[string]bool, error) {
	return doctor.ToolkitStatus()
}

var ansiRe = regexp.MustCompile("\x1b\\[[0-9;]*m")

// stripANSI removes SGR escape sequences for measurement and safe truncation.
func stripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

// countLines counts newline-separated rows.
func countLines(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// fitLines truncates a block to at most maxH rows (ANSI-safe: lines are never
// split mid-row, only whole rows are dropped). Trailing empty rows are
// trimmed so padded boxes never grow a phantom extra line.
func fitLines(s string, maxH int) string {
	if maxH <= 0 {
		return ""
	}
	lines := strings.Split(s, "\n")
	// Drop one trailing artifact row from blocks ending in "\n". Only truly
	// empty rows are dropped — whitespace-only rows (e.g. particle field
	// spacing) are significant for height and must be preserved.
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) > maxH {
		lines = lines[:maxH]
	}
	return strings.Join(lines, "\n")
}

// fitWidthLines truncates every row to at most maxW visible cells. Rows that
// fit keep their styling; overflowing rows fall back to plain text so no
// dangling escape sequence can bleed across the screen.
func fitWidthLines(s string, maxW int) string {
	if maxW <= 0 {
		return ""
	}
	lines := strings.Split(s, "\n")
	for i, ln := range lines {
		if lipgloss.Width(ln) > maxW {
			plain := []rune(stripANSI(ln))
			if len(plain) > maxW {
				plain = append(plain[:maxW-1], '…')
			}
			lines[i] = string(plain)
		}
	}
	return strings.Join(lines, "\n")
}
