package ui

import (
	"testing"

	"vantrilex/internal/catalog"
)

func BenchmarkVListFilter(b *testing.B) {
	v := NewVList(catalog.MCPs())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.SetSearch("sql")
		v.SetSearch("")
	}
}

func BenchmarkVListWindow(b *testing.B) {
	v := NewVList(catalog.MCPs())
	v.PageSize = 12
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Move(1)
		_ = v.Visible()
	}
}
