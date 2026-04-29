package dotenv

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationLevel controls how strict validation is.
type ValidationLevel string

const (
	ValidationWarn  ValidationLevel = "warn"
	ValidationError ValidationLevel = "error"
)

// ValidationIssue represents a single validation problem found in a diff entry.
type ValidationIssue struct {
	Key     string
	Message string
	Level   ValidationLevel
}

// ValidateOptions configures the validation behaviour.
type ValidateOptions struct {
	// RequireKeys fails if any of these keys are missing from both sides.
	RequireKeys []string
	// ForbidPattern marks keys matching this regex as issues.
	ForbidPattern string
	// WarnOnChanged emits a warning for every changed key.
	WarnOnChanged bool
}

// Validate inspects a slice of DiffEntry values and returns any issues found.
func Validate(entries []DiffEntry, opts ValidateOptions) []ValidationIssue {
	var issues []ValidationIssue

	presentKeys := make(map[string]bool, len(entries))
	for _, e := range entries {
		presentKeys[e.Key] = true
	}

	for _, req := range opts.RequireKeys {
		if !presentKeys[req] {
			issues = append(issues, ValidationIssue{
				Key:     req,
				Message: fmt.Sprintf("required key %q is missing from both sources", req),
				Level:   ValidationError,
			})
		}
	}

	var forbidRe *regexp.Regexp
	if opts.ForbidPattern != "" {
		forbidRe = regexp.MustCompile(opts.ForbidPattern)
	}

	for _, e := range entries {
		if forbidRe != nil && forbidRe.MatchString(e.Key) {
			issues = append(issues, ValidationIssue{
				Key:     e.Key,
				Message: fmt.Sprintf("key %q matches forbidden pattern %q", e.Key, opts.ForbidPattern),
				Level:   ValidationError,
			})
		}
		if opts.WarnOnChanged && e.Status == StatusChanged {
			issues = append(issues, ValidationIssue{
				Key:     e.Key,
				Message: fmt.Sprintf("key %q has changed", e.Key),
				Level:   ValidationWarn,
			})
		}
	}

	return issues
}

// HasErrors returns true if any issue has ErrorLevel.
func HasErrors(issues []ValidationIssue) bool {
	for _, i := range issues {
		if i.Level == ValidationError {
			return true
		}
	}
	return false
}

// FormatIssues returns a human-readable summary of all issues.
func FormatIssues(issues []ValidationIssue) string {
	if len(issues) == 0 {
		return "no validation issues found"
	}
	var sb strings.Builder
	for _, i := range issues {
		fmt.Fprintf(&sb, "[%s] %s: %s\n", strings.ToUpper(string(i.Level)), i.Key, i.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}
