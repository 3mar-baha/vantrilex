package ui

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"vantrilex/internal/catalog"
	"vantrilex/internal/doctor"
	"vantrilex/internal/particles"
	"vantrilex/internal/runner"
	"vantrilex/internal/scaffold"
)

// Stage is the wizard step.
type Stage int

const (
	StageIntro Stage = iota
	StageDoctor
	StageHistory
	StageRunner
	StageModel
	StageEffort
	StageWorkspace
	StageAgents
	StageSkillsPick
	StagePlugins
	StageHooks
	StageMCP
	StageSkills
	StageLaunch
)

func (s Stage) label() string {
	switch s {
	case StageIntro:
		return "Ignition"
	case StageDoctor:
		return "Preflight Doctor"
	case StageHistory:
		return "History"
	case StageRunner:
		return "Runner"
	case StageModel:
		return "Model Matrix"
	case StageEffort:
		return "Cognitive Effort"
	case StageWorkspace:
		return "Workspace"
	case StageAgents:
		return "Agents"
	case StageSkillsPick:
		return "Skills"
	case StagePlugins:
		return "Plugins"
	case StageHooks:
		return "Hooks"
	case StageMCP:
		return "MCP Servers"
	case StageSkills:
		return "Skills Gate"
	case StageLaunch:
		return "Launch"
	default:
		return ""
	}
}

func stageOrder() []Stage {
	return []Stage{StageDoctor, StageHistory, StageRunner, StageModel, StageEffort, StageWorkspace, StageAgents, StageSkillsPick, StagePlugins, StageHooks, StageMCP, StageSkills, StageLaunch}
}

// Messages.
type tickMsg time.Time
type doctorResultMsg struct{ statuses []doctor.Status }
type installDoneMsg struct {
	key string
	err error
}
type toolkitLineMsg struct{ line string }
type toolkitDoneMsg struct{ err error }
type logoMsg struct{ art string }
type introArtMsg struct{ art string }
type runnerUpdateMsg struct {
	key string
	err error
}
type runnersDoneMsg struct{ errs map[string]error }
type liveModelsMsg struct{ models []catalog.Model }

func tickCmd() tea.Cmd {
	return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Model is the root Bubble Tea model.
type Model struct {
	stage  Stage
	width  int
	height int

	engine *particles.Engine
	logo   string

	// cinematic intro
	introStart time.Time
	introBurst bool
	introArt   string
	warp       []warpStar
	warpRng    *rand.Rand

	// doctor
	statuses    []doctor.Status
	doctorCheck bool
	doctorBusy  bool
	doctorLog   []string
	toolkitLog  []string
	toolkitBusy bool
	toolkitDone bool

	// history
	sessions   []runner.Session
	histCursor int

	// runner
	runnerCursor int
	curRunner    catalog.Runner
	hasRunner    bool

	// model matrix
	tabCursor   int
	modelCursor int
	search      string
	filtered    []catalog.Model
	curModel    catalog.Model
	hasModel    bool

	// effort
	efforts      []catalog.Effort
	effortCursor int

	// workspace
	wsInput   textinput.Model
	wsAsk     bool // asking create folder Y/N
	wsErr     string
	workspace string

	// skills
	scaffoldItems []catalog.Model // unused placeholder to keep struct tidy
	skillItems    []scaffold.Item
	skillDone     bool

	// enterprise registry multi-select stages (virtualized)
	agentsList  VList
	skillsList  VList
	pluginsList VList
	hooksList   VList
	mcpList     VList
	listsInit   bool

	// mouse + background services
	mouseX     int
	mouseY     int
	runnerUpd  doctor.RunnerUpdateStatus
	liveMods   []catalog.Model
	liveLoaded bool

	info string
	err  string

	// launch handover
	pendingLaunch bool
	launched      bool
}

// NewModel builds the initial model.
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "C:\\projects\\my-app  or  ~/projects/my-app"
	ti.CharLimit = 512
	ti.Width = 60

	efforts := catalog.Efforts()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	m := Model{
		stage:        StageIntro,
		engine:       particles.New(100, 24),
		introStart:   time.Now(),
		warp:         newWarpField(90, rng),
		warpRng:      rng,
		runnerCursor: 0,
		wsInput:      ti,
		efforts:      efforts,
		effortCursor: 1,
		width:        100,
		height:       30,
		runnerUpd:    doctor.RunnerUpdateStatus{Updating: true},
	}
	m.refreshModels()
	return m
}

// PendingLaunch reports whether main should exec the runner after quit.
func (m Model) PendingLaunch() bool { return m.pendingLaunch }

// LaunchSpec returns runner/model/effort/workspace for main to exec.
func (m Model) LaunchSpec() (catalog.Runner, string, string, string) {
	return m.curRunner, m.curModel.ID, m.efforts[m.effortCursor].ID, m.workspace
}

// Init kicks off 60 FPS ticks, the doctor scan, background runner
// self-updates, live OpenRouter ingestion, and logo renders.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		func() tea.Msg {
			return doctorResultMsg{statuses: doctor.CheckAll()}
		},
		func() tea.Msg {
			return logoMsg{art: loadLogo()}
		},
		func() tea.Msg {
			return introArtMsg{art: loadLogoSized(36, 18)}
		},
		m.updateRunnersCmd(),
		m.fetchLiveModelsCmd(),
	)
}

// updateRunnersCmd upgrades all three runners in the background.
func (m Model) updateRunnersCmd() tea.Cmd {
	return func() tea.Msg {
		errs := doctor.UpdateAllRunners(func(string) {})
		return runnersDoneMsg{errs: errs}
	}
}

// fetchLiveModelsCmd ingests the live OpenRouter catalog once per session.
func (m Model) fetchLiveModelsCmd() tea.Cmd {
	return func() tea.Msg {
		models := catalog.FetchOpenRouterModels(context.Background())
		return liveModelsMsg{models: models}
	}
}

// introArtCmd renders the hero splash emblem at the strict 2:1 size 36x18.
func introArtCmd() tea.Cmd {
	return func() tea.Msg {
		return introArtMsg{art: loadLogoSized(36, 18)}
	}
}

func (m *Model) refreshModels() {
	tabs := catalog.Categories()
	if m.tabCursor < 0 {
		m.tabCursor = 0
	}
	if m.tabCursor >= len(tabs) {
		m.tabCursor = len(tabs) - 1
	}
	rid := m.curRunner.ID
	if !m.hasRunner {
		rid = catalog.RunnerOpenCode // default broadest catalog for preview
	}
	base := catalog.FilterModels(rid, tabs[m.tabCursor], m.search)
	if m.liveLoaded && len(m.liveMods) > 0 {
		// Overlay live pricing/context onto the filtered static matrix.
		byLive := map[string]catalog.Model{}
		for _, l := range m.liveMods {
			byLive[l.ID] = l
		}
		for i, item := range base {
			if l, ok := byLive[item.ID]; ok {
				if l.InputPerM != "" && l.InputPerM != "varies" {
					base[i].InputPerM = l.InputPerM
				}
				if l.OutputPerM != "" && l.OutputPerM != "varies" {
					base[i].OutputPerM = l.OutputPerM
				}
				if l.Context != "" && l.Context != "n/a" {
					base[i].Context = l.Context
				}
			}
		}
		// Append live-only models matching the current filter is handled by
		// MergeModels at the catalog layer; order/filter stays stable here.
	}
	m.filtered = base
	if m.modelCursor >= len(m.filtered) {
		m.modelCursor = 0
	}
	if len(m.filtered) == 0 {
		m.modelCursor = 0
	}
}

func (m *Model) burst() {
	m.engine.Supernova(90)
}

// enterIntro resets the cinematic sequence.
func (m *Model) enterIntro() {
	m.stage = StageIntro
	m.introStart = time.Now()
	m.introBurst = false
	if m.warpRng == nil {
		m.warpRng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	m.warp = newWarpField(90, m.warpRng)
	m.engine.SetSize(m.width, m.height)
}

// skipIntro jumps straight to the Preflight Doctor.
func (m *Model) skipIntro() {
	m.stage = StageDoctor
	m.onEnterStage()
}

func (m *Model) advance() {
	if m.stage == StageIntro {
		m.skipIntro()
		return
	}
	order := stageOrder()
	for i, s := range order {
		if s == m.stage && i+1 < len(order) {
			m.stage = order[i+1]
			m.onEnterStage()
			return
		}
	}
}

func (m *Model) goBack() {
	order := stageOrder()
	for i, s := range order {
		if s == m.stage && i-1 >= 0 {
			m.stage = order[i-1]
			m.onEnterStage()
			return
		}
	}
}

func (m *Model) onEnterStage() {
	m.err = ""
	m.info = ""
	m.ensureLists()
	// Stage-transition perimeter supernova: shockwave erupts along the full
	// terminal perimeter and travels inwards (skipped for the cinematic
	// intro, which owns its own ignition sequence).
	if m.stage != StageIntro {
		m.engine.PerimeterSupernova()
	}
	switch m.stage {
	case StageIntro:
		m.enterIntro()
	case StageDoctor:
		// Compact menu particle field; header layout owns the height budget.
		m.engine.SetSize(m.width-4, 3)
	case StageHistory:
		m.sessions = runner.LoadSessions()
		if m.histCursor >= len(m.sessions) {
			m.histCursor = 0
		}
	case StageModel:
		m.refreshModels()
	case StageWorkspace:
		m.wsInput.Focus()
		if m.workspace != "" {
			m.wsInput.SetValue(m.workspace)
		}
		m.wsAsk = false
		m.wsErr = ""
	case StageSkills:
		if m.workspace == "" {
			m.skillItems = nil
		} else {
			m.skillItems = scaffold.Check(m.workspace)
		}
		m.skillDone = scaffold.AllPresent(m.workspace) && m.workspace != ""
	}
}

// ensureLists lazily builds the five virtualized registry lists once.
func (m *Model) ensureLists() {
	if m.listsInit {
		return
	}
	m.agentsList = NewVList(catalog.Agents())
	m.skillsList = NewVList(catalog.SkillsRegistry())
	m.pluginsList = NewVList(catalog.Plugins())
	m.hooksList = NewVList(catalog.Hooks())
	m.mcpList = NewVList(catalog.MCPs())
	m.listsInit = true
}

func (m *Model) activeList() *VList {
	switch m.stage {
	case StageAgents:
		return &m.agentsList
	case StageSkillsPick:
		return &m.skillsList
	case StagePlugins:
		return &m.pluginsList
	case StageHooks:
		return &m.hooksList
	case StageMCP:
		return &m.mcpList
	}
	return nil
}

// doctorMissingAny counts absent deps (incl. optional).
func (m Model) doctorMissingAny() []doctor.Status {
	var out []doctor.Status
	for _, s := range m.statuses {
		if !s.Found {
			out = append(out, s)
		}
	}
	return out
}

func (m Model) doctorMissingRequired() []doctor.Status {
	var out []doctor.Status
	for _, s := range m.statuses {
		if !s.Found && !s.Dep.Optional {
			out = append(out, s)
		}
	}
	return out
}

// Update handles all messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.wsInput.Width = min(msg.Width-24, 70)
		if m.stage == StageIntro {
			m.engine.SetSize(msg.Width, msg.Height)
			return m, introArtCmd()
		}
		m.engine.SetSize(msg.Width-4, 3)
		return m, nil

	case tickMsg:
		m.engine.Tick(1.0 / 60.0)
		if m.stage == StageIntro {
			elapsed := time.Since(m.introStart).Seconds()
			// Gentle drift in acts 1-3, hyperspace streaks in act 4.
			accel := 0.5
			if elapsed >= act4At {
				accel = 8.0
			}
			if m.warpRng == nil {
				m.warpRng = rand.New(rand.NewSource(time.Now().UnixNano()))
			}
			advanceWarp(m.warp, 1.0/60.0, accel, m.warpRng)
			if !m.introBurst && elapsed >= act4At {
				m.introBurst = true
				m.engine.Supernova(130)
			}
			if elapsed >= introDuration {
				m.skipIntro()
			}
		}
		return m, tickCmd()

	case logoMsg:
		m.logo = msg.art
		return m, nil

	case introArtMsg:
		m.introArt = msg.art
		return m, nil

	case doctorResultMsg:
		m.statuses = msg.statuses
		m.doctorCheck = true
		m.doctorBusy = false
		return m, nil

	case installDoneMsg:
		m.doctorLog = append(m.doctorLog, fmt.Sprintf("%s: %s", msg.key, doneText(msg.err)))
		// rescan after each install
		return m, func() tea.Msg {
			return doctorResultMsg{statuses: doctor.CheckAll()}
		}

	case toolkitLineMsg:
		m.toolkitLog = append(m.toolkitLog, msg.line)
		if len(m.toolkitLog) > 8 {
			m.toolkitLog = m.toolkitLog[len(m.toolkitLog)-8:]
		}
		return m, nil

	case toolkitDoneMsg:
		m.toolkitBusy = false
		m.toolkitDone = msg.err == nil
		if msg.err != nil {
			m.toolkitLog = append(m.toolkitLog, "toolkit error: "+msg.err.Error())
		} else {
			m.toolkitLog = append(m.toolkitLog, "toolkit sync complete")
		}
		return m, nil

	case runnerUpdateMsg:
		m.doctorLog = append(m.doctorLog, fmt.Sprintf("%s: %s", msg.key, doneText(msg.err)))
		return m, nil

	case runnersDoneMsg:
		m.runnerUpd = doctor.RunnerUpdateStatus{Updating: false, Done: true, Errors: msg.errs}
		if len(msg.errs) > 0 {
			m.doctorLog = append(m.doctorLog, "runner auto-update partial (see pill)")
		} else {
			m.doctorLog = append(m.doctorLog, "runners up-to-date")
		}
		return m, nil

	case liveModelsMsg:
		if len(msg.models) > 0 {
			m.liveMods = msg.models
			m.liveLoaded = true
			if m.stage == StageModel {
				m.refreshModels()
			}
		}
		return m, nil

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Text input routing for workspace stage.
	if m.stage == StageWorkspace && !m.wsAsk {
		var cmd tea.Cmd
		m.wsInput, cmd = m.wsInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func doneText(err error) string {
	if err == nil {
		return "installed"
	}
	return "FAILED: " + err.Error()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// handleMouse routes full mouse interaction: hover trails, click starbursts
// with instant selection/toggling, wheel pagination, search focus.
// Keyboard remains fully functional; mouse is additive.
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	m.mouseX, m.mouseY = msg.X, msg.Y
	switch msg.Action {
	case tea.MouseActionMotion:
		m.engine.EmitTrail(msg.X, msg.Y)
		return m, nil
	case tea.MouseActionPress:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			m.scrollBy(-1)
			return m, nil
		case tea.MouseButtonWheelDown:
			m.scrollBy(1)
			return m, nil
		case tea.MouseButtonRight, tea.MouseButtonMiddle:
			// Right/middle click cycles model tabs on the matrix stage.
			if m.stage == StageModel {
				tabs := catalog.Categories()
				m.tabCursor = (m.tabCursor + 1) % len(tabs)
				m.modelCursor = 0
				m.refreshModels()
			}
			m.engine.ClickBurst(msg.X, msg.Y)
			return m, nil
		default: // left click
			m.engine.ClickBurst(msg.X, msg.Y)
			return m.clickSelect(), nil
		}
	case tea.MouseActionRelease:
		return m, nil
	}
	return m, nil
}

// scrollBy paginates the active list (wheel support).
func (m *Model) scrollBy(delta int) {
	switch m.stage {
	case StageAgents, StageSkillsPick, StagePlugins, StageHooks, StageMCP:
		m.ensureLists()
		if v := m.activeList(); v != nil {
			v.Move(delta)
		}
	case StageModel:
		if len(m.filtered) > 0 {
			m.modelCursor = (m.modelCursor + delta + len(m.filtered)) % len(m.filtered)
		}
	case StageRunner:
		m.runnerCursor = (m.runnerCursor + delta + 3) % 3
	case StageHistory:
		n := min(len(m.sessions), 5)
		if n > 0 {
			m.histCursor = (m.histCursor + delta + n) % n
		}
	}
}

// clickSelect performs instant left-click selection per stage.
func (m Model) clickSelect() tea.Model {
	switch m.stage {
	case StageIntro:
		m.skipIntro()
	case StageRunner:
		idx := m.runnerCursor + 1
		if r, ok := catalog.RunnerByIndex(idx); ok {
			m.curRunner = r
			m.hasRunner = true
			m.burst()
			m.stage = StageModel
			m.onEnterStage()
		}
	case StageModel:
		if len(m.filtered) > 0 {
			m.curModel = m.filtered[m.modelCursor]
			m.hasModel = true
			m.burst()
			m.stage = StageEffort
			m.onEnterStage()
			m.effortCursor = 1
		}
	case StageAgents, StageSkillsPick, StagePlugins, StageHooks, StageMCP:
		m.ensureLists()
		if v := m.activeList(); v != nil {
			v.Toggle()
		}
	case StageWorkspace:
		m.wsInput.Focus()
	case StageLaunch:
		if !m.launched {
			m.burst()
			m.pendingLaunch = true
			m.launched = true
			m.info = "Handing over terminal to runner..."
			// NOTE: cannot return tea.Quit from here without Cmd plumbing;
			// the launch happens on the next Enter or via pendingLaunch.
		}
	}
	return m
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	// Global quit.
	if key == "ctrl+c" {
		return m, tea.Quit
	}
	// Intro skip: any keypress (Space, Enter, Esc, ...) jumps to the Doctor.
	if m.stage == StageIntro {
		m.skipIntro()
		return m, nil
	}
	switch m.stage {
	case StageDoctor:
		return m.updateDoctor(key)
	case StageHistory:
		return m.updateHistory(key, msg)
	case StageRunner:
		return m.updateRunner(key)
	case StageModel:
		return m.updateModel(key, msg)
	case StageEffort:
		return m.updateEffort(key)
	case StageWorkspace:
		return m.updateWorkspace(key, msg)
	case StageAgents, StageSkillsPick, StagePlugins, StageHooks, StageMCP:
		return m.updateVList(key, msg)
	case StageSkills:
		return m.updateSkills(key)
	case StageLaunch:
		return m.updateLaunch(key)
	}
	return m, nil
}

// ---- Stage updates ----

func (m Model) updateDoctor(key string) (tea.Model, tea.Cmd) {
	switch strings.ToLower(key) {
	case "a":
		if m.doctorBusy || m.toolkitBusy {
			return m, nil
		}
		missing := m.doctorMissingAny()
		if len(missing) == 0 {
			m.info = "All dependencies present. Syncing toolkit..."
			return m, m.syncToolkitCmd()
		}
		m.doctorBusy = true
		m.doctorLog = append(m.doctorLog, "Auto-fix approved: installing missing tools...")
		return m, m.installAllCmd(missing)
	case "c":
		m.burst()
		m.stage = StageHistory
		m.onEnterStage()
		return m, nil
	case "r":
		m.doctorBusy = true
		m.info = "Rescanning..."
		return m, func() tea.Msg {
			return doctorResultMsg{statuses: doctor.CheckAll()}
		}
	case "s":
		if m.toolkitBusy {
			return m, nil
		}
		return m, m.syncToolkitCmd()
	}
	return m, nil
}

func (m Model) installAllCmd(missing []doctor.Status) tea.Cmd {
	return func() tea.Msg {
		// Install sequentially, then report last; UI rescans after each via chaining.
		// Run first missing item here; the rest continue through repeated rescan? For