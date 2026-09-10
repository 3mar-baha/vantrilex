package ui

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Cinematic intro timeline (seconds): 4 acts over 10 seconds.
const (
	introDuration = 10.0
	act2At        = 2.5 // Act 2: orbital singularity
	act3At        = 5.0 // Act 3: implosion + laser trace
	act4At        = 7.5 // Act 4: supernova ignition
	act4TeaseAt   = 9.2 // Act 4 tail: dock header, tease the Doctor card
)

// introPhase maps elapsed seconds to an act index:
// 0 deep void, 1 orbital singularity, 2 laser trace, 3 ignition, 4 finished.
func introPhase(elapsed float64) int {
	switch {
	case elapsed < act2At:
		return 0
	case elapsed < act3At:
		return 1
	case elapsed < act4At:
		return 2
	case elapsed < introDuration:
		return 3
	default:
		return 4
	}
}

func introPhaseName(ph int) string {
	switch ph {
	case 0:
		return "DEEP VOID — TELEMETRY LOCK"
	case 1:
		return "ORBITAL SINGULARITY"
	case 2:
		return "COLLAPSE — LASER TRACE"
	default:
		return "SUPERNOVA IGNITION"
	}
}

// warpStar is one hyperspace star flying past the camera.
type warpStar struct {
	x, y float64 // -1..1 across the view
	z    float64 // depth: 0 near -> 1 far
	hue  int     // 0 white, 1 cyan, 2 violet, 3 dim
}

func newWarpField(n int, rng *rand.Rand) []warpStar {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	stars := make([]warpStar, 0, n)
	for i := 0; i < n; i++ {
		stars = append(stars, warpStar{
			x:   rng.Float64()*2 - 1,
			y:   rng.Float64()*2 - 1,
			z:   rng.Float64(),
			hue: rng.Intn(4),
		})
	}
	return stars
}

// advanceWarp moves stars toward the camera; speed ramps up for the
// hyper-jump acceleration feel. dt is seconds.
func advanceWarp(stars []warpStar, dt, accel float64, rng *rand.Rand) {
	for i := range stars {
		s := &stars[i]
		s.z -= dt * (0.25 + accel) * (1.1 - s.z)
		if s.z <= 0.03 {
			s.x = rng.Float64()*2 - 1
			s.y = rng.Float64()*2 - 1
			s.z = 0.9 + rng.Float64()*0.1
			s.hue = rng.Intn(4)
		}
	}
}

var (
	warpWhite  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	warpCyan   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	warpViolet = lipgloss.NewStyle().Foreground(lipgloss.Color("#6366F1"))
	warpDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#1E3A5F"))
	introTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)
	introSub   = lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8"))
	introHint  = lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B"))
	flashWhite = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	// Act palettes: cobalt -> ice blue -> electric violet.
	actCobalt = lipgloss.NewStyle().Foreground(lipgloss.Color("#1E1B4B"))
	actIce    = lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8"))
	actLaser  = lipgloss.NewStyle().Foreground(lipgloss.Color("#7DD3FC")).Bold(true)
	lockGreen = lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E")).Bold(true)
	scanAmber = lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B")).Bold(true)
)

func warpGlyph(z float64) string {
	switch {
	case z < 0.25:
		return "✦"
	case z < 0.5:
		return "✧"
	case z < 0.75:
		return "⋆"
	default:
		return "·"
	}
}

func warpStyle(hue int, near bool) lipgloss.Style {
	if near {
		return warpWhite
	}
	switch hue {
	case 1:
		return warpCyan
	case 2:
		return warpViolet
	case 3:
		return warpDim
	default:
		return warpWhite
	}
}

// introCell is one composited screen cell. Multi-cell chunks (emblem rows,
// titles) occupy one cell slot and mark the cells they visually cover as
// skip so flattened rows are always exactly w cells wide (no terminal wrap).
type introCell struct {
	text string // already styled
	set  bool
	skip bool
	// chunk marks a multi-cell chunk origin; at most one chunk owns a row.
	chunk bool
}

func (m *Model) viewIntro() string {
	w, h := m.width, m.height
	if w < 20 {
		w = 20
	}
	if h < 10 {
		h = 10
	}
	elapsed := time.Since(m.introStart).Seconds()
	ph := introPhase(elapsed)

	grid := make([][]introCell, h)
	for y := range grid {
		grid[y] = make([]introCell, w)