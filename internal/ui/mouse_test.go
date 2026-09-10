package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMouseWheelScrollsVList(t *testing.T) {
	m := NewModel()
	m.stage = StageMCP
	m.width, m.height = 120, 30
	m.onEnterStage()
	before := m.mcpList.Cursor
	updated, _ := m.Update(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonWheelDown})
	fm := updated.(Model)
	if fm.mcpList.Cursor == before && len(fm.mcpList.Filtered) > 1 {
		t.Fatal("wheel did not scroll MCP list")
	}
}

func TestMouseClickTogglesVList(t *testing.T) {
	m := NewModel()
	m.stage = StageAgents
	m.width, m.height = 120, 30
	m.onEnterStage()
	before := m.agentsList.SelectedCount()
	updated, _ := m.Update(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: 10, Y: 10})
	fm := updated.(Model)
	after := fm.agentsList.SelectedCount()
	if after == before {
		t.Fatal("click did not toggle agent selection")
	}
	if fm.engine.Count() == 0 {
		t.Fatal("click produced no starburst particles")
	}
}

func TestMouseMotionTrail(t *testing.T) {
	m := NewModel()
	m.width, m.height = 120, 30
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
