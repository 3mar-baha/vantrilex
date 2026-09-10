package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// stageModel120 builds a 120x30 model on the given stage with layout ready.
func stageModel120(stage Stage) Model {
	m := NewModel()
	m.width, m.height = 120, 30
	m.stage = stage
	m.onEnterStage()
	return m
}

func leftClick(x, y int) tea.MouseMsg {
	return tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: x, Y: y}
}

func TestMouseWheelScrollsVList(t *testing.T) {
	m := stageModel120(StageMCP)
	before := m.mcpList.Cursor
	updated, _ := m.Update(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonWheelDown})
	fm := updated.(Model)
	if fm.mcpList.Cursor == before && len(fm.mcpList.Filtered) > 1 {
		t.Fatal("wheel did not scroll MCP list")
	}
}

func TestMouseClickTogglesVListRow(t *testing.T) {
	m := stageModel120(StageAgents)
	lo := m.layout()
	y := lo.cardY + 2 + 4 // first visible item row
	x := lo.boxX + 3 + 1  // inside the "[X]" checkbox cell
	before := m.agentsList.SelectedCount()
	updated, _ := m.Update(leftClick(x, y))
	fm := updated.(Model)
	if fm.agentsList.SelectedCount() == before {
		t.Fatal("click on item row did not toggle agent selection")
	}
	if fm.engine.Count() == 0 {
		t.Fatal("click produced no starburst particles")
	}
	if fm.pressID != "row" || !fm.pressFresh() {
		t.Fatal("click produced no pressed-row feedback")
	}
}

func TestFooterDoctorContinueClick(t *testing.T) {
	m := stageModel120(StageDoctor)
	lo := m.layout()
	var target fhit
	found := false
	for _, h := range m.footerHits(lo) {
		if h.key == "c" {
			target, found = h, true
		}
	}
	if !found {
		t.Fatal("no [C] footer hit rect on Doctor stage")
	}
	updated, _ := m.Update(leftClick(target.x+target.w/2, target.y))
	fm := updated.(Model)
	if fm.stage != StageHistory {
		t.Fatalf("clicking [C] advanced to %v, want History", fm.stage)
	}
	if fm.pressID != "btn:c" {
		t.Fatalf("pressID=%q, want btn:c", fm.pressID)
	}
}

func TestFooterRunnerBackClick(t *testing.T) {
	m := stageModel120(StageRunner)
	lo := m.layout()
	for _, h := range m.footerHits(lo) {
		if h.key == "b" {
			updated, _ := m.Update(leftClick(h.x, h.y))
			fm := updated.(Model)
			if fm.stage != StageHistory {
				t.Fatalf("clicking [B] went to %v, want History", fm.stage)
			}
			return
		}
	}
	t.Fatal("no [B] footer hit rect on Runner stage")
}

func TestRunnerRowClickSelects(t *testing.T) {
	m := stageModel120(StageRunner)
	lo := m.layout()
	y := lo.cardY + 2 + 3 + 1*3 // second runner row
	updated, _ := m.Update(leftClick(lo.boxX+4, y))
	fm := updated.(Model)
	if !fm.hasRunner || fm.curRunner.Index != 2 {
		t.Fatalf("row click selected %+v, want runner index 2", fm.curRunner)
	}
	if fm.stage != StageModel {
		t.Fatalf("row click landed on %v, want Model", fm.stage)
	}
}

func TestModelRowClickSelects(t *testing.T) {
	m := stageModel120(StageModel)
	if len(m.filtered) == 0 {
		t.Skip("no models in preview catalog")
	}
	lo := m.layout()
	y := lo.cardY + 2 + 5 // first model entry, first line
	updated, _ := m.Update(leftClick(lo.boxX+4, y))
	fm := updated.(Model)
	if !fm.hasModel {
		t.Fatal("row click did not select a model")
	}
	if fm.stage != StageEffort {
		t.Fatalf("row click landed on %v, want Effort", fm.stage)
	}
}

func TestModelTabRowClickSwitches(t *testing.T) {
	m := stageModel120(StageModel)
	lo := m.layout()
	y := lo.cardY + 2 + 2 // tabs row
	xx := lo.boxX + 3
	segs := modelTabSegs()
	if len(segs) < 2 {
		t.Skip("need 2+ tabs")
	}
	// Click inside the second tab label.
	x := xx + len(segs[0]) + 2
	updated, _ := m.Update(leftClick(x, y))
	fm := updated.(Model)
	if fm.tabCursor != 1 {
		t.Fatalf("tab click set cursor %d, want 1", fm.tabCursor)
	}
}

func TestStageBarJumpBack(t *testing.T) {
	m := stageModel120(StageModel)
	lo := m.layout()
	var target shit
	for _, s := range m.stageBarHits(lo) {
		if s.stage == StageDoctor && s.past {
			target = s
		}
	}
	if target.w == 0 {
		t.Fatal("no clickable Doctor pipeline step")
	}
	updated, _ := m.Update(leftClick(target.x+target.w/2, lo.barY))
	fm := updated.(Model)
	if fm.stage != StageDoctor {
		t.Fatalf("pipeline click landed on %v, want Doctor", fm.stage)
	}
}

func TestStageBarFutureIgnored(t *testing.T) {
	m := NewModel()
	m.width, m.height = 200, 30
	m.stage = StageModel
	m.onEnterStage()
	lo := m.layout()
	for _, s := range m.stageBarHits(lo) {
		if s.stage == StageEffort && !s.past {
			updated, _ := m.Update(leftClick(s.x+s.w/2, lo.barY))
			fm := updated.(Model)
			if fm.stage != StageModel {
				t.Fatalf("future pipeline click moved to %v", fm.stage)
			}
			return
		}
	}
	t.Fatal("no future Effort pipeline step found")
}

func TestStageBarTruncatedHonestly(t *testing.T) {
	// On a 120-wide terminal the 13-step bar overflows; only the visible
	// prefix may yield hit rects, never phantom cells past the edge.
	m := stageModel120(StageModel)
	lo := m.layout()
	for _, s := range m.stageBarHits(lo) {
		if s.x < 0 || s.x+s.w > m.width-1 {
			t.Fatalf("pipeline hit %+v exceeds visible bar", s)
		}
	}
}

func TestClickEmptySpaceIgnored(t *testing.T) {
	m := stageModel120(StageRunner)
	updated, _ := m.Update(leftClick(0, 0)) // top margin row: no target
	fm := updated.(Model)
	if fm.stage != StageRunner {
		t.Fatalf("margin click moved to %v", fm.stage)
	}
	if fm.engine.Count() == 0 {
		t.Fatal("margin click produced no burst feedback")
	}
}

func TestFooterHitsInsideTerminal(t *testing.T) {
	for _, st := range stageOrder() {
		m := stageModel120(st)
		lo := m.layout()
		for _, h := range m.footerHits(lo) {
			if h.x < 0 || h.w <= 0 || h.x+h.w > m.width || h.y != lo.footerY {
				t.Fatalf("stage %v footer hit %+v outside terminal", st, h)
			}
		}
	}
}

func TestMouseMotionTrail(t *testing.T) {
	m := stageModel120(StageAgents)
	before := m.engine.Count()
	updated, _ := m.Update(tea.MouseMsg{Action: tea.MouseActionMotion, X: 5, Y: 5})
	fm := updated.(Model)
	if fm.engine.Count() <= before {
		t.Fatal("hover produced no trail particles")
	}
}

func TestStageTransitionPerimeterSupernova(t *testing.T) {
	m := NewModel()
	m.width, m.height = 120, 30
	m.engine.SetSize(116, 3)
	m.stage = StageDoctor
	m.onEnterStage()
	if m.engine.Bursts() == 0 {
		t.Fatal("stage entry produced no perimeter supernova")
	}
}

func TestMenuEdgeOnlyEnabled(t *testing.T) {
	m := stageModel120(StageDoctor)
	if !m.engine.EdgeOnly() {
		t.Fatal("menu stage did not enable edge-only ambient field")
	}
	m.enterIntro()
	if m.engine.EdgeOnly() {
		t.Fatal("intro must keep the full-field ambient mode")
	}
}

func TestLiveModelsMsg(t *testing.T) {
	m := NewModel()
	m.stage = StageModel
	updated, _ := m.Update(liveModelsMsg{models: nil})
	_ = updated
	updated2, _ := m.Update(runnersDoneMsg{errs: map[string]error{}})
	fm := updated2.(Model)
	if fm.runnerUpd.Pill() != "[RUNNERS UP-TO-DATE]" {
		t.Fatalf("pill=%q", fm.runnerUpd.Pill())
	}
}
