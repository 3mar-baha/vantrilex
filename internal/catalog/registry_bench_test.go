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
