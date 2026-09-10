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

func BenchmarkEmitTrail(b *testing.B) {
	e := New(100, 24)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.EmitTrail(50, 12)
	}
}

func BenchmarkPerimeterSupernova(b *testing.B) {
	e := New(100, 24)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.PerimeterSupernova()
	}
}
