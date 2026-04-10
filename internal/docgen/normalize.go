package docgen

import (
	"regexp"
	"strings"
)

// Normalize applies all normalization rules to the input string.
// It processes rules in order and trims leading/trailing whitespace.
func Normalize(s string, rules []NormalizeRule) string {
	for _, r := range rules {
		re := regexp.MustCompile(r.Pattern)
		s = re.ReplaceAllString(s, r.Replace)
	}
	return strings.TrimSpace(s)
}