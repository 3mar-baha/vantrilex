package ui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"vantrilex/internal/doctor"
)

func TestIntroStartsFirst(t *testing.T) {
	m := NewModel()
	if m.stage != StageIntro {
		t.Fatalf("expected StageIntro on startup, got %v", m.stage)
	}
}

func TestIntroPhaseBoundaries(t *testing.T) {
	cases := []struct {
		elapsed float64
		want    int
	}{
		{0.0, 0}, {2.49, 0}, {2.5, 1}, {4.99, 1},
		{5.0, 2}, {7.49, 2}, {7.5, 3}, {9.99, 3}, {10.0, 4}, {12.0, 4},
	}
	for _, c := range cases {
		if got := introPhase(c.elapsed); got != c.want {
			t.Fatalf("introPhase(%v) = %d, want %d", c.elapsed, got, c.want)
		}
	}
}

func TestIntroSkipKeys(t *testing.T) {
	keys := []tea.KeyMsg{
		{Type: tea.KeySpace},
		{Type: tea.KeyEnter},
		{Type: tea.KeyEsc},
		{Type: tea.KeyRunes, Runes: []rune{'x'}},
	}
	for _, k := range keys {
		m := NewModel()
		updated, _ := m.Update(k)
		fm := updated.(Model)
		if fm.stage != StageDoctor {
			t.Fatalf("key %q did not skip intro (stage=%v)", k.String(), fm.stage)
		}
	}
}

func TestIntroAutoAdvance(t *testing.T) {
	m := NewModel()
	m.introStart = time.Now().Add(-11 * time.Second)
	updated, _ := m.Update(tickMsg(time.Now()))
	fm := updated.(Model)
	if fm.stage != StageDoctor {
		t.Fatalf("intro did not auto-advance after 10s (stage=%v)", fm.stage)
	}
}

func TestIntroDetonationBurst(t *testing.T) {
	m := NewModel()
	m.width, m.height = 120, 30
	m.introStart = time.Now().Add(-8 * time.Second)
	updated, _ := m.Update(tickMsg(time.Now()))
	fm := updated.(Model)
	if !fm.introBurst {
		t.Fatal("supernova burst did not trigger at t>=7.5s")
	}
	if n := len(fm.engine.Snap()); n == 0 {
		t.Fatal("engine has no particles after detonation")
	}
}

func TestIntroViewExactHeight(t *testing.T) {
	for _, size := range [][2]int{{120, 30}, {100, 24}, {80, 20}} {
		m := NewModel()
		m.width, m.height = size[0], size[1]
		m.engine.SetSize(size[0], size[1])
		view := m.View()
		if got := countLines(view); got != size[1] {
			t.Fatalf("intro view at %dx%d = %d lines, want exactly %d", size[0], size[1], got, size[1])
		}
	}
}

// TestIntroRowWidths ensures no intro row exceeds the terminal width in
// visible cells (which would cause terminal wrap = visual scrolling).
// It exercises all four phases with emblem art and burst particles live.
func TestIntroRowWidths(t *testing.T) {
	phases := []time.Duration{1 * time.Second, 7 * time.Second / 2, 6 * time.Second, 17 * time.Second / 2}
	for _, size := range [][2]int{{120, 30}, {100, 24}, {80, 20}} {
		for _, age := range phases {
			m := NewModel()
			m.width, m.height = size[0], size[1]
			m.engine.SetSize(size[0], size[1])
			m.engine.Supernova(130)
			m.introArt = fallbackLogo()
			m.introStart = time.Now().Add(-age)
			for i, row := range strings.Split(m.View(), "\n") {
				if got := len([]rune(stripANSI(row))); got != size[0] {
					t.Fatalf("intro %dx%d phase %v row %d visible width %d, want %d",
						size[0], size[1], age, i, got, size[0])
				}
			}
		}
	}
}

func TestMenuViewNeverExceedsBudget(t *testing.T) {
	stages := []Stage{StageDoctor, StageHistory, StageRunner, StageModel, StageEffort, StageWorkspace, StageAgents, StageSkillsPick, StagePlugins, StageHooks, StageMCP, StageSkills, StageLaunch}
	sizes := [][2]int{{140, 40}, {120, 30}, {100, 24}, {80, 20}}
	for _, st := range stages {
		for _, size := range sizes {
			m := NewModel()
			m.stage = st
			m.width, m.height = size[0], size[1]
			m.statuses = doctor.CheckAll()
			m.doctorCheck = true
			m.onEnterStage()
			if st == StageWorkspace {
				m.wsInput.SetValue("C:\\projects\\demo")
			}
			view := m.View()
			if got := countLines(view); got > size[1]-1 {
				t.Fatalf("stage %v at %dx%d = %d lines, exceeds budget %d", st, size[0], size[1], got, size[1]-1)
			}
		}
	}
}

func TestOrbitAngleGrowsExponentially(t *testing.T) {
	a0 := orbitAngle(0, 1, 1.2, 1.1, 0)
	a1 := orbitAngle(0, 1, 1.2, 1.1, 1)
	a2 := orbitAngle(0, 1, 1.2, 1.1, 2)
	if a0 != 0 {
		t.Fatalf("orbitAngle(tau=0) = %v, want 0", a0)
	}
	if !(a2-a1 > a1-a0) {
		t.Fatal("orbital rings are not accelerating")
	}
	if back := orbitAngle(0, -1, 1.2, 1.1, 2); back >= 0 {
		t.Fatal("counter-rotating ring does not move backwards")
	}
}

func TestBlueprintDefined(t *testing.T) {
	if len(crestBlueprint) < 5 {
		t.Fatalf("blueprint has %d polylines, want >=5", len(crestBlueprint))
	}
	for i, line := range crestBlueprint {
		if len(line) < 2 {
			t.Fatalf("blueprint line %d has <2 points", i)
		}
	}
}

func TestFitHelpers(t *testing.T) {
	if got := countLines("a\nb\nc"); got != 3 {
		t.Fatalf("countLines = %d, want 3", got)
	}
	if got := fitLines("a\nb\nc\nd", 2); got != "a\nb" {
		t.Fatalf("fitLines = %q, want %q", got, "a\nb")
	}
	if got := clampInt(99, 8, 18); got != 18 {
		t.Fatalf("clampInt = %d, want 18", got)
	}
	if strings.Contains(stripANSI("\x1b[36mok\x1b[0m"), "\x1b") {
		t.Fatal("stripANSI left escape sequences")
	}
	long := strings.Repeat("x", 200)
	if got := fitWidthLines(long, 40); len([]rune(stripANSI(got))) > 40 {
		t.Fatal("fitWidthLines did not truncate")
	}
}
