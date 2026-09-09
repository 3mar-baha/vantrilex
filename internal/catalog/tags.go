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
	}
	return false
}

func hasVisionSignal(hay string) bool {
	hay = strings.ToLower(hay)
	for _, k := range []string{"vision", "multimodal", "image", "maverick", "gemini", "glm", "minimax", "llama-4", "opus-4-6-vision"} {
		if strings.Contains(hay, k) {
			return true
		}
	}
	return false
}

func hasCodingSignal(hay string) bool {
	hay = strings.ToLower(hay)
	for _, k := range []string{"coder", "codex", "coding", "sonnet", "haiku", "kimi", "qwen", "flash", "code"} {
		if strings.Contains(hay, k) {
			return true
		}
	}
	return false
}

func priceValue(s string) (float64, bool) {
	s = strings.TrimSpace(strings.TrimPrefix(s, "$"))
	if s == "" || s == "varies" || s == "n/a" {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

func isFree(m Model) bool {
	if in, ok := priceValue(m.InputPerM); ok && in == 0 {
		if out, ok := priceValue(m.OutputPerM); ok && out == 0 {
			return true
		}
	}
	return false
}

func isFast(m Model, hay string) bool {
	if strings.ToUpper(m.Latency) == "LOW" {
		return true
	}
	for _, k := range []string{"flash", "haiku", "fast", "instant", "turbo", "kimi", "qwen"} {
		if strings.Contains(hay, k) {
			return true
		}
	}
	return false
}
