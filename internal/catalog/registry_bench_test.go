package catalog

import "testing"

func BenchmarkSearchRegistry(b *testing.B) {
	items := MCPs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SearchRegistry(items, "sql database")
	}
}

func BenchmarkRegistryLoad(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if n := len(MCPs()); n < 1000 {
			b.Fatalf("want >=1000 MCPs, got %d", n)
		}
	}
}

func BenchmarkFilterModels(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FilterModels(RunnerOpenCode, CatAll, "zen")
	}
}

func BenchmarkModelTags(b *testing.B) {
	m := Model{ID: "test/coder-flash", Short: "Coder Flash", Blurb: "fast vision coder", InputPerM: "$0.50", OutputPerM: "$1.50", Latency: "LOW", Reasoning: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ModelTags(m)
	}
}
