// Command check verifies embedded registry capacities offline.
package main

import (
	"fmt"
	"os"

	"vantrilex/internal/catalog"
)

func main() {
	targets := []struct {
		name  string
		count int
		want  int
	}{
		{"mcp", len(catalog.MCPs()), 1000},
		{"plugins", len(catalog.Plugins()), 100},
		{"skills", len(catalog.SkillsRegistry()), 300},
		{"hooks", len(catalog.Hooks()), 50},
		{"agents", len(catalog.Agents()), 300},
	}
	failed := false
	for _, t := range targets {
		status := "ok"
		if t.count < t.want {
			status = "SHORT"
			failed = true
		}
		fmt.Printf("%-8s %5d (want >=%d) %s\n", t.name, t.count, t.want, status)
	}
	checkInstallRefs()
	if failed {
		os.Exit(1)
	}
}
