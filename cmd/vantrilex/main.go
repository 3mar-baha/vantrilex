// Command vantrilex: Vantrilex Workflow Launcher — fully autonomous
// self-bootstrapping agent and workflow orchestrator.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"vantrilex/internal/runner"
	"vantrilex/internal/ui"
)

func main() {
	m := ui.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	final, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Vantrilex Workflow Launcher error: "+err.Error())
		os.Exit(1)
	}
	fm, ok := final.(ui.Model)
	if !ok {
		return
	}
	if !fm.PendingLaunch() {
		fmt.Println("Vantrilex Workflow Launcher closed. Goodbye.")
		return
	}
	r, modelID, effort, workspace := fm.LaunchSpec()
	fmt.Printf("Launching %s (%s) in %s ...\n", r.Name, modelID, workspace)
	if err := runner.Launch(r, modelID, effort, workspace); err != nil {
		fmt.Fprintln(os.Stderr, "Launch failed: "+err.Error())
		os.Exit(1)
	}
}
