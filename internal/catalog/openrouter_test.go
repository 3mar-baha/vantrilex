package catalog

import (
	"strings"
	"testing"
)

func TestModelTags(t *testing.T) {
	cases := []struct {
		model Model
		want  []string
	}{
		{Model{ID: "deepseek/deepseek-r1", Short: "R1", Blurb: "reasoning", Reasoning: true, InputPerM: "$2.00", OutputPerM: "$6.00", Latency: "MED"}, []string{"[REASONING]"}},
		{Model{ID: "meta/llama-4-maverick-vision", Short: "Maverick", Blurb: "vision", InputPerM: "$0.40", OutputPerM: "$0.80", Latency: "LOW"}, []string{"[VISION]", "[ULTRA FAST]"}},
		{Model{ID: "test/coder-flash", Short: "Coder Flash", Blurb: "code generation fast", InputPerM: "$0.50", OutputPerM: "$1.50", Latency: "LOW"}, []string{"[CODING]", "[ULTRA FAST]"}},
		{Model{ID: "x/free-model", Short: "Free", Blurb: "experimentation", InputPerM: "$0.00", OutputPerM: "$0.00", Latency: "MED"}, []string{"[FREE TIER]"}},
	}
	for _, c := range cases {
		got := ModelTags(c.model)
		for _, w := range c.want {
			found := false
			for _, g := range got {
				if g == w {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("model %s: want tag %s in %v", c.model.ID, w, got)
			}
		}
	}
}

func TestParseOpenRouter(t *testing.T) {
	raw := `{"data":[{"id":"test/alpha","name":"Alpha","description":"fast coder","context_length":128000,"pricing":{"prompt":"0.0000005","completion":"0.0000015"}},{"id":"test/free","name":"Free","description":"vision experiment","context_length":64000,"pricing":{"prompt":"0","completion":"0"}}]}`
	models, err := ParseOpenRouter([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if len(models) != 2 {
		t.Fatalf("want 2 models, got %d", len(models))
	}
	if models[0].InputPerM != "$0.50" {
		t.Fatalf("price got %q", models[0].InputPerM)
	}
	if models[0].Context != "128k" {
		t.Fatalf("context got %q", models[0].Context)
	}
	tags := ModelTags(models[1])
	joined := strings.Join(tags, " ")
	if !strings.Contains(joined, "[FREE TIER]") || !strings.Contains(joined, "[VISION]") {
		t.Fatalf("free vision tags missing: %v", tags)
	}
}

func TestMergeModels(t *testing.T) {
	static := []Model{{ID: "a/x", Short: "X", Runners: []RunnerID{RunnerOpenCode}, Zen: true, InputPerM: "$1.00", OutputPerM: "$3.00", Context: "200k"}}
	live := []Model{{ID: "a/x", Short: "X live", InputPerM: "$2.00", OutputPerM: "$4.00", Context: "256k", Blurb: "live"}}
	merged := MergeModels(static, live)
	if len(merged) != 1 {
		t.Fatalf("want 1 merged, got %d", len(merged))
	}
	if merged[0].InputPerM != "$2.00" || !merged[0].Zen {
		t.Fatalf("merge failed: %+v", merged[0])
	}
}
