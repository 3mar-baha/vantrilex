// Package catalog defines runners, model matrix, Zen catalog and effort levels.
package catalog

import "strings"

// RunnerID identifies an agent runner.
type RunnerID string

const (
	RunnerClaude  RunnerID = "claude"
	RunnerOpenCode RunnerID = "opencode"
	RunnerCodex   RunnerID = "codex"
)

// Runner describes an agent harness.
type Runner struct {
	ID    RunnerID
	Index int
	Name  string
	Tag   string
	Desc  string
	Bin   string // binary to invoke
}

var runners = []Runner{
	{ID: RunnerClaude, Index: 1, Name: "Claude Code CLI", Tag: "ANTHROPIC", Desc: "Anthropic agent harness, streaming tools, sub-agents.", Bin: "claude"},
	{ID: RunnerOpenCode, Index: 2, Name: "OpenCode Agent", Tag: "SST + ZEN", Desc: "SST autonomous agent with multi-provider routing and OpenCode Zen.", Bin: "opencode"},
	{ID: RunnerCodex, Index: 3, Name: "OpenAI Codex CLI", Tag: "OPENAI", Desc: "Autonomous agent for GPT-5.6 / Codex execution.", Bin: "codex"},
}

// Runners returns all runners.
func Runners() []Runner { return runners }

// RunnerByIndex 1-based lookup.
func RunnerByIndex(i int) (Runner, bool) {
	for _, r := range runners {
		if r.Index == i {
			return r, true
		}
	}
	return Runner{}, false
}

// RunnerByID lookup.
func RunnerByID(id RunnerID) (Runner, bool) {
	for _, r := range runners {
		if r.ID == id {
			return r, true
		}
	}
	return Runner{}, false
}

// Category tabs for the model matrix.
type Category string

const (
	CatFrontier Category = "Frontier Reasoning"
	CatFast     Category = "Ultra-Fast Coding"
	CatVision   Category = "Multimodal Vision"
	CatFree     Category = "Free Tier"
	CatAll      Category = "All Models"
)

// Categories in tab order.
func Categories() []Category {
	return []Category{CatAll, CatFrontier, CatFast, CatVision, CatFree}
}

// Model describes one selectable model.
type Model struct {
	ID         string   // provider-qualified id passed to runner
	Short      string   // display name
	Category   Category
	Runners    []RunnerID // empty = all runners
	Zen        bool       // part of OpenCode Zen catalog
	InputPerM  string     // pricing input per 1M
	OutputPerM string     // pricing output per 1M
	Context    string
	Latency    string // "LOW", "MED", "HIGH"
	Reasoning  bool   // supports extended reasoning tokens