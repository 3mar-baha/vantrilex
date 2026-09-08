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
		// OpenCode Zen catalog (routed via opencode/zen).
		{ID: "opencode/zen", Short: "Zen Auto-Router", Category: CatFrontier, Runners: []RunnerID{RunnerOpenCode}, Zen: true, InputPerM: "varies", OutputPerM: "varies", Context: "512k", Latency: "MED", Reasoning: true, Blurb: "Zen smart router, best-value frontier."},
		{ID: "kimi-k2.6", Short: "Kimi K2.6 (Zen)", Category: CatFast, Runners: []RunnerID{RunnerOpenCode}, Zen: true, InputPerM: "$1.20", OutputPerM: "$3.60", Context: "256k", Latency: "LOW", Reasoning: true, Blurb: "Agentic coding specialist."},
		{ID: "qwen3.6-plus", Short: "Qwen 3.6 Plus (Zen)", Category: CatFast, Runners: []RunnerID{RunnerOpenCode}, Zen: true, InputPerM: "$0.90", OutputPerM: "$2.70", Context: "256k", Latency: "LOW", Reasoning: true, Blurb: "Fast open coder."},
		{ID: "minimax-m3", Short: "MiniMax M3 (Zen)", Category: CatVision, Runners: []RunnerID{RunnerOpenCode}, Zen: true, InputPerM: "$1.00", OutputPerM: "$3.00", Context: "200k", Latency: "LOW", Reasoning: true, Blurb: "Vision + tool use."},
		{ID: "deepseek-v4-pro", Short: "DeepSeek V4 Pro (Zen)", Category: CatFrontier, Runners: []RunnerID{RunnerOpenCode}, Zen: true, InputPerM: "$2.00", OutputPerM: "$6.00", Context: "256k", Latency: "MED", Reasoning: true, Blurb: "Deep reasoning at low price."},
		{ID: "glm-5.1", Short: "GLM-5.1 (Zen)", Category: CatFrontier, Runners: []RunnerID{RunnerOpenCode}, Zen: true, InputPerM: "$1.50", OutputPerM: "$4.50", Context: "200k", Latency: "MED", Reasoning: true, Blurb: "Bilingual reasoner."},
		{ID: "deepseek/deepseek-v3-free", Short: "DeepSeek V3 (Free)", Category: CatFree, Runners: []RunnerID{RunnerOpenCode}, InputPerM: "$0.00", OutputPerM: "$0.00", Context: "128k", Latency: "MED", Reasoning: false, Blurb: "Zero-cost experimentation."},
		{ID: "moonshot/kimi-k2-free", Short: "Kimi K2 (Free)", Category: CatFree, Runners: []RunnerID{RunnerOpenCode}, InputPerM: "$0.00", OutputPerM: "$0.00", Context: "128k", Latency: "MED", Reasoning: false, Blurb: "Free agentic runs."},
		{ID: "meta/llama-4-maverick-vision", Short: "Llama 4 Maverick", Category: CatVision, Runners: []RunnerID{RunnerOpenCode}, InputPerM: "$0.40", OutputPerM: "$0.80", Context: "512k", Latency: "LOW", Reasoning: false, Blurb: "Open vision workhorse."},
		{ID: "anthropic/claude-opus-4-6-vision", Short: "Opus 4.6 Vision", Category: CatVision, InputPerM: "$15.00", OutputPerM: "$75.00", Context: "200k", Latency: "HIGH", Reasoning: true, Blurb: "Screenshot + diagram analysis."},
	}
}

// ZenModels returns only Zen catalog entries.
func ZenModels() []Model {
	var out []Model
	for _, m := range Models() {
		if m.Zen {
			out = append(out, m)
		}
	}
	return out
}

// FilterModels applies runner compat + category tab + fuzzy keyword.
func FilterModels(runner RunnerID, cat Category, query string) []Model {