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
		case StageHistory:
			return fmt.Sprintf("%s resume  %s new  %s del", k("[1-5]"), k("[N]"), k("[D]"))
		case StageRunner:
			return fmt.Sprintf("%s select  %s back", k("[1-3]"), k("[B]"))
		case StageModel:
			return fmt.Sprintf("%s select  %s tabs  %s back", k("[Enter]"), k("[[]]"), k("[B]"))
		case StageEffort:
			return fmt.Sprintf("%s move  %s ok  %s back", k("[←→]"), k("[Enter]"), k("[B]"))
		case StageWorkspace:
			return fmt.Sprintf("%s ok  %s back", k("[Enter]"), k("[Esc]"))
		case StageAgents, StageSkillsPick, StagePlugins, StageHooks, StageMCP:
			return fmt.Sprintf("%s toggle  %s all  %s confirm  %s search", k("[Space]"), k("[A]"), k("[C]"), k("[/]"))
		case StageSkills:
			return fmt.Sprintf("%s provision  %s skip", k("[Enter]"), k("[S]"))
		case StageLaunch:
			return fmt.Sprintf("%s launch  %s back", k("[Enter]"), k("[B]"))
		}
		return ""
	}
	switch m.stage {
	case StageDoctor:
		return fmt.Sprintf("%s approve + auto-install   %s continue   %s rescan   %s sync toolkit   %s quit", k("[A]"), k("[C]"), k("[R]"), k("[S]"), k("[Ctrl+C]"))
	case StageHistory:
		return fmt.Sprintf("%s resume   %s new project   %s delete   %s navigate   %s quit", k("[1-5]/Enter"), k("[N]"), k("[D]"), k("[↑↓]"), k("[Ctrl+C]"))
	case StageRunner:
		return fmt.Sprintf("%s select   %s navigate   %s back", k("[1-3]/Enter"), k("[↑↓]"), k("[B]"))
	case StageModel:
		return fmt.Sprintf("%s select   %s tabs   %s search   %s back", k("[Enter]"), k("[ [ ] ]"), k("type or [/]"), k("[B/Esc]"))
	case StageEffort:
		return fmt.Sprintf("%s move slider   %s confirm   %s back", k("[←→]"), k("[Enter]"), k("[B]"))
	case StageWorkspace:
		return fmt.Sprintf("%s confirm   %s back", k("[Enter]"), k("[Esc]"))
	case StageAgents, StageSkillsPick, StagePlugins, StageHooks, StageMCP:
		return fmt.Sprintf("%s toggle   %s toggle-all-filtered   %s confirm   %s search   %s back", k("[Space]"), k("[A]"), k("[C]"), k("[/]"), k("[B/Esc]"))
	case StageSkills:
		return fmt.Sprintf("%s provision + continue   %s skip if ready   %s back", k("[Enter]"), k("[S]"), k("[B]"))
	case StageLaunch:
		return fmt.Sprintf("%s launch now   %s open folder   %s back", k("[Enter]"), k("[O]"), k("[B]"))
	}
	return ""
}

// ---- Stage views ----

func (m Model) viewDoctor() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("STAGE 0 — PREFLIGHT DEPENDENCY DOCTOR"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Scanning host for runtime tools. Press A to auto-install missing tools (confirmation required, non-destructive)."))
	b.WriteString("\n")
	b.WriteString(subStyle.Render("Runners auto-update in background  " + m.runnerUpd.Pill()))
	b.WriteString("\n\n")
	if !m.doctorCheck {
		b.WriteString(subStyle.Render("Scanning environment..."))
		return b.String()
	}
	for _, s := range m.statuses {
		mark := pill(s.Found, s.Dep.Optional)
		ver := s.Version
		if !s.Found {
			ver = s.Err
		}
		cmds := ""
		if !s.Found {
			cmds = mutedStyle.Render("  → " + strings.Join(doctorInstallHint(s.Dep.Key), " | "))
		}
		b.WriteString(fmt.Sprintf("%s  %s  %s%s\n", mark, whiteStyle.Render(s.Dep.Label), mutedStyle.Render(ver), cmds))
	}
	miss := m.doctorMissingRequired()
	b.WriteString("\n")
	if len(miss) == 0 {
		b.WriteString(okPill.Render("READY") + "  " + whiteStyle.Render("All required tools present."))
	} else {
		b.WriteString(missPill.Render(fmt.Sprintf("%d MISSING", len(miss))) + "  " + whiteStyle.Render("Press A to install automatically, or C to continue anyway."))
	}
	if len(m.doctorLog) > 0 {
		b.WriteString("\n" + mutedStyle.Render(strings.Join(tail(m.doctorLog, 4), "\n")))
	}
	// Toolkit status (all designated repos, dynamic count).
	root, present, err := toolkitStatusSafe()
	if err == nil {
		b.WriteString("\n\n" + subStyle.Render(fmt.Sprintf("Toolkit (%d): "+root, len(doctor.ToolkitRepos))))
		names := []string{}
		for _, r := range doctor.ToolkitRepos {
			if present[r.Dir] {
				names = append(names, okPill.Render(r.Dir))
			} else {
				names = append(names, missPill.Render(r.Dir))
			}
		}
		b.WriteString("\n" + strings.Join(names, " "))
	}
	if len(m.toolkitLog) > 0 {
		b.WriteString("\n" + mutedStyle.Render(strings.Join(tail(m.toolkitLog, 4), "\n")))
	}
	if m.toolkitBusy {
		b.WriteString("\n" + subStyle.Render("Syncing toolkit (shallow clone, parallel)..."))
	}
	if m.doctorBusy {
		b.WriteString("\n" + subStyle.Render("Installing... (spinner running, do not close)"))
	}
	return b.String()
}

func doctorInstallHint(key string) []string {
	switch key {
	case "git":
		return []string{"winget install -e --id Git.Git"}
	case "node":
		return []string{"winget install -e --id OpenJS.NodeJS.LTS"}
	case "chafa":
		return []string{"winget install -e --id HansPetterJansson.Chafa"}
	case "claude":
		return []string{"npm install -g @anthropic-ai/claude-code"}
	case "opencode":
		return []string{"npm install -g opencode-ai"}
	case "codex":
		return []string{"npm install -g @openai/codex"}
	default:
		return []string{"system package manager"}
	}
}

func tail(lines []string, n int) []string {
	if len(lines) <= n {
		return lines
	}
	return lines[len(lines)-n:]
}

func (m Model) viewHistory() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("STAGE 1 — HISTORY vs NEW PROJECT"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Resume a recent session [1-5], delete with [D], or press [N] for a new project."))
	b.WriteString("\n\n")
	n := min(len(m.sessions), 5)
	if n == 0 {
		b.WriteString(mutedStyle.Render("No previous sessions. Press [N] to create a new project."))
		return b.String()
	}
	for i := 0; i < n; i++ {
		s := m.sessions[i]