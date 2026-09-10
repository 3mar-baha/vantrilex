package particles

import (
	"strings"
	"testing"
)

func TestSupernovaBurstCount(t *testing.T) {
	e := New(100, 6)
	e.Supernova(80)
	if e.Count() < 80 {
		t.Fatalf("expected >=80 particles, got %d", e.Count())
	}
	// Tick 5 seconds at 60fps: bursts should decay away.
	for i := 0; i < 300; i++ {
		e.Tick(1.0 / 60.0)
	}
	if e.Count() > 60 {
		t.Fatalf("bursts did not decay, count=%d", e.Count())
	}
	if e.Bursts() != 1 {
		t.Fatalf("expected 1 burst, got %d", e.Bursts())
	}
	if got := e.RenderField(); len(got) == 0 {
		t.Fatal("empty field render")
	}
}

func TestClickBurstCount(t *testing.T) {
	e := New(100, 30)
	e.ClickBurst(50, 15)
	if got := e.Count(); got < 12 || got > 16 {
		t.Fatalf("click burst = %d, want 12-16", got)
	}
}

func TestPerimeterSupernova(t *testing.T) {
	e := New(80, 24)
	e.PerimeterSupernova()
	if e.Count() == 0 {
		t.Fatal("perimeter supernova emitted nothing")
	}
	if e.Bursts() != 1 {
		t.Fatalf("bursts=%d", e.Bursts())
	}
	found := false
	for _, s := range e.Snap() {
		if s.X <= 1 || s.Y <= 1 {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("no edge particles")
	}
}

func TestTrailEmits(t *testing.T) {
	e := New(100, 30)
	e.EmitTrail(10, 5)
	if e.Count() != 2 {
		t.Fatalf("trail=%d, want 2", e.Count())
	}
}

func TestOverlayFlags(t *testing.T) {
	e := New(100, 30)
	e.EmitTrail(10, 5)
	e.ClickBurst(50, 15)
	e.PerimeterSupernova()
	e.Supernova(80)
	var trail, click, perim, center int
	for _, s := range e.Snap() {
		if !s.Overlay {
			center++
			continue
		}
		switch {
		case s.Glyph == '·' || s.Glyph == '⋆' || s.Glyph == '∘':
			trail++
		case s.Glyph == '✦' || s.Glyph == '✧' || s.Glyph == '⚡':
			click++
		default:
			perim++
		}
	}
	if trail == 0 || click == 0 || perim == 0 {
		t.Fatalf("overlay classes missing: trail=%d click=%d perim=%d", trail, click, perim)
	}
	if center == 0 {
		t.Fatal("expected non-overlay center supernova particles")
	}
}

func TestOverlaySnapsCap(t *testing.T) {
	e := New(100, 30)
	for i := 0; i < 50; i++ {
		e.EmitTrail(i, 5)
	}
	got := e.OverlaySnaps()
	if len(got) > 80 {
		t.Fatalf("overlay snaps=%d, want cap 80", len(got))
	}
	if len(got) == 0 {
		t.Fatal("expected overlay snaps")
	}
	for _, s := range got {
		if !s.Overlay {
			t.Fatal("non-overlay snapshot leaked into overlay set")
		}
	}
}

func TestEdgeOnlyHidesCenter(t *testing.T) {
	e := New(80, 8)
	e.SetEdgeOnly(true)
	if !e.EdgeOnly() {
		t.Fatal("edge flag not set")
	}
	// Deterministic probes: 'C' dead-center, 'E' on the far-left margin.
	e.parts = append(e.parts,
		&Particle{X: 40, Y: 4, Life: 5, MaxLife: 7, Glyph: 'C'},
		&Particle{X: 2, Y: 4, Life: 5, MaxLife: 7, Glyph: 'E'},
	)
	rendered := e.RenderField()
	if strings.Contains(rendered, "C") {
		t.Fatal("center ambient star leaked through edge-only field")
	}
	if !strings.Contains(rendered, "E") {
		t.Fatal("margin ambient star missing from edge-only field")
	}
	// Overlay FX still draw inside the center frame.
	e.ClickBurst(40, 4)
	found := false
	for _, s := range e.OverlaySnaps() {
		if s.X == 40 && s.Y == 4 {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("click overlay missing at center under edge-only mode")
	}
}

func TestEdgeOnlyAmbientSpawnOnMargins(t *testing.T) {
	e := New(80, 8)
	e.SetEdgeOnly(true)
	for i := 0; i < 200; i++ {
		x, y := e.ambientPos()
		if e.inCenterFrame(x, y, e.Width, e.Height) {
			t.Fatalf("ambient spawn inside center frame at %.1f,%.1f", x, y)
		}
	}
}
