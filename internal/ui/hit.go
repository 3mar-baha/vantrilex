package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"vantrilex/internal/catalog"
)

// fseg is one footer button segment: a styled head, a plain tail, and the
// key fed to the stage handler when clicked ("" = informational only).
type fseg struct {
	head string
	tail string
	key  string
}

func (s fseg) plain() string { return s.head + s.tail }

// footerSegs returns the exact button segments rendered by footerKeys,
// shared by the renderer and the mouse hit-tester so geometry can never
// drift from what is on screen.
func (m Model) footerSegs() []fseg {
	narrow := m.width < 100
	if narrow {
		switch m.stage {
		case StageDoctor:
			return []fseg{{"[A]", " fix", "a"}, {"[C]", " go", "c"}, {"[R]", " scan", "r"}, {"[S]", " sync", "s"}}
		case StageHistory:
			return []fseg{{"[1-5]", " resume", "enter"}, {"[N]", " new", "n"}, {"[D]", " del", "d"}}
		case StageRunner:
			return []fseg{{"[1-3]", " select", "enter"}, {"[B]", " back", "b"}}
		case StageModel:
			return []fseg{{"[Enter]", " select", "enter"}, {"[[]]", " tabs", "]"}, {"[B]", " back", "b"}}
		case StageEffort:
			return []fseg{{"[←→]", " move", ""}, {"[Enter]", " ok", "enter"}, {"[B]", " back", "b"}}
		case StageWorkspace:
			return []fseg{{"[Enter]", " ok", "enter"}, {"[Esc]", " back", "esc"}}
		case StageAgents, StageSkillsPick, StagePlugins, StageHooks, StageMCP:
			return []fseg{{"[Space]", " toggle", " "}, {"[A]", " all", "a"}, {"[C]", " confirm", "c"}, {"[/]", " search", ""}}
		case StageSkills:
			return []fseg{{"[Enter]", " provision", "enter"}, {"[S]", " skip", "s"}}
		case StageLaunch:
			return []fseg{{"[Enter]", " launch", "enter"}, {"[B]", " back", "b"}}
		}
		return nil
	}
	switch m.stage {
	case StageDoctor:
		return []fseg{{"[A]", " approve + auto-install", "a"}, {"[C]", " continue", "c"}, {"[R]", " rescan", "r"}, {"[S]", " sync toolkit", "s"}, {"[Ctrl+C]", " quit", "ctrl+c"}}
	case StageHistory:
		return []fseg{{"[1-5]/Enter", " resume", "enter"}, {"[N]", " new project", "n"}, {"[D]", " delete", "d"}, {"[↑↓]", " navigate", ""}, {"[Ctrl+C]", " quit", "ctrl+c"}}
	case StageRunner:
		return []fseg{{"[1-3]/Enter", " select", "enter"}, {"[↑↓]", " navigate", ""}, {"[B]", " back", "b"}}
	case StageModel:
		return []fseg{{"[Enter]", " select", "enter"}, {"[ [ ] ]", " tabs", "]"}, {"type or [/]", " search", ""}, {"[B/Esc]", " back", "b"}}
	case StageEffort:
		return []fseg{{"[←→]", " move slider", ""}, {"[Enter]", " confirm", "enter"}, {"[B]", " back", "b"}}
	case StageWorkspace:
		return []fseg{{"[Enter]", " confirm", "enter"}, {"[Esc]", " back", "esc"}}
	case StageAgents, StageSkillsPick, StagePlugins, StageHooks, StageMCP:
		return []fseg{{"[Space]", " toggle", " "}, {"[A]", " toggle-all-filtered", "a"}, {"[C]", " confirm", "c"}, {"[/]", " search", ""}, {"[B/Esc]", " back", "b"}}
	case StageSkills:
		return []fseg{{"[Enter]", " provision + continue", "enter"}, {"[S]", " skip if ready", "s"}, {"[B]", " back", "b"}}
	case StageLaunch:
		return []fseg{{"[Enter]", " launch now", "enter"}, {"[O]", " open folder", "o"}, {"[B]", " back", "b"}}
	}
	return nil
}

// footerSep is the exact joiner used between footer segments per mode.
func (m Model) footerSep() string {
	if m.width < 100 {
		return "  "
	}
	return "   "
}

// layout mirrors View's vertical budget exactly: one shared computation
// drives both rendering and mouse hit-testing.
type layout struct {
	logoH, fieldH int
	headerH       int
	barY          int
	cardY, cardH  int
	footerY, fh   int
	contentW      int
	boxX, boxW    int
	innerMax      int
}

func (m Model) layout() layout {
	var lo layout
	w, h := m.width, m.height
	if w <= 0 || h <= 0 {
		return lo
	}
	footer := fitWidthLines(m.renderFooter(), w)
	lo.fh = countLines(footer)
	logoH, fieldH := 14, 3
	for logoH > 4 && 1+logoH+fieldH+1+1+3+lo.fh > h-1 {
		logoH--
	}
	for fieldH > 1 && 1+logoH+fieldH+1+1+3+lo.fh > h-1 {
		fieldH--
	}
	lo.logoH, lo.fieldH = logoH, fieldH
	lo.headerH = countLines(m.renderHeader(logoH, fieldH))
	lo.barY = 1 + lo.headerH
	lo.cardY = lo.barY + 1
	lo.cardH = h - 1 - (1 + lo.headerH + 1 + lo.fh)
	if lo.cardH < 3 {
		lo.cardH = 3
	}
	lo.footerY = lo.cardY + lo.cardH
	lo.contentW = min(92, w-8)
	if lo.contentW < 16 {
		lo.contentW = 16
	}
	lo.boxW = lo.contentW + 6 // border (2) + horizontal padding (4)
	lo.boxX = (w - lo.boxW) / 2
	lo.innerMax = lo.cardH - 4
	if lo.innerMax < 1 {
		lo.innerMax = 1
	}
	return lo
}

// fhit is a clickable footer button rectangle (terminal cells).
type fhit struct {
	x, y, w int
	key     string
}

// footerHits returns clickable footer button rects, or nil when the footer
// is truncated (geometry would be untrustworthy).
func (m Model) footerHits(lo layout) []fhit {
	segs := m.footerSegs()
	if len(segs) == 0 || m.width <= 0 {
		return nil
	}
	sepW := lipgloss.Width(m.footerSep())
	keysW := 0
	for i, s := range segs {
		keysW += lipgloss.Width(s.plain())
		if i+1 < len(segs) {
			keysW += sepW
		}
	}
	if keysW > m.width {
		return nil
	}
	// The footer block is centered by its widest row (keys or message).
	msgW := 0
	if m.err != "" {
		msgW = len("! " + m.err)
	} else if m.info != "" {
		msgW = len(m.info)
	}
	if msgW > m.width {
		return nil
	}
	off := (m.width - max(keysW, msgW)) / 2
	var out []fhit
	x := off
	for i, s := range segs {
		w := lipgloss.Width(s.plain())
		if s.key != "" {
			out = append(out, fhit{x: x, y: lo.footerY, w: w, key: s.key})
		}
		x += w
		if i+1 < len(segs) {
			x += sepW
		}
	}
	return out
}

// shit is a clickable stage-pipeline step rectangle.
type shit struct {
	x, w  int
	stage Stage
	past  bool
}

// stageBarHits returns clickable pipeline steps (visited stages only).
// Narrow terminals render a single summary line with no targets. When the
// full bar overflows the terminal, only the visible prefix cells (matching
// the renderer's own truncation) yield hit rects.
func (m Model) stageBarHits(lo layout) []shit {
	if m.width < 100 {
		return nil
	}
	order := stageOrder()
	cur := 0
	for i, s := range order {
		if s == m.stage {
			cur = i
		}
	}
	sepW := lipgloss.Width(" → ")
	total := 0
	labels := make([]string, len(order))
	for i, s := range order {
		labels[i] = s.label()
		total += lipgloss.Width(labels[i])
		if i+1 < len(order) {
			total += sepW
		}
	}
	x := (m.width - total) / 2
	clip := total
	if total > m.width {
		// Mirrors fitWidthLines: visible prefix keeps untruncated offsets,
		// the last cell is the ellipsis (never clickable).
		x, clip = 0, m.width-1
	}
	var out []shit
	for i, s := range order {
		w := lipgloss.Width(labels[i])
		if x >= 0 && x+w <= clip {
			out = append(out, shit{x: x, w: w, stage: s, past: i < cur})
		}
		x += w
		if i+1 < len(order) {
			x += sepW
		}
	}
	return out
}

// setPress records visual click feedback (hover/pressed highlight).
func (m *Model) setPress(id string, row int) {
	m.pressID = id
	m.pressRow = row
	m.pressAt = time.Now()
}

// pressFresh reports whether pressed feedback is still visible.
func (m Model) pressFresh() bool {
	return !m.pressAt.IsZero() && time.Since(m.pressAt) < 350*time.Millisecond
}

// pressBtn reports whether a footer button is in pressed state.
func (m Model) pressBtn(key string) bool {
	return m.pressFresh() && m.pressID == "btn:"+key
}

// pressRowAt reports whether a content row is in pressed state.
func (m Model) pressRowAt(row int) bool {
	return m.pressFresh() && m.pressID == "row" && m.pressRow == row
}

func inRect(x, y int, rx, ry, rw int) bool {
	return y == ry && x >= rx && x < rx+rw
}

// clickAt routes a left click by exact terminal coordinates: footer
// buttons, pipeline steps, then card content rows.
func (m Model) clickAt(x, y int) (tea.Model, tea.Cmd) {
	if m.stage == StageIntro {
		m.skipIntro()
		return m, nil
	}
	if m.width <= 0 || m.height <= 0 {
		return m, nil
	}
	lo := m.layout()
	for _, h := range m.footerHits(lo) {
		if inRect(x, y, h.x, h.y, h.w) {
			m.setPress("btn:"+h.key, -1)
			return m.clickKey(h.key)
		}
	}
	for _, s := range m.stageBarHits(lo) {
		if s.past && inRect(x, y, s.x, lo.barY, s.w) {
			m.setPress("stage", -1)
			m.stage = s.stage
			m.onEnterStage()
			return m, nil
		}
	}
	if out, cmd, ok := m.clickCard(x, y, lo); ok {
		return out, cmd
	}
	return m, nil
}

// clickKey fires a footer button action through the existing keyboard
// handlers, so mouse and keyboard can never disagree.
func (m Model) clickKey(key string) (tea.Model, tea.Cmd) {
	if key == "ctrl+c" {
		return m, tea.Quit
	}
	zero := tea.KeyMsg{}
	switch m.stage {
	case StageDoctor:
		return m.updateDoctor(key)
	case StageHistory:
		return m.updateHistory(key, zero)
	case StageRunner:
		return m.updateRunner(key)
	case StageModel:
		if key == "]" {
			tabs := catalog.Categories()
			m.tabCursor = (m.tabCursor + 1) % len(tabs)
			m.modelCursor = 0
			m.refreshModels()
			return m, nil
		}
		return m.updateModel(key, zero)
	case StageEffort:
		return m.updateEffort(key)
	case StageWorkspace:
		return m.updateWorkspace(key, zero)
	case StageAgents, StageSkillsPick, StagePlugins, StageHooks, StageMCP:
		m.ensureLists()
		v := m.activeList()
		if v == nil {
			return m, nil
		}
		switch key {
		case " ":
			v.Toggle()
			return m, nil
		case "a":
			v.ToggleAllFiltered()
			return m, nil
		case "c":
			m.burst()
			m.advance()
			return m, nil
		default:
			return m.updateVList(key, zero)
		}
	case StageSkills:
		return m.updateSkills(key)
	case StageLaunch:
		return m.updateLaunch(key)
	}
	return m, nil
}

// modelTabSegs returns the plain tab labels rendered on the model row,
// shared by the view and the hit-tester.
func modelTabSegs() []string {
	tabs := catalog.Categories()
	out := make([]string, len(tabs))
	for i, t := range tabs {
		out[i] = fmt.Sprintf("%d %s", i+1, string(t))
	}
	return out
}

// modelWindow mirrors viewModel's visible window computation.
func (m Model) modelWindow() (start, shown int) {
	shown = min(len(m.filtered), 10)
	start = 0
	if m.modelCursor >= shown {
		start = m.modelCursor - shown + 1
	}
	return start, shown
}

// effortPlain mirrors viewEffort's slider cell text (unstyled).
func (m Model) effortPlain(i int) string {
	e := m.efforts[i]
	compat := true
	if m.hasModel {
		compat = catalog.EffortCompatible(m.curModel, e)
	}
	switch {
	case i == m.effortCursor && compat:
		return e.Label
	case i == m.effortCursor && !compat:
		return e.Label + " LOCKED"
	case !compat:
		return e.Label + " ✕"
	default:
		return e.Label
	}
}

// clickCard handles clicks on card content rows. Coordinates are exact
// terminal cells; rows follow each stage view's line layout.
func (m Model) clickCard(x, y int, lo layout) (tea.Model, tea.Cmd, bool) {
	x0 := lo.boxX + 3 // border (1) + left padding (2)
	if x < x0 || x >= x0+lo.contentW {
		return m, nil, false
	}
	rel := y - (lo.cardY + 2) // border top (1) + top padding (1)
	if rel < 0 || rel >= lo.innerMax {
		return m, nil, false
	}
	switch m.stage {
	case StageRunner:
		return m.clickRunnerRow(rel)
	case StageHistory:
		return m.clickHistoryRow(rel)
	case StageModel:
		return m.clickModelRow(x, x0, rel)
	case StageEffort:
		return m.clickEffortSeg(x, x0, rel)
	case StageAgents, StageSkillsPick, StagePlugins, StageHooks, StageMCP:
		return m.clickVListRow(rel)
	case StageWorkspace:
		if rel == 3 {
			m.wsInput.Focus()
			return m, nil, true
		}
	}
	return m, nil, false
}

// clickRunnerRow selects the runner whose two-line entry was clicked.
func (m Model) clickRunnerRow(rel int) (tea.Model, tea.Cmd, bool) {
	if rel < 3 {
		return m, nil, false
	}
	k := (rel - 3) / 3
	if (rel-3)%3 > 1 || k < 0 || k >= len(catalog.Runners()) {
		return m, nil, false
	}
	m.runnerCursor = k
	m.setPress("row", k)
	out, cmd := m.updateRunner("enter")
	return out, cmd, true
}

// clickHistoryRow resumes the clicked session, or starts a new project.
func (m Model) clickHistoryRow(rel int) (tea.Model, tea.Cmd, bool) {
	zero := tea.KeyMsg{}
	n := min(len(m.sessions), 5)
	if n == 0 {
		return m, nil, false
	}
	if rel >= 3 && rel < 3+n {
		m.histCursor = rel - 3
		m.setPress("row", m.histCursor)
		out, cmd := m.updateHistory("enter", zero)
		return out, cmd, true
	}
	if rel == 3+n+1 {
		m.setPress("row", n)
		m.burst()
		m.stage = StageRunner
		m.onEnterStage()
		return m, nil, true
	}
	return m, nil, false
}

// clickModelRow switches tabs from the tab row or selects a model entry.
// Model entries occupy two physical lines (summary + detail).
func (m Model) clickModelRow(x, x0, rel int) (tea.Model, tea.Cmd, bool) {
	zero := tea.KeyMsg{}
	if rel == 2 && len(m.filtered) > 0 {
		xx := x0
		for i, lbl := range modelTabSegs() {
			w := lipgloss.Width(lbl)
			if x >= xx && x < xx+w {
				m.tabCursor = i
				m.modelCursor = 0
				m.refreshModels()
				m.setPress("row", -2)
				return m, nil, true
			}
			xx += w + 2
		}
		return m, nil, false
	}
	base := 5
	if m.liveLoaded && len(m.liveMods) > 0 {
		base = 6
	}
	if rel < base || len(m.filtered) == 0 {
		return m, nil, false
	}
	start, shown := m.modelWindow()
	idx := start + (rel-base)/2
	if idx < start+shown && idx < len(m.filtered) {
		m.modelCursor = idx
		m.setPress("row", idx)
		out, cmd := m.updateModel("enter", zero)
		return out, cmd, true
	}
	return m, nil, false
}

// clickEffortSeg moves the slider to the clicked effort segment.
func (m Model) clickEffortSeg(x, x0, rel int) (tea.Model, tea.Cmd, bool) {
	if rel != 3 {
		return m, nil, false
	}
	sepW := lipgloss.Width(" ── ")
	xx := x0
	for i := range m.efforts {
		w := lipgloss.Width(m.effortPlain(i))
		if x >= xx && x < xx+w {
			m.effortCursor = i
			m.setPress("row", i)
			return m, nil, true
		}
		xx += w + sepW
	}
	return m, nil, false
}

// clickVListRow toggles the clicked checkbox row.
func (m Model) clickVListRow(rel int) (tea.Model, tea.Cmd, bool) {
	m.ensureLists()
	v := m.activeList()
	if v == nil || rel < 4 {
		return m, nil, false
	}
	vis := v.Visible()
	if r := rel - 4; r >= 0 && r < len(vis) {
		v.Cursor = v.Offset + r
		v.Toggle()
		m.setPress("row", v.Cursor)
		return m, nil, true
	}
	return m, nil, false
}
