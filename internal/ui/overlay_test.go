package ui

import (
	"strings"
	"testing"

	"vantrilex/internal/particles"
)

func TestOverlaySpliceReplacesCell(t *testing.T) {
	base := "abc\ndef\nghi"
	snaps := []particles.Snapshot{{X: 1, Y: 1, Glyph: 'X', Age: 0.1, Burst: true}}
	out := overlayFX(base, snaps, 10, 4)
	rows := strings.Split(out, "\n")
	if len(rows) != 3 {
		t.Fatalf("row count=%d, want 3", len(rows))
	}
	if !strings.Contains(rows[1], "X") {
		t.Fatalf("overlay glyph missing: %q", rows[1])
	}
	if strings.Contains(rows[0], "X") || strings.Contains(rows[2], "X") {
		t.Fatal("overlay glyph leaked to other rows")
	}
}

func TestOverlayIgnoresOutOfBounds(t *testing.T) {
	base := "abc\ndef"
	snaps := []particles.Snapshot{
		{X: -1, Y: 0, Glyph: 'X', Burst: true},
		{X: 99, Y: 0, Glyph: 'X', Burst: true},
		{X: 0, Y: 99, Glyph: 'X', Burst: true},
	}
	if got := overlayFX(base, snaps, 10, 3); got != base {
		t.Fatalf("out-of-bounds snaps altered frame: %q", got)
	}
}

func TestOverlayPreservesBudget(t *testing.T) {
	m := stageModel120(StageMCP)
	// Seed live overlay particles across the frame.
	m.engine.ClickBurst(60, 15)
	m.engine.EmitTrail(20, 25)
	view := m.View()
	rows := strings.Split(view, "\n")
	if len(rows) > m.height-1 {
		t.Fatalf("overlay grew frame to %d rows, budget %d", len(rows), m.height-1)
	}
	for i, r := range rows {
		if got := len([]rune(stripANSI(r))); got > m.width {
			t.Fatalf("row %d width %d exceeds %d", i, got, m.width)
		}
	}
}

func TestOverlayEmptyIsIdentity(t *testing.T) {
	base := "ab\ncd"
	if got := overlayFX(base, nil, 10, 3); got != base {
		t.Fatalf("empty overlay changed frame: %q", got)
	}
}
