package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"vantrilex/internal/catalog"
	"vantrilex/internal/doctor"
	"vantrilex/internal/runner"
	"vantrilex/internal/scaffold"
)

func (m Model) renderHeader(logoH, fieldH int) string {
	logo := m.logo
	if logo == "" {
		logo = fallbackLogo()
	}
	logo = fitLines(logo, logoH)
	field := fitLines(m.engine.RenderField(), fieldH)
	return logo + "\n" + field
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) renderStageBar() string {
	stages := stageOrder()
	cur := 0
	for i, s := range stages {
		if s == m.stage {
			cur = i
		}
	}
	// Narrow terminals: single compact tracker line (always exactly 1 row).
	if m.width < 100 {
		return mutedStyle.Render(fmt.Sprintf("STAGE %d/%d · %s", cur+1, len(stages), stages[cur].label()))
	}
	parts := make([]string, 0, len(stages))
	for _, s := range stages {
		lbl := s.label()
		if s == m.stage {
			parts = append(parts, selStyle.Render(lbl))
		} else if s < m.stage {
			parts = append(parts, subStyle.Render(lbl))
		} else {
			parts = append(parts, mutedStyle.Render(lbl))
		}
	}
	return strings.Join(parts, mutedStyle.Render(" → "))
}

// renderCard boxes stage content into exactly cardH rows.
func (m Model) renderCard(cardH int) string {
	if cardH < 3 {
		cardH = 3
	}
	// Box chrome costs 4 rows (top/bottom border + vertical padding).
	innerMax := cardH - 4
	if innerMax < 1 {
		innerMax = 1
	}
	contentW := min(92, m.width-8)
	if contentW < 16 {
		contentW = 16
	}
	raw := m.renderStageContent()
	// Fit loop: shrink content until the boxed block (borders + padding
	// included) fits cardH. This absorbs any lipgloss word-wrap growth
	// from rows whose measured width disagrees with render width.
	boxed := ""
	for n := innerMax; n >= 1; n-- {
		content := fitWidthLines(fitLines(raw, n), contentW)
		boxed = boxStyle.Copy().Width(contentW).Render(content)
		if countLines(boxed) <= cardH {
			break
		}
	}
	return lipgloss.Place(m.width, cardH, lipgloss.Center, lipgloss.Top, fitLines(boxed, cardH))
}

func (m Model) renderStageContent() string {
	var content string
	switch m.stage {
	case StageDoctor:
		content = m.viewDoctor()
	case StageHistory:
		content = m.viewHistory()
	case StageRunner:
		content = m.viewRunner()
	case StageModel:
		content = m.viewModel()
	case StageEffort:
		content = m.viewEffort()
	case StageWorkspace:
		content = m.viewWorkspace()
	case StageAgents:
		content = m.viewVList("AGENTS", "Architecture, Security, Frontend, Backend, DevOps, QA personas.", &m.agentsList)
	case StageSkillsPick:
		content = m.viewVList("SKILLS", "Workflow skills from ECC, mattpocock, ponytail, guard-skills, vercel.", &m.skillsList)
	case StagePlugins:
		content = m.viewVList("PLUGINS", "Official + community plugin marketplace.", &m.pluginsList)
	case StageHooks:
		content = m.viewVList("HOOKS", "Lifecycle hooks: PreToolUse, PostToolUse, PreCompact, Stop.", &m.hooksList)
	case StageMCP:
		content = m.viewVList("MCP SERVERS", "1000+ Model Context Protocol servers.", &m.mcpList)
	case StageSkills:
		content = m.viewSkills()
	case StageLaunch:
		content = m.viewLaunch()
	default:
		content = ""
	}
	return content
}

func (m Model) renderFooter() string {
	keys := m.footerKeys()
	// Footer is at most 2 rows: keys + one message line (error wins).
	msg := ""
	if m.err != "" {
		msg = lipgloss.NewStyle().Foreground(Red).Render("! " + m.err)
	} else if m.info != "" {
		msg = subStyle.Render(m.info)
	}
	if msg == "" {
		return keys
	}
	return keys + "\n" + msg
}

func (m Model) footerKeys() string {
	k := func(keys string) string { return keyStyle.Render(keys) }
	// Narrow terminals get compact single-line hints (width-budgeted).
	if m.width < 100 {
		switch m.stage {
		case StageDoctor:
			return fmt.Sprintf("%s fix  %s go  %s scan  %s sync", k("[A]"), k("[C]"), k("[R]"), k("[S]"))