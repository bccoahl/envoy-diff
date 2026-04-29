package dotenv

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ValidationFormatter formats ValidationIssue slices for output.
type ValidationFormatter struct {
	Format string // "text", "json"
}

// Render converts issues to a formatted string.
func (vf ValidationFormatter) Render(issues []ValidationIssue) (string, error) {
	switch strings.ToLower(vf.Format) {
	case "json":
		return vf.renderJSON(issues)
	default:
		return vf.renderText(issues), nil
	}
}

func (vf ValidationFormatter) renderText(issues []ValidationIssue) string {
	if len(issues) == 0 {
		return "✔ no validation issues"
	}
	var sb strings.Builder
	for _, i := range issues {
		symbol := "✖"
		if i.Level == ValidationWarn {
			symbol = "⚠"
		}
		fmt.Fprintf(&sb, "%s [%s] %s — %s\n",
			symbol,
			strings.ToUpper(string(i.Level)),
			i.Key,
			i.Message,
		)
	}
	return strings.TrimRight(sb.String(), "\n")
}

type jsonIssue struct {
	Key     string `json:"key"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

func (vf ValidationFormatter) renderJSON(issues []ValidationIssue) (string, error) {
	payload := make([]jsonIssue, len(issues))
	for idx, i := range issues {
		payload[idx] = jsonIssue{
			Key:     i.Key,
			Level:   string(i.Level),
			Message: i.Message,
		}
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", fmt.Errorf("validation json render: %w", err)
	}
	return string(b), nil
}
