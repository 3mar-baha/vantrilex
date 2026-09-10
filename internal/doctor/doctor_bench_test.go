package doctor

import "testing"

func BenchmarkPill(b *testing.B) {
	s := RunnerUpdateStatus{Updating: false, Done: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Pill()
	}
}

func BenchmarkRunnerUpdates(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RunnerUpdates()
	}
}
