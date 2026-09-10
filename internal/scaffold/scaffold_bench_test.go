package scaffold

import "testing"

func BenchmarkPlan(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Plan()
	}
}

func BenchmarkBridgeRelPath(b *testing.B) {
	t := Toolkit{Slug: "vercel-skills", Name: "vercel-labs/skills"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BridgeRelPath(t)
	}
}

func BenchmarkSanitize(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sanitize("Lead System Architect")
	}
}
