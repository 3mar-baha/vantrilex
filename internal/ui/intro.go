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
	}
	plot := func(x, y int, text string) {
		if x < 0 || x >= w || y < 0 || y >= h {
			return
		}
		if grid[y][x].skip || grid[y][x].chunk {
			return // never break a multi-cell chunk apart
		}
		grid[y][x] = introCell{text: text, set: true}
	}
	cx, cy := float64(w)/2, float64(h)/2

	// Layer 1: warp field. Gentle drift in acts 1-3; hyperspace streaks in
	// act 4 (the tick loop drives accel to 8.0 there).
	for _, s := range m.warp {
		scale := (1 - s.z + 0.08)
		sx := int(math.Round(cx + s.x/scale*float64(w)*0.28))
		sy := int(math.Round(cy + s.y/scale*float64(h)*0.42))
		if ph == 3 && s.z < 0.35 && sx+2 < w {
			streak := "━━"
			st := warpCyan
			if s.hue == 2 {
				st = warpViolet
			} else if s.hue == 3 {
				streak = "──"
				st = warpDim
			}
			placeChunk(grid, w, sx, sy, st.Render(streak), 2)
			continue
		}
		plot(sx, sy, warpStyle(s.hue, s.z < 0.2).Render(warpGlyph(s.z)))
	}

	// Layer 2: supernova burst particles from the physics engine.
	for _, p := range m.engine.Snap() {
		if !p.Burst {
			continue
		}
		var st lipgloss.Style
		switch {
		case p.Age < 0.25:
			st = flashWhite
		case p.Age < 0.6:
			st = warpCyan
		default:
			st = warpViolet
		}
		plot(p.X, p.Y, st.Render(string(p.Glyph)))
	}

	// Layer 3: blinding flash disc right after ignition.
	flash := m.engine.Flash
	if ph == 3 && flash > 0.02 {
		r := int(flash * float64(w) * 0.22)
		fill := "▒"
		if flash > 0.6 {
			fill = "█"
		} else if flash > 0.3 {
			fill = "▓"
		}
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dx := float64(x) - cx
				dy := (float64(y) - cy) * 2.1
				if math.Sqrt(dx*dx+dy*dy) < float64(r) {
					plot(x, y, flashWhite.Render(fill))
				}
			}
		}
	}

	// Layer 4: expanding shockwave ring during ignition.
	if ph == 3 {
		t := (elapsed - act4At) / (introDuration - act4At)
		r := 3 + t*float64(w)*0.45
		for deg := 0; deg < 360; deg += 2 {
			a := float64(deg) * math.Pi / 180
			x := int(math.Round(cx + math.Cos(a)*r))
			y := int(math.Round(cy + math.Sin(a)*r/2.1))
			plot(x, y, warpCyan.Render("○"))
		}
	}

	// Layer 5: act content.
	switch ph {
	case 0:
		// Act 1: telemetry lock in the top corner over the void.
		m.overlayTelemetry(grid, w, elapsed)
	case 1:
		// Act 2: three counter-rotating orbital rings, accelerating.
		m.overlayOrbits(grid, w, h, elapsed, false)
	case 2:
		// Act 3: rings collapse, laser traces the crest blueprint.
		m.overlayOrbits(grid, w, h, elapsed, true)
		m.overlayBlueprint(grid, w, h, elapsed)
	case 3:
		if elapsed > act4TeaseAt {
			// Hyperspace leap tail: hero dissolves, compact header docks
			// on top as the Doctor card fades in.
			m.overlayCrestTop(grid, w, h)
			m.overlayFadingCard(grid, w, h, "PREFLIGHT DOCTOR")
			break
		}
		// Act 4: the true-color crest materializes with its neon halo.
		endY := m.overlayCrest(grid, w, h, m.introArt, true)
		if endY < 0 {
			endY = h/2 + 2
		}
		m.overlayTitle(grid, w, endY+1, 1.0)
	}

	// Layer 6: progress bar + skip hint (bottom rows, always visible).
	m.overlayProgress(grid, w, h, elapsed, ph)

	// Flatten to exactly h lines of exactly w visible cells.
	var sb strings.Builder
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			cell := grid[y][x]
			if cell.skip {
				continue
			}
			if cell.set {
				sb.WriteString(cell.text)
			} else {
				sb.WriteByte(' ')
			}
		}
		sb.WriteString("\x1b[0m")
		if y < h-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// overlayCrest centers the full emblem art, breathing via a neon aura line.
// Returns the last art row (for title placement), or -1 when nothing drawn.
func (m *Model) overlayCrest(grid [][]introCell, w, h int, art string, breathe bool) int {
	lines := strings.Split(art, "\n")
	// Trim empty edge lines chafa may emit.
	for len(lines) > 0 && strings.TrimSpace(stripANSI(lines[0])) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(stripANSI(lines[len(lines)-1])) == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return -1
	}
	startY := h/2 - len(lines)/2 - 1
	for i, ln := range lines {
		y := startY + i
		if y < 1 || y >= h-2 {
			continue
		}
		visW := lipgloss.Width(ln)
		x := (w - visW) / 2
		if x < 0 || x+visW > w {