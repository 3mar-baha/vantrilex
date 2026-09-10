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
			continue // never corrupt ANSI or overflow the row
		}
		placeChunk(grid, w, x, y, ln+"\x1b[0m", visW)
	}
	if breathe {
		// Breathing aura: a soft divider glowing under the crest.
		ay := startY + len(lines) + 1
		if ay > 0 && ay < h-2 {
			pulse := 0.5 + 0.5*math.Sin(time.Since(m.introStart).Seconds()*3.2)
			glyph := "─"
			if pulse > 0.66 {
				glyph = "═"
			}
			aura := strings.Repeat(glyph, min(46, w-8))
			st := warpCyan
			if pulse < 0.4 {
				st = warpDim
			}
			putCentered(grid, w, ay, st.Render(aura))
			return ay // title must go below the aura, never on it
		}
	}
	return startY + len(lines) - 1
}

// overlayCrestTop docks the compact header crest at the top during the
// hyperspace-leap tail (same 28x14 art as menu headers).
func (m *Model) overlayCrestTop(grid [][]introCell, w, h int) {
	art := m.logo
	if art == "" {
		art = fallbackLogo()
	}
	lines := strings.Split(fitLines(art, 14), "\n")
	y := 1 // 1-line top safety margin
	for _, ln := range lines {
		if y >= h-3 {
			break
		}
		visW := lipgloss.Width(ln)
		x := (w - visW) / 2
		if x < 0 || x+visW > w {
			y++
			continue
		}
		placeChunk(grid, w, x, y, ln+"\x1b[0m", visW)
		y++
	}
}

// overlayTitle draws the Vantrilex wordmark at row y with the Workflow
// Launcher subtitle: together they read "Vantrilex Workflow Launcher".
func (m *Model) overlayTitle(grid [][]introCell, w, y int, strength float64) {
	title := "V  A  N  T  R  I  L  E  X"
	sub := "WORKFLOW LAUNCHER"
	st := introTitle
	if strength < 1 {
		st = introSub
	}
	putCentered(grid, w, y, st.Render(title))
	putCentered(grid, w, y+1, introHint.Render(sub))
}

// overlayFadingCard teases the incoming Doctor card.
func (m *Model) overlayFadingCard(grid [][]introCell, w, h int, label string) {
	putCentered(grid, w, h-5, warpCyan.Render("◆ "+label+" ◆"))
}

// overlayTelemetry draws Act 1 deep-void telemetry in the top-left corner:
// sector tag, radar pulse and the green signal-lock indicator.
func (m *Model) overlayTelemetry(grid [][]introCell, w int, elapsed float64) {
	tag := introHint.Render("[DEEP VOID: SECTOR 0x7F]")
	placeChunk(grid, w, 2, 1, tag, lipgloss.Width(tag))
	radar := []string{"◉", "◎", "○", "◌"}
	pulse := radar[int(elapsed*4)%len(radar)]
	lock := "SCANNING··"
	st := scanAmber
	if elapsed > 1.2 {
		lock = "SIGNAL ACQUIRED // INITIATING QUANTUM CORE"
		st = lockGreen
	} else if int(elapsed*2)%2 == 0 {
		lock = "SCANNING..."
	}
	line := st.Render(pulse + " " + lock)
	placeChunk(grid, w, 2, 2, line, lipgloss.Width(line))
}

// orbitAngle integrates exponential angular acceleration in closed form:
// theta(tau) = base + dir*w0*(e^(lambda*tau)-1)/lambda. Stateless, so the
// 60 FPS render stays perfectly smooth with zero stored animation state.
func orbitAngle(base, dir, w0, lambda, tau float64) float64 {
	if tau < 0 {
		tau = 0
	}
	return base + dir*w0*(math.Exp(lambda*tau)-1)/lambda
}

// ringColor shifts cobalt -> ice blue -> electric violet across Act 2.
func ringColor(frac float64) lipgloss.Style {
	if frac < 0 {
		frac = 0
	}
	if frac > 1 {
		frac = 1
	}
	lerp := func(a, b int, f float64) int { return int(float64(a) + float64(b-a)*f) }
	var r, g, bl int
	switch {
	case frac < 0.5:
		f := frac / 0.5
		r, g, bl = lerp(0x1E, 0x38, f), lerp(0x1B, 0xBD, f), lerp(0x4B, 0xF8, f)
	default:
		f := (frac - 0.5) / 0.5
		r, g, bl = lerp(0x38, 0x63, f), lerp(0xBD, 0x66, f), lerp(0xF8, 0xF1, f)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(
		"#" + hex2(r) + hex2(g) + hex2(bl)))
}

func hex2(v int) string {
	const digits = "0123456789ABCDEF"
	if v < 0 {
		v = 0
	}
	if v > 255 {
		v = 255
	}
	return string([]byte{digits[v>>4], digits[v&0xF]})
}

// overlayOrbits draws Act 2's three counter-rotating stardust rings. When
// collapse is true (Act 3) the rings shrink toward the singularity.
func (m *Model) overlayOrbits(grid [][]introCell, w, h int, elapsed float64, collapse bool) {
	tau := elapsed - act2At
	frac := tau / (act3At - act2At)
	if collapse {
		frac = 1
	}
	base := math.Min(float64(w)*0.30, float64(h)*0.60)
	if base < 6 {
		base = 6
	}
	shrink := 1.0
	if collapse {
		p := (elapsed - act3At) / (act4At - act3At)
		shrink = 1 - p
		if shrink < 0.04 {
			shrink = 0.04
		}
		tau = (act3At - act2At) + p*1.2 // rings spin up as they fall in
	}
	cx, cy := float64(w)/2, float64(h)/2-1
	rings := []struct {
		mult float64
		dir  float64
		base float64
	}{
		{0.55, 1, 0.0},
		{0.80, -1, 2.1},
		{1.05, 1, 4.2},
	}
	glyphs := []string{"◈", "✦", "∘", "✧"}
	st := ringColor(frac)
	for _, r := range rings {
		rad := base * r.mult * shrink
		th := orbitAngle(r.base, r.dir, 1.2, 1.1, tau)
		for i := 0; i < 26; i++ {
			a := th + float64(i)/26*2*math.Pi
			x := int(math.Round(cx + math.Cos(a)*rad))
			y := int(math.Round(cy + math.Sin(a)*rad/2))
			if x < 0 || x >= w || y < 0 || y >= h {
				continue
			}
			if !grid[y][x].set {
				grid[y][x] = introCell{text: st.Render(glyphs[i%len(glyphs)]), set: true}
			}
		}
	}
	if collapse && shrink <= 0.05 {
		// Singularity point: hot white core flashing at the center.
		if int(elapsed*12)%2 == 0 {
			plotAt(grid, w, h, int(cx), int(cy), flashWhite.Render("✦"))
		}
	}
}

func plotAt(grid [][]introCell, w, h, x, y int, text string) {
	if x < 0 || x >= w || y < 0 || y >= h {
		return
	}
	if grid[y][x].skip || grid[y][x].chunk {
		return
	}
	grid[y][x] = introCell{text: text, set: true}
}

// crestPoint is one blueprint vertex in normalized crest space (y down).
type crestPoint struct{ x, y float64 }

// crestBlueprint sketches the holographic wireframe: twin horns curving up,
// stepped side wings, central diamond spire, base chevron.
var crestBlueprint = [][]crestPoint{
	{{-0.10, -0.95}, {-0.32, -0.68}, {-0.48, -0.32}, {-0.40, 0.05}},
	{{0.10, -0.95}, {0.32, -0.68}, {0.48, -0.32}, {0.40, 0.05}},
	{{0, -0.62}, {0.09, -0.18}, {0, 0.08}, {-0.09, -0.18}, {0, -0.62}},
	{{-0.40, 0.05}, {-0.72, 0.22}, {-0.60, 0.55}, {-0.34, 0.92}},
	{{0.40, 0.05}, {0.72, 0.22}, {0.60, 0.55}, {0.34, 0.92}},
	{{-0.34, 0.92}, {0, 1.0}, {0.34, 0.92}},
}

// overlayBlueprint laser-traces the crest wireframe with electrostatic
// micro-sparks at the sketch tip.
func (m *Model) overlayBlueprint(grid [][]introCell, w, h int, elapsed float64) {
	p := (elapsed - act3At) / (act4At - act3At)
	if p < 0 {
		return
	}
	if p > 1 {
		p = 1
	}
	reveal := p * p // accelerating trace
	base := math.Min(float64(w)*0.30, float64(h)*0.60)
	if base < 6 {
		base = 6
	}
	cx, cy := float64(w)/2, float64(h)/2-1
	type pt struct{ x, y int }
	var trace []pt
	for _, line := range crestBlueprint {
		for i := 0; i+1 < len(line); i++ {
			a, b := line[i], line[i+1]
			steps := 16
			for s := 0; s < steps; s++ {
				f := float64(s) / float64(steps)
				x := int(math.Round(cx + (a.x+(b.x-a.x)*f)*base*1.15))
				y := int(math.Round(cy + (a.y+(b.y-a.y)*f)*base*0.55))
				trace = append(trace, pt{x, y})
			}
		}
	}
	drawn := int(reveal * float64(len(trace)))
	for i := 0; i < drawn && i < len(trace); i++ {
		glyph := "·"
		st := actIce
		if i%7 == 0 {
			glyph = "✦"
			st = actLaser
		}
		plotAt(grid, w, h, trace[i].x, trace[i].y, st.Render(glyph))
	}
	// Micro-sparks crackling at the sketch tip.
	if drawn > 0 && drawn < len(trace) && m.warpRng != nil {
		tip := trace[drawn-1]
		sparks := []struct {
			dx, dy int
			g      string
			st     lipgloss.Style
		}{
			{0, 0, "⚡", warpViolet},
			{1, -1, "✦", flashWhite},
			{-1, 1, "✦", warpCyan},
			{2, 0, "·", actLaser},
		}
		j := m.warpRng.Intn(2)
		for i := j; i < len(sparks); i++ {
			s := sparks[(i+m.warpRng.Intn(len(sparks)))%len(sparks)]
			plotAt(grid, w, h, tip.x+s.dx, tip.y+s.dy, s.st.Render(s.g))
		}
	}
}

// overlayProgress draws the timeline bar + skip hint on the bottom rows.
func (m *Model) overlayProgress(grid [][]introCell, w, h int, elapsed float64, ph int) {
	frac := elapsed / introDuration
	if frac > 1 {
		frac = 1
	}
	barW := min(44, w-10)
	if barW < 10 {
		barW = 10
	}
	filled := int(math.Round(frac * float64(barW)))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barW-filled)
	putCentered(grid, w, h-3, warpViolet.Render(introPhaseName(ph)))
	putCentered(grid, w, h-2, warpCyan.Render("["+bar+"]"))
	putCentered(grid, w, h-1, introHint.Render("press SPACE / ENTER / ESC to skip"))
}

// putCentered writes a pre-styled chunk centered on a row.
func putCentered(grid [][]introCell, w, y int, styled string) {
	if y < 0 || y >= len(grid) {
		return
	}
	visW := lipgloss.Width(styled)
	x := (w - visW) / 2
	if x < 0 || x+visW > w {
		return
	}
	placeChunk(grid, w, x, y, styled+"\x1b[0m", visW)
}

// placeChunk stores a multi-cell chunk and marks covered cells as skipped.
// At most one chunk owns a row: any previous chunk (and its skip marks) is
// removed first so overlapping chunks can never double-emit.
func placeChunk(grid [][]introCell, w, x, y int, text string, visW int) {
	if y < 0 || y >= len(grid) || x < 0 || x >= w {
		return
	}
	for k := 0; k < w; k++ {
		if grid[y][k].chunk || grid[y][k].skip {
			grid[y][k] = introCell{}
		}
	}
	grid[y][x] = introCell{text: text, set: true, chunk: true}
	for k := x + 1; k < x+visW && k < w; k++ {
		grid[y][k] = introCell{skip: true}
	}
}
