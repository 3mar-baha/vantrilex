// Package particles implements a 60 FPS supernova particle engine.
// Physics runs on fixed 16ms ticks: radial bursts with velocity decay plus
// continuous low-cost ambient drift. Rendering produces decorative
// starfield strips composited by the UI layer.
package particles

import (
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var glyphs = []rune{'✦', '✧', '⋆', '*', '•', '·', '✺', '❋'}

// Particle is a single stardust mote.
type Particle struct {
	X, Y     float64 // position in field coordinates
	VX, VY   float64 // velocity
	Life     float64 // remaining life 0..1
	MaxLife  float64
	Glyph    rune
	Size     float64
	Burst    bool // true for supernova burst particles
	BornAt   time.Time
	HuePhase float64
}

// Engine holds ambient + burst particles.
type Engine struct {
	Width, Height int
	parts         []*Particle
	bursts        int
	rng           *rand.Rand
	lastAmbient   time.Time
	Flash         float64 // 0..1 screen flash right after a burst
	frame         uint64
}

// New creates an engine with a default field size.
func New(w, h int) *Engine {
	if w <= 0 {
		w = 100
	}
	if h <= 0 {
		h = 8
	}
	return &Engine{
		Width:  w,
		Height: h,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SetSize resizes the field.
func (e *Engine) SetSize(w, h int) {
	if w > 4 {
		e.Width = w
	}
	if h > 2 {
		e.Height = h
	}
}

// Count returns live particle count.
func (e *Engine) Count() int { return len(e.parts) }

// Bursts returns total bursts triggered.
func (e *Engine) Bursts() int { return e.bursts }

var trailGlyphs = []rune{'·', '⋆', '∘'}
var burstGlyphs = []rune{'✦', '✧', '⚡'}

// EmitTrail leaves 1-2 decaying micro-particles behind the cursor path.
func (e *Engine) EmitTrail(x, y int) {
	if os.Getenv("VANTRILEX_NO_FX") != "" {
		return
	}
	for i := 0; i < 2; i++ {
		e.parts = append(e.parts, &Particle{
			X:        float64(x) + (e.rng.Float64()-0.5)*1.2,
			Y:        float64(y) + (e.rng.Float64()-0.5)*0.8,
			VX:       (e.rng.Float64() - 0.5) * 1.4,
			VY:       -0.4 - e.rng.Float64()*0.8,
			Life:     0.5 + e.rng.Float64()*0.5,
			MaxLife:  1.0,
			Glyph:    trailGlyphs[e.rng.Intn(len(trailGlyphs))],
			Burst:    true,
			BornAt:   time.Now(),
			HuePhase: e.rng.Float64(),
		})
	}
}

// ClickBurst detonates 12-16 sparkling stars radially from click coords.
func (e *Engine) ClickBurst(x, y int) {
	if os.Getenv("VANTRILEX_NO_FX") != "" {
		return
	}
	n := 12 + e.rng.Intn(5)
	for i := 0; i < n; i++ {
		ang := e.rng.Float64() * 2 * math.Pi
		speed := 4 + e.rng.Float64()*14
		life := 0.5 + e.rng.Float64()*0.6
		e.parts = append(e.parts, &Particle{
			X:        float64(x),
			Y:        float64(y),
			VX:       math.Cos(ang) * speed,
			VY:       math.Sin(ang) * speed * 0.55,
			Life:     life,
			MaxLife:  life,
			Glyph:    burstGlyphs[e.rng.Intn(len(burstGlyphs))],
			Burst:    true,
			BornAt:   time.Now(),
			HuePhase: e.rng.Float64(),
		})
	}
	e.bursts++
}

// PerimeterSupernova erupts stardust along the full terminal perimeter,
// traveling inwards and fading into deep space. Call on stage transitions.
func (e *Engine) PerimeterSupernova() {
	if os.Getenv("VANTRILEX_NO_FX") != "" {
		return
	}
	w, h := e.Width, e.Height
	if w < 8 {
		w = 8
	}
	if h < 4 {
		h = 4
	}
	perim := [][2]int{}
	for x := 0; x < w; x += 2 {
		perim = append(perim, [2]int{x, 0}, [2]int{x, h - 1})
	}
	for y := 0; y < h; y++ {
		perim = append(perim, [2]int{0, y}, [2]int{w - 1, y})
	}
	cx, cy := float64(w)/2, float64(h)/2
	for _, p := range perim {
		dx, dy := cx-float64(p[0]), cy-float64(p[1])
		dist := math.Hypot(dx, dy)
		if dist == 0 {
			continue
		}
		speed := 6 + e.rng.Float64()*10
		life := 0.8 + e.rng.Float64()*0.8
		e.parts = append(e.parts, &Particle{
			X:        float64(p[0]),
			Y:        float64(p[1]),
			VX:       dx / dist * speed,
			VY:       dy / dist * speed * 0.55,
			Life:     life,
			MaxLife:  life,
			Glyph:    glyphs[e.rng.Intn(len(glyphs))],
			Burst:    true,
			BornAt:   time.Now(),
			HuePhase: e.rng.Float64(),
		})
	}
	e.bursts++
	if e.Flash < 0.9 {
		e.Flash = 0.7
	}
}

// Supernova triggers a radial explosion of n stardust particles from center.
func (e *Engine) Supernova(n int) {
	if n < 80 {
		n = 80
	}
	cx := float64(e.Width) / 2
	cy := float64(e.Height) / 2
	for i := 0; i < n; i++ {
		ang := e.rng.Float64() * 2 * math.Pi
		speed := 6 + e.rng.Float64()*26
		life := 0.7 + e.rng.Float64()*0.9
		e.parts = append(e.parts, &Particle{
			X:        cx,
			Y:        cy,
			VX:       math.Cos(ang) * speed,
			VY:       math.Sin(ang) * speed * 0.55,
			Life:     life,
			MaxLife:  life,
			Glyph:    glyphs[e.rng.Intn(len(glyphs))],
			Size:     0.5 + e.rng.Float64(),
			Burst:    true,
			BornAt:   time.Now(),
			HuePhase: e.rng.Float64(),
		})
	}
	e.bursts++
	if e.Flash < 0.9 {
		e.Flash = 1.0
	}
}

// Tick advances physics by dt seconds. Called at 60 FPS.
func (e *Engine) Tick(dt float64) {
	e.frame++
	if dt <= 0 || dt > 0.1 {
		dt = 1.0 / 60.0
	}
	// Decay flash.
	if e.Flash > 0 {
		e.Flash -= dt * 1.8
		if e.Flash < 0 {
			e.Flash = 0
		}