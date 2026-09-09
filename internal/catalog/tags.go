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
	}
	if hasCodingSignal(hay) {
		tags = append(tags, "[CODING]")
	}
	if isFree(m) {
		tags = append(tags, "[FREE TIER]")
	}
	if isFast(m, hay) {
		tags = append(tags, "[ULTRA FAST]")
	}
	return tags
}

func hasReasoningSignal(hay string) bool {
	hay = strings.ToLower(hay)
	for _, k := range []string{"reasoning", "thinking", "think", "r1", "o1", "o3", "o4", "deepseek", "opus", "grok", "qwen3", "kimi-k2", "glm", "minimax"} {
		if strings.Contains(hay, k) {
			return true
		}