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