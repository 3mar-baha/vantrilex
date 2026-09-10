package runner

import (
	"testing"

	"vantrilex/internal/catalog"
)

func BenchmarkBuildCommand(b *testing.B) {
	r, _ := catalog.RunnerByID(catalog.RunnerCodex)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = BuildCommand(r, "openai/gpt-5.6", "max", "C:\\tmp\\x")
	}
}

func BenchmarkTruncPath(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TruncPath("C:\\very\\long\\path\\to\\project\\folder\\name", 20)
	}
}
