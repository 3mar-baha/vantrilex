// Automated capability tagging engine for model badges.
package catalog

import (
	"strconv"
	"strings"
)

// ModelTags evaluates model metadata and returns visual badges:
// [REASONING] [VISION] [CODING] [FREE TIER] [ULTRA FAST].
func ModelTags(m Model) []string {
	hay := strings.ToLower(m.ID + " " + m.Short + " " + m.Blurb)
	var tags []string
	if m.Reasoning || hasReasoningSignal(hay) {
		tags = append(tags, "[REASONING]")
	}
	if hasVisionSignal(hay) {
		tags = append(tags, "[VISION]")