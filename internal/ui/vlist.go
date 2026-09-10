package ui

import (
	"strings"

	"vantrilex/internal/catalog"
)

// VList is a virtualized multi-select list: only ~12 rows render, filtering
// is token-AND fuzzy over name+desc+tags, cursor stays in the filtered view.
type VList struct {
	Items    []catalog.RegistryItem
	Filtered []int // indices into Items
	Cursor   int
	Offset   int
	Search   string
	Selected map[string]bool
	PageSize int
}

func NewVList(items []catalog.RegistryItem) VList {
	sel := map[string]bool{}
	for _, it := range items {
		if it.Selected {
			sel[it.Name] = true
		}
	}
	v := VList{Items: items, Selected: sel, PageSize: 12}
	v.refilter()
	return v
}

func (v *VList) refilter() {
	v.Filtered = v.Filtered[:0]
	q := strings.ToLower(strings.TrimSpace(v.Search))
	toks := strings.Fields(q)
	for i, it := range v.Items {
		if q == "" {
			v.Filtered = append(v.Filtered, i)
			continue
		}
		hay := strings.ToLower(it.Name + " " + it.Desc + " " + strings.Join(it.Tags, " ") + " " + it.Source)
		ok := true
		for _, t := range toks {
			if !strings.Contains(hay, t) {
				ok = false
				break
			}
		}
		if ok {
			v.Filtered = append(v.Filtered, i)
		}
	}
	if v.Cursor >= len(v.Filtered) {
		v.Cursor = 0
	}
	if v.Cursor < 0 {
		v.Cursor = 0
	}
	v.clampOffset()
}

func (v *VList) clampOffset() {
	if v.PageSize <= 0 {
		v.PageSize = 12
	}
	if v.Cursor < v.Offset {
		v.Offset = v.Cursor
	}
	if v.Cursor >= v.Offset+v.PageSize {
		v.Offset = v.Cursor - v.PageSize + 1
	}
	if v.Offset < 0 {
		v.Offset = 0
	}
}

// SetSearch updates the fuzzy query and resets cursor.
func (v *VList) SetSearch(s string) {
	v.Search = s
	v.Cursor = 0
	v.Offset = 0
	v.refilter()
}

// Move steps the cursor with wraparound.
func (v *VList) Move(delta int) {
	n := len(v.Filtered)
	if n == 0 {
		return
	}
	v.Cursor = (v.Cursor + delta + n) % n
	v.clampOffset()
}

// Toggle flips the cursor item.
func (v *VList) Toggle() {
	if len(v.Filtered) == 0 {
		return
	}
	name := v.Items[v.Filtered[v.Cursor]].Name
	if v.Selected[name] {
		delete(v.Selected, name)
	} else {
		v.Selected[name] = true
	}
}

// ToggleAllFiltered selects all filtered if any unselected, else clears them.
func (v *VList) ToggleAllFiltered() {
	anyOff := false
	for _, idx := range v.Filtered {
		if !v.Selected[v.Items[idx].Name] {
			anyOff = true
			break
		}
	}
	for _, idx := range v.Filtered {
		name := v.Items[idx].Name
		if anyOff {
			v.Selected[name] = true
		} else {
			delete(v.Selected, name)
		}
	}
}

// SelectedItems returns checked items in catalog order.
func (v *VList) SelectedItems() []catalog.RegistryItem {
	out := []catalog.RegistryItem{}
	for _, it := range v.Items {
		if v.Selected[it.Name] {
			out = append(out, it)
		}
	}
	return out
}

// SelectedCount returns checked count.
func (v *VList) SelectedCount() int { return len(v.Selected) }

// Visible returns the current window of filtered indices.
func (v *VList) Visible() []int {
	if v.PageSize <= 0 {
		v.PageSize = 12
	}
	end := v.Offset + v.PageSize
	if end > len(v.Filtered) {
		end = len(v.Filtered)
	}
	if v.Offset > end {
		v.Offset = 0
	}
	return v.Filtered[v.Offset:end]
}
