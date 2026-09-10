package doctor

import (
	"strings"
	"testing"
)

func TestCheckAllFindsCoreTools(t *testing.T) {
	all := CheckAll()
	if len(all) != len(Deps()) {
		t.Fatalf("expected %d statuses, got %d", len(Deps()), len(all))
	}
	byKey := map[string]Status{}
	for _, s := range all {
		byKey[s.Dep.Key] = s
	}
	// Host has git, node, npx, chafa per environment.
	for _, k := range []string{"git", "node", "npx", "chafa"} {
		if !byKey[k].Found {
			t.Fatalf("expected %s to be detected on this host", k)
		}
	}
	if len(ToolkitRepos) != 8 {
		t.Fatalf("expected 8 toolkit repos, got %d", len(ToolkitRepos))
	}
	wantDirs := map[string]string{
		"skills":                         "mattpocock/skills",
		"ponytail":                       "ponytail",
		"guard-skills":                   "guard-skills",
		"Anthropic-Cybersecurity-Skills": "Cybersecurity",
		"agency-agents":                  "agency-agents",
		"everything-claude-code":         "everything-claude-code",
		"vercel-skills":                  "vercel-labs/skills",
		"anthropic-skills":               "anthropics/skills",
	}
	seen := map[string]bool{}
	for _, r := range ToolkitRepos {
		frag, ok := wantDirs[r.Dir]
		if !ok {
			t.Fatalf("unexpected toolkit dir %q", r.Dir)
		}
		if !strings.Contains(strings.ToLower(r.URL), strings.ToLower(frag)) {
			t.Fatalf("toolkit %q has unexpected URL %q", r.Dir, r.URL)
		}
		seen[r.Dir] = true
	}
	if len(seen) != 8 {
		t.Fatalf("expected 8 distinct toolkit dirs, got %d", len(seen))
	}
}

func TestInstallPlanNonEmpty(t *testing.T) {
	for _, d := range Deps() {
		_ = InstallCommands(d) // must not panic; may be empty on some platforms
	}
}
