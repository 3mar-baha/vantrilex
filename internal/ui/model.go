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