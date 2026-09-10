package ui

import (
	"testing"

	"vantrilex/internal/catalog"
)

func mkItems() []catalog.RegistryItem {
	return []catalog.RegistryItem{
		{Name: "alpha", Desc: "first tool", Selected: true},
		{Name: "beta", Desc: "second tool"},
		{Name: "gamma-filesystem", Desc: "files"},
	}
}

func TestVListToggleAndSearch(t *testing.T) {
	v := NewVList(mkItems())
	if v.SelectedCount() != 1 {
		t.Fatalf("want 1 preselected, got %d", len(v.Selected))
	}
	v.Move(1)
	v.Toggle()
	if v.SelectedCount() != 2 {
		t.Fatalf("toggle failed: %v", v.Selected)
	}
	v.SetSearch("filesystem")
	if len(v.Filtered) != 1 {
		t.Fatalf("search got %d", len(v.Filtered))
	}
	v.ToggleAllFiltered()
	if v.SelectedCount() != 3 { // gamma selected, others untouched
		t.Fatalf("toggle-all-filtered got %d", v.SelectedCount())
	}
	v.ToggleAllFiltered()
	if v.SelectedCount() != 2 { // gamma cleared again
		t.Fatalf("toggle-all-clear got %d", v.SelectedCount())
	}
}

func TestVListWindow(t *testing.T) {
	items := []catalog.RegistryItem{}
	for i := 0; i < 50; i++ {
		items = append(items, catalog.RegistryItem{Name: string(rune('a'+i%26)) + string(rune('0'+i%10)) + "-", Desc: "x"})
	}
	v := NewVList(items)
	v.PageSize = 12
	for i := 0; i < 20; i++ {
		v.Move(1)
	}
	if len(v.Visible()) != 12 {
		t.Fatalf("window=%d", len(v.Visible()))
	}
}
