package catalog

import "testing"

func TestZenIncludedForOpenCode(t *testing.T) {
	got := FilterModels(RunnerOpenCode, CatAll, "")
	foundZen := false
	for _, m := range got {
		if m.Zen {
			foundZen = true
		}
		if !m.SupportsRunner(RunnerOpenCode) {
			t.Fatalf("model %s not compatible with opencode", m.ID)
		}
	}
	if !foundZen {
		t.Fatal("expected Zen catalog entries for opencode")
	}
	if len(ZenModels()) < 6 {
		t.Fatalf("expected >=6 zen models, got %d", len(ZenModels()))
	}
}

func TestEffortGating(t *testing.T) {
	var noReason, reason Model
	for _, m := range Models() {
		if m.Reasoning && reason.ID == "" {
			reason = m
		}
		if !m.Reasoning && noReason.ID == "" {
			noReason = m
		}
	}
	max := Efforts()[len(Efforts())-1]
	if !EffortCompatible(reason, max) {
		t.Fatal("reasoning model should allow max effort")
	}
	if EffortCompatible(noReason, max) {
		t.Fatalf("non-reasoning model %s should lock max effort", noReason.ID)
	}
	low := Efforts()[0]
	if !EffortCompatible(noReason, low) {
		t.Fatal("low effort must always be allowed")
	}
}

func TestFuzzySearch(t *testing.T) {
	got := FilterModels(RunnerOpenCode, CatAll, "zen kimi")
	if len(got) == 0 {
		t.Fatal("expected fuzzy match for 'zen kimi'")
	}
}
