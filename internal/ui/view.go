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
		cursor := "  "
		row := fmt.Sprintf("[%d] %s  %s  %s  %s", i+1,
			runnerBadge(s.Runner), modelTag(s.Model),
			runner.TruncPath(s.Workspace, 34),
			mutedStyle.Render(s.At.Format("2006-01-02 15:04")))
		if i == m.histCursor {
			cursor = keyStyle.Render("▸ ")
			row = selStyle.Render(fmt.Sprintf("[%d] %s %s %s", i+1, s.Runner, shortModel(s.Model), runner.TruncPath(s.Workspace, 30)))
		}
		b.WriteString(cursor + row + "\n")
	}
	b.WriteString("\n" + mutedStyle.Render("[N] New Project"))
	return b.String()
}

func runnerBadge(r string) string {
	return lipgloss.NewStyle().Foreground(Navy).Background(Sky).Bold(true).Padding(0, 1).Render(strings.ToUpper(r))
}

func modelTag(id string) string {
	return lipgloss.NewStyle().Foreground(Cyan).Render(shortModel(id))
}

func shortModel(id string) string {
	if i := strings.LastIndex(id, "/"); i >= 0 {
		return id[i+1:]
	}
	return id
}

func (m Model) viewRunner() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("STAGE 2 — AGENT RUNNER SELECTION"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Choose the harness that will own the terminal session."))
	b.WriteString("\n\n")
	for i, r := range catalog.Runners() {
		marker := "⟦ " + itoa(r.Index) + " ⟧"
		line := fmt.Sprintf("%s  %s  %s\n      %s", marker, whiteStyle.Render(r.Name), runnerBadge(r.Tag), mutedStyle.Render(r.Desc))
		if m.runnerCursor == i {
			line = selStyle.Render(fmt.Sprintf("%s  %s", marker, r.Name)) + "  " + mutedStyle.Render(r.Desc)
		}
		b.WriteString(line + "\n\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m Model) viewModel() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("STAGE 3 — MODEL MATRIX"))
	if m.hasRunner && m.curRunner.ID == catalog.RunnerOpenCode {
		b.WriteString("  " + lipgloss.NewStyle().Foreground(Navy).Background(Cyan).Bold(true).Padding(0, 1).Render("ZEN CATALOG INCLUDED"))
	}
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Fuzzy search: just type, or press [/]. Tabs [1-5] or [ [ ] ]. Pricing per 1M in/out · context · latency."))
	b.WriteString("\n")
	// Tabs.
	tabs := catalog.Categories()
	var tabRow []string
	for i, t := range tabs {
		lbl := fmt.Sprintf("%d %s", i+1, string(t))
		if i == m.tabCursor {
			tabRow = append(tabRow, selStyle.Render(lbl))
		} else {
			tabRow = append(tabRow, mutedStyle.Render(lbl))
		}
	}
	b.WriteString(strings.Join(tabRow, "  ") + "\n")
	b.WriteString(subStyle.Render("Search: "+m.search+"▌") + "\n\n")
	if len(m.filtered) == 0 {
		b.WriteString(mutedStyle.Render("No models match. Clear search with Esc."))
		return b.String()
	}
	if m.liveLoaded && len(m.liveMods) > 0 {
		b.WriteString(subStyle.Render(fmt.Sprintf("Live OpenRouter catalog: %d models merged.", len(m.liveMods))) + "\n")
	}
	shown := min(len(m.filtered), 10)
	start := 0
	if m.modelCursor >= shown {
		start = m.modelCursor - shown + 1
	}
	for i := start; i < start+shown && i < len(m.filtered); i++ {
		md := m.filtered[i]
		lat := latencyBadge(md.Latency)
		zen := ""
		if md.Zen {
			zen = " " + lipgloss.NewStyle().Foreground(Navy).Background(Cyan).Padding(0, 1).Render("ZEN")
		}
		tags := ""
		for _, tg := range catalog.ModelTags(md) {
			tags += " " + tagPill(tg)
		}
		reason := ""
		if !md.Reasoning {
			reason = " " + mutedStyle.Render("(no-reasoning)")
		}
		line := fmt.Sprintf("%s  $%s/$%s  ctx %s  %s%s%s%s\n    %s",
			whiteStyle.Render(md.Short),
			md.InputPerM, md.OutputPerM, md.Context, lat, zen, tags, reason,
			mutedStyle.Render(md.Blurb+" · "+md.ID))
		prefix := "  "
		if i == m.modelCursor {
			prefix = keyStyle.Render("▸ ")
			line = selStyle.Render(md.Short) + mutedStyle.Render(fmt.Sprintf("  $%s/$%s ctx %s", md.InputPerM, md.OutputPerM, md.Context)) + " " + lat + zen + tags
		}
		b.WriteString(prefix + line + "\n")
	}
	b.WriteString(mutedStyle.Render(fmt.Sprintf("\n%d model(s) · showing %d", len(m.filtered), shown)))
	return b.String()
}

func latencyBadge(l string) string {
	switch l {
	case "LOW":
		return okPill.Render("LOW")
	case "MED":
		return optPill.Render("MED")
	default:
		return missPill.Render(l)
	}
}

// tagPill renders automated capability badges in English.
func tagPill(t string) string {
	switch t {
	case "[REASONING]":
		return lipgloss.NewStyle().Foreground(Navy).Background(Cyan).Bold(true).Padding(0, 1).Render("REASONING")
	case "[VISION]":
		return lipgloss.NewStyle().Foreground(White).Background(Violet).Bold(true).Padding(0, 1).Render("VISION")
	case "[CODING]":
		return lipgloss.NewStyle().Foreground(Navy).Background(Sky).Bold(true).Padding(0, 1).Render("CODING")
	case "[FREE TIER]":
		return okPill.Render("FREE TIER")
	case "[ULTRA FAST]":
		return optPill.Render("ULTRA FAST")
	default:
		return mutedStyle.Render(t)
	}
}

func (m Model) viewEffort() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("STAGE 4 — COGNITIVE EFFORT GATING"))
	b.WriteString("\n")
	if m.hasModel {
		b.WriteString(mutedStyle.Render(fmt.Sprintf("Model: %s  (reasoning=%v)", m.curModel.Short, m.curModel.Reasoning)))
	} else {
		b.WriteString(mutedStyle.Render("Model: (none)"))
	}
	b.WriteString("\n\n")
	// Slider.
	var segs []string
	for i, e := range m.efforts {
		compat := true
		if m.hasModel {
			compat = catalog.EffortCompatible(m.curModel, e)
		}
		lbl := e.Label
		var cell string
		switch {
		case i == m.effortCursor && compat:
			cell = selStyle.Render(lbl)
		case i == m.effortCursor && !compat:
			cell = missPill.Render(lbl + " LOCKED")
		case !compat:
			cell = mutedStyle.Render(lbl + " ✕")
		default:
			cell = whiteStyle.Render(lbl)
		}
		segs = append(segs, cell)
	}
	b.WriteString(strings.Join(segs, mutedStyle.Render(" ── ")) + "\n\n")
	e := m.efforts[m.effortCursor]
	b.WriteString(whiteStyle.Render("Selected: "+e.Label) + "  " + mutedStyle.Render(e.Desc) + "\n")
	if m.hasModel && !catalog.EffortCompatible(m.curModel, e) {
		b.WriteString(missPill.Render("LOCKED") + "  " + lipgloss.NewStyle().Foreground(Red).Render(m.curModel.Short+" lacks extended reasoning tokens. Pick low/medium."))
	}
	return b.String()
}

func (m Model) viewWorkspace() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("STAGE 5 — WORKSPACE TARGET DIRECTORY"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Enter the project folder. It will be created on confirm if missing."))
	b.WriteString("\n\n")
	b.WriteString(m.wsInput.View() + "\n\n")
	if m.wsAsk {
		b.WriteString(lipgloss.NewStyle().Foreground(Amber).Bold(true).Render("Create folder? [Y/N]") + "\n")
		b.WriteString(mutedStyle.Render(m.info) + "\n")
	}
	if m.wsErr != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(Red).Render(m.wsErr) + "\n")
	}
	if m.workspace != "" {
		b.WriteString(mutedStyle.Render("Current: "+m.workspace) + "\n")
	}
	return b.String()
}

// viewVList renders a virtualized multi-select stage (~12 visible rows).
func (m Model) viewVList(title, hint string, v *VList) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("STAGE — " + title))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render(hint + "  [/] search · [Space] toggle · [A] all · [C] confirm."))
	b.WriteString("\n")
	b.WriteString(subStyle.Render(fmt.Sprintf("Search: %s▌   Selected: [%d/%d]", v.Search, v.SelectedCount(), len(v.Items))) + "\n\n")
	if len(v.Filtered) == 0 {
		b.WriteString(mutedStyle.Render("No matches. Clear search with Esc."))
		return b.String()
	}
	for row, idx := range v.Visible() {
		it := v.Items[idx]
		box := "[ ]"
		if v.Selected[it.Name] {
			box = "[X]"
		}
		line := fmt.Sprintf("%s  %s  %s", box, whiteStyle.Render(truncV(it.Name, 34)), mutedStyle.Render(truncV(it.Desc, 52)))
		prefix := "  "
		if v.Offset+row == v.Cursor {
			prefix = keyStyle.Render("▸ ")
			line = selStyle.Render(box+" "+truncV(it.Name, 34)) + "  " + mutedStyle.Render(truncV(it.Desc, 52))
		}
		b.WriteString(prefix + line + "\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func truncV(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func (m Model) viewSkills() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("STAGE 6 — WORKFLOW & SKILLS PROVISIONING GATE"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Workspace: " + m.workspace))
	b.WriteString("\n\n")
	items := m.skillItems
	if items == nil && m.workspace != "" {
		items = scaffold.Check(m.workspace)
	}
	for _, it := range items {
		mark := okPill.Render("FOUND")
		if !it.Present {
			mark = missPill.Render("MISSING")
		}
		b.WriteString(fmt.Sprintf("%s  %s  %s\n", mark, whiteStyle.Render(it.Path), mutedStyle.Render(it.Note)))
	}
	b.WriteString("\n" + mutedStyle.Render("External: vercel-labs/skills@find-skills · Anthropic skill-creator") + "\n")
	b.WriteString(mutedStyle.Render(fmt.Sprintf("Toolkit: $HOME/.vantrilex/toolkit (%d repos, shallow clones)", len(doctor.ToolkitRepos))) + "\n")
	if m.skillDone {
		b.WriteString("\n" + okPill.Render("GATE READY") + "  " + whiteStyle.Render("Press S or Enter to continue."))
	} else {
		b.WriteString("\n" + whiteStyle.Render("Press Enter to scaffold missing files (never overwrites)."))
	}
	return b.String()
}

func (m Model) viewLaunch() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("STAGE 7 — SUBPROCESS LAUNCH"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Confirm and hand over the terminal to the runner."))
	b.WriteString("\n\n")
	name, args := runner.BuildCommand(m.curRunner, m.curModel.ID, m.efforts[m.effortCursor].ID, m.workspace)
	rows := []string{
		kv("Runner", string(m.curRunner.ID)+" — "+m.curRunner.Name),
		kv("Model", m.curModel.ID),
		kv("Effort", m.efforts[m.effortCursor].ID),
		kv("Workspace", m.workspace),
		kv("Command", name+" "+strings.Join(args, " ")),
	}
	b.WriteString(strings.Join(rows, "\n") + "\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(Cyan).Bold(true).Render("Press ENTER to launch — the runner takes over this terminal."))
	return b.String()
}

func kv(k, v string) string {
	return subStyle.Render(k+": ") + whiteStyle.Render(v)
}
