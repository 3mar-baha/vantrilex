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