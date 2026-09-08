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
	Blurb      string
}

// SupportsRunner reports compat.
func (m Model) SupportsRunner(r RunnerID) bool {
	if len(m.Runners) == 0 {
		return true
	}
	for _, x := range m.Runners {
		if x == r {
			return true
		}
	}
	return false
}

// Models returns the full matrix.
func Models() []Model {
	return []Model{
		{ID: "anthropic/claude-opus-4-6", Short: "Claude Opus 4.6", Category: CatFrontier, InputPerM: "$15.00", OutputPerM: "$75.00", Context: "200k", Latency: "HIGH", Reasoning: true, Blurb: "Deepest reasoning, agentic coding."},
		{ID: "anthropic/claude-sonnet-4-5", Short: "Claude Sonnet 4.5", Category: CatFast, InputPerM: "$3.00", OutputPerM: "$15.00", Context: "1M", Latency: "LOW", Reasoning: true, Blurb: "Workhorse coder, fast + smart."},
		{ID: "anthropic/claude-haiku-4-5", Short: "Claude Haiku 4.5", Category: CatFast, InputPerM: "$0.80", OutputPerM: "$4.00", Context: "200k", Latency: "LOW", Reasoning: false, Blurb: "Cheapest, sub-second latency."},
		{ID: "openai/gpt-5.6", Short: "GPT-5.6", Category: CatFrontier, Runners: []RunnerID{RunnerCodex, RunnerOpenCode}, InputPerM: "$10.00", OutputPerM: "$30.00", Context: "400k", Latency: "MED", Reasoning: true, Blurb: "OpenAI flagship reasoner."},
		{ID: "openai/gpt-5.6-codex", Short: "GPT-5.6 Codex", Category: CatFast, Runners: []RunnerID{RunnerCodex}, InputPerM: "$8.00", OutputPerM: "$24.00", Context: "400k", Latency: "LOW", Reasoning: true, Blurb: "Codex-tuned execution variant."},
		{ID: "google/gemini-3-pro", Short: "Gemini 3 Pro", Category: CatVision, InputPerM: "$5.00", OutputPerM: "$15.00", Context: "1M", Latency: "MED", Reasoning: true, Blurb: "Long context + vision."},
		{ID: "google/gemini-3-flash", Short: "Gemini 3 Flash", Category: CatFast, InputPerM: "$0.50", OutputPerM: "$1.50", Context: "1M", Latency: "LOW", Reasoning: false, Blurb: "Ultra-cheap multimodal."},
		{ID: "x-ai/grok-4-1", Short: "Grok 4.1", Category: CatFrontier, Runners: []RunnerID{RunnerOpenCode}, InputPerM: "$6.00", OutputPerM: "$18.00", Context: "256k", Latency: "MED", Reasoning: true, Blurb: "Contrarian reasoner."},