package particles

import "testing"

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
