package docgen

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// globalContext is the shared context for template execution
var globalContext = NewContext()

// FuncMap returns the template function map.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"config":        configFunc,
		"launch_game":   launchGameFunc,
		"end_game":      endGameFunc,
		"command":       commandFunc,
		"delay":         delayFunc,
		"verifyOutputs": verifyOutputsFunc,
		"gocode":        gocodeFunc,
		"list":          listFunc,
		"dict":          dictFunc,
	}
}

// ProcessTemplate processes a template file.
func ProcessTemplate(tmplPath string) (string, error) {
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	// Reset global context for each template execution
	globalContext = NewContext()

	return ProcessTemplateString(string(content), nil)
}

// ProcessTemplateString processes a template string with the given data.
func ProcessTemplateString(tmplContent string, data any) (string, error) {
	tmpl, err := template.New("doc").Funcs(FuncMap()).Parse(tmplContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// configFunc sets game directory and normalization rules.
func configFunc(gameDir string, normalizeRules ...map[string]any) (string, error) {
	cfg := &Config{GameDir: gameDir}

	// Parse normalize rules if provided
	if len(normalizeRules) > 0 && normalizeRules[0] != nil {
		rules := normalizeRules[0]
		// Handle []NormalizeRule (from Go code)
		if rulesList, ok := rules["rules"].([]NormalizeRule); ok {
			cfg.Normalize = rulesList
		}
		// Handle []map[string]any (from template dict/list)
		if rulesList, ok := rules["rules"].([]map[string]any); ok {
			for _, r := range rulesList {
				pattern, _ := r["Pattern"].(string)
				replace, _ := r["Replace"].(string)
				if pattern != "" {
					cfg.Normalize = append(cfg.Normalize, NormalizeRule{
						Pattern: pattern,
						Replace: replace,
					})
				}
			}
		}
	}

	globalContext.SetConfig(cfg)
	return "", nil
}

// launchGameFunc starts a game session.
func launchGameFunc(args ...string) (string, error) {
	session, err := LaunchGame(globalContext, args...)
	if err != nil {
		return "", err
	}
	globalContext.GameSession = session
	return "", nil
}

// endGameFunc ends the current game session.
func endGameFunc() (string, error) {
	if globalContext.GameSession == nil {
		return "", nil
	}
	err := EndGame(globalContext.GameSession)
	globalContext.GameSession = nil
	return "", err
}

// commandFunc executes an autoebiten command.
func commandFunc(cmdName string, flags ...map[string]any) (string, error) {
	if globalContext.GameSession == nil {
		return "", fmt.Errorf("no active game session")
	}
	var flagMap map[string]any
	if len(flags) > 0 {
		flagMap = flags[0]
	}
	return ExecuteCommand(globalContext.GameSession, cmdName, flagMap)
}

// delayFunc waits for the specified duration.
func delayFunc(duration string) (string, error) {
	Delay(duration)
	return "", nil
}

// verifyOutputsFunc verifies outputs match.
func verifyOutputsFunc(outputs ...string) (string, error) {
	return "", VerifyOutputs(outputs...)
}

// gocodeFunc extracts Go code from a file.
func gocodeFunc(filePath, target string, transforms ...[]string) (string, error) {
	var t []string
	if len(transforms) > 0 {
		t = transforms[0]
	}
	return ExtractGoCode(filePath, target, t)
}

// listFunc creates a slice from the given items.
func listFunc(items ...any) []any {
	return items
}

// dictFunc creates a map from alternating key-value pairs.
func dictFunc(items ...any) map[string]any {
	d := make(map[string]any)
	for i := 0; i+1 < len(items); i += 2 {
		key, ok := items[i].(string)
		if ok {
			d[key] = items[i+1]
		}
	}
	return d
}