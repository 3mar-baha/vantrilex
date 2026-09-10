package particles

import "testing"

func BenchmarkSupernovaTick(b *testing.B) {
	e := New(100, 24)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Supernova(130)
		e.Tick(1.0 / 60.0)
	}
}

func BenchmarkRenderField(b *testing.B) {
	e := New(100, 24)
	e.Supernova(130)
	e.Tick(1.0 / 60.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.RenderField()
	}
}
