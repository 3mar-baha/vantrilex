package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"vantrilex/internal/particles"
)

// FX overlay palette (mirrors the particle age ramp, tuned for overlay).
var (
	fxWhite  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	fxCyan   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)
	fxViolet = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B5CF6"))
	fxDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8"))
)

// snapStyle maps an overlay snapshot to its glow style by burst age.
func snapStyle(s particles.Snapshot) lipgloss.Style {
	if s.Burst {
		switch {
		case s.Age < 0.25:
			return fxWhite
		case s.Age < 0.6:
			return fxCyan
		default:
			return fxViolet
		}
	}
	return fxDim
}

// spliceCell replaces the visible cell at column x with a styled glyph,
// preserving ANSI escapes and wide-rune alignment. Rows shorter than x
// are space-padded (never beyond the terminal width, enforced by caller).
func spliceCell(row string, x int, styled string) string {
	newW := lipgloss.Width(styled)
	var b strings.Builder
	vx := 0
	done := false
	rs := []rune(row)
	for i := 0; i < len(rs); {
		r := rs[i]
		if r == '\x1b' {
			j := i + 1
			if j < len(rs) && rs[j] == '[' {
				j++
				for j < len(rs) && !(rs[j] >= '@' && rs[j] <= '~') {
					j++
				}
				if j < len(rs) {
					j++
				}
			}
			b.WriteString(string(rs[i:j]))
			i = j
			continue
		}
		w := lipgloss.Width(string(r))
		if !done && vx == x {
			b.WriteString(styled)
			for k := newW; k < w; k++ {
				b.WriteByte(' ')
			}
			done = true
			i++
			vx += w
			continue
		}
		b.WriteRune(r)
		i++
		vx += w
	}
	if !done {
		for vx < x {
			b.WriteByte(' ')
			vx++
		}
		b.WriteString(styled)
	}
	return b.String()
}

// overlayFX composites transient mouse FX snapshots (trail, click bursts,
// perimeter waves) over a finished menu frame. Row count and row widths
// never change, so layout budgets and tests are unaffected.
func overlayFX(base string, snaps []particles.Snapshot, w, h int) string {
	if len(snaps) == 0 || w <= 0 || h <= 0 {
		return base
	}
	want := h - 1
	lines := strings.Split(base, "\n")
	for len(lines) < want {
		lines = append(lines, "")
	}
	if len(lines) > want {
		lines = lines[:want]
	}
	occ := map[[2]int]particles.Snapshot{}
	for _, s := range snaps { // newest-first: first occupant wins
		if s.X < 0 || s.X >= w || s.Y < 0 || s.Y >= want {
			continue
		}
		k := [2]int{s.X, s.Y}
		if _, ok := occ[k]; !ok {
			occ[k] = s
		}
	}
	for k, s := range occ {
		lines[k[1]] = spliceCell(lines[k[1]], k[0], snapStyle(s).Render(string(s.Glyph)))
	}
	return strings.Join(lines, "\n")
}
