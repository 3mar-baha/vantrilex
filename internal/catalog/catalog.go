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