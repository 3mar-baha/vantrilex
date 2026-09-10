package catalog

import "testing"

func TestRegistryCapacities(t *testing.T) {
	if got := len(MCPs()); got < 1000 {
		t.Fatalf("MCPs = %d, want >=1000", got)
	}
	if got := len(Plugins()); got < 100 {
		t.Fatalf("plugins = %d, want >=100", got)
	}
	if got := len(SkillsRegistry()); got < 300 {
		t.Fatalf("skills = %d, want >=300", got)
	}
	if got := len(Hooks()); got < 50 {
		t.Fatalf("hooks = %d, want >=50", got)
	}
	if got := len(Agents()); got < 300 {
		t.Fatalf("agents = %d, want >=300", got)
	}
}

func TestRegistryDefaults(t *testing.T) {
	for _, want := range []struct {
		items []RegistryItem
		name  string
	}{
		{Agents(), "Lead System Architect"},
		{SkillsRegistry(), "find-skills"},
		{SkillsRegistry(), "skill-creator"},
		{Plugins(), "commit-commands"},
		{Plugins(), "circuit-breaker-guard"},
		{Plugins(), "context-primer"},
		{Hooks(), "pre-compact-checkpoint"},
		{Hooks(), "dangerous-command-guard"},
		{Hooks(), "format-on-edit"},
		{MCPs(), "sequential-thinking"},
		{MCPs(), "filesystem"},
		{MCPs(), "fetch"},
		{MCPs(), "memory"},
	} {
		found := false
		for _, it := range want.items {
			if it.Name == want.name && it.Selected {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("default %q not preselected", want.name)
		}
	}
}

func TestRegistrySearch(t *testing.T) {
	got := SearchRegistry(MCPs(), "filesystem")
	if len(got) == 0 {
		t.Fatal("expected filesystem search hit")
	}
	if got := SearchRegistry(Agents(), "zzz-no-such-asset"); len(got) != 0 {
		t.Fatal("expected empty result for nonsense query")
	}
}
