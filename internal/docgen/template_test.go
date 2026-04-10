package docgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFuncMapHasAllFunctions(t *testing.T) {
	fm := FuncMap()

	// Core functions
	assert.Contains(t, fm, "config")
	assert.Contains(t, fm, "launch_game")
	assert.Contains(t, fm, "end_game")
	assert.Contains(t, fm, "command")
	assert.Contains(t, fm, "delay")
	assert.Contains(t, fm, "verifyOutputs")
	assert.Contains(t, fm, "gocode")

	// Utility functions
	assert.Contains(t, fm, "list")
	assert.Contains(t, fm, "dict")
}

func TestProcessSimpleTemplate(t *testing.T) {
	tmplContent := `Hello {{ .Name }}!`

	result, err := ProcessTemplateString(tmplContent, map[string]string{"Name": "World"})
	require.NoError(t, err)
	assert.Equal(t, "Hello World!", result)
}

func TestProcessTemplateWithVariables(t *testing.T) {
	tmplContent := `Name: {{ .Name }}
Age: {{ .Age }}`

	data := map[string]any{
		"Name": "Test",
		"Age":  42,
	}

	result, err := ProcessTemplateString(tmplContent, data)
	require.NoError(t, err)
	assert.Equal(t, "Name: Test\nAge: 42", result)
}

func TestUtilityFunctions(t *testing.T) {
	fm := FuncMap()

	// Test list
	listFunc := fm["list"].(func(...any) []any)
	list := listFunc("a", "b", "c")
	assert.Equal(t, []any{"a", "b", "c"}, list)

	// Test dict
	dictFunc := fm["dict"].(func(...any) map[string]any)
	d := dictFunc("key1", "val1", "key2", "val2")
	assert.Equal(t, map[string]any{"key1": "val1", "key2": "val2"}, d)
}

func TestListFuncWithMixedTypes(t *testing.T) {
	fm := FuncMap()
	listFunc := fm["list"].(func(...any) []any)

	list := listFunc("string", 42, true, 3.14)
	assert.Equal(t, []any{"string", 42, true, 3.14}, list)
}

func TestDictFuncWithNonStringKey(t *testing.T) {
	fm := FuncMap()
	dictFunc := fm["dict"].(func(...any) map[string]any)

	// Non-string keys should be skipped
	d := dictFunc(123, "val", "valid", "value")
	assert.Equal(t, map[string]any{"valid": "value"}, d)
}

func TestDictFuncWithOddArgs(t *testing.T) {
	fm := FuncMap()
	dictFunc := fm["dict"].(func(...any) map[string]any)

	// Odd number of args - last value should be ignored
	d := dictFunc("key1", "val1", "key2")
	assert.Equal(t, map[string]any{"key1": "val1"}, d)
}

func TestProcessTemplateStringResetsContext(t *testing.T) {
	// Create a new global context
	globalContext = NewContext()

	// Process a template (this should reset context)
	_, err := ProcessTemplateString("test", nil)
	require.NoError(t, err)

	// Context should be fresh
	assert.Nil(t, globalContext.Config)
	assert.Nil(t, globalContext.GameSession)
}

func TestDelayFunc(t *testing.T) {
	fm := FuncMap()
	delayFunc := fm["delay"].(func(string) (string, error))

	// Test with valid duration (very short to keep test fast)
	result, err := delayFunc("1ms")
	require.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestVerifyOutputsFunc(t *testing.T) {
	fm := FuncMap()
	verifyOutputsFunc := fm["verifyOutputs"].(func(...string) (string, error))

	// Test with matching outputs
	result, err := verifyOutputsFunc("same", "same", "same")
	require.NoError(t, err)
	assert.Equal(t, "", result)

	// Test with no outputs
	result, err = verifyOutputsFunc()
	require.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestVerifyOutputsFuncMismatch(t *testing.T) {
	fm := FuncMap()
	verifyOutputsFunc := fm["verifyOutputs"].(func(...string) (string, error))

	// Test with mismatched outputs
	_, err := verifyOutputsFunc("same", "different")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "output mismatch")
}

func TestConfigFunc(t *testing.T) {
	// Reset context
	globalContext = NewContext()

	fm := FuncMap()
	configFunc := fm["config"].(func(string, ...map[string]any) (string, error))

	result, err := configFunc("/tmp/game")
	require.NoError(t, err)
	assert.Equal(t, "", result)
	assert.NotNil(t, globalContext.Config)
	assert.Equal(t, "/tmp/game", globalContext.Config.GameDir)
}

func TestEndGameFuncWithNilSession(t *testing.T) {
	// Reset context with no game session
	globalContext = NewContext()

	fm := FuncMap()
	endGameFunc := fm["end_game"].(func() (string, error))

	// Should not error with nil session
	result, err := endGameFunc()
	require.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestCommandFuncWithNoSession(t *testing.T) {
	// Reset context with no game session
	globalContext = NewContext()

	fm := FuncMap()
	commandFunc := fm["command"].(func(string, ...map[string]any) (string, error))

	// Should error with no active session
	_, err := commandFunc("screenshot")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active game session")
}

func TestProcessTemplateFile(t *testing.T) {
	// Create temp template file
	tmpDir := t.TempDir()
	tmplPath := filepath.Join(tmpDir, "test.md.tmpl")

	templateContent := "# Test\nHello {{ .Name }}!"
	err := os.WriteFile(tmplPath, []byte(templateContent), 0644)
	require.NoError(t, err)

	// Process template from file
	result, err := ProcessTemplate(tmplPath)
	require.NoError(t, err)
	// Go templates output "<no value>" for missing fields when data is nil
	assert.Contains(t, result, "# Test")
	assert.Contains(t, result, "Hello")

	// Context should be reset
	assert.Nil(t, globalContext.Config)
}

func TestProcessTemplateFileNotFound(t *testing.T) {
	_, err := ProcessTemplate("/nonexistent/path/template.md.tmpl")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read template")
}

func TestProcessTemplateStringParseError(t *testing.T) {
	// Invalid template syntax
	tmplContent := `Hello {{ .Name`

	_, err := ProcessTemplateString(tmplContent, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse template")
}

func TestConfigFuncWithNormalizeRules(t *testing.T) {
	// Reset context
	globalContext = NewContext()

	fm := FuncMap()
	configFunc := fm["config"].(func(string, ...map[string]any) (string, error))

	// Call with normalize rules
	rules := map[string]any{
		"rules": []NormalizeRule{
			{Pattern: "\\d+", Replace: "N"},
		},
	}
	result, err := configFunc("/tmp/game", rules)
	require.NoError(t, err)
	assert.Equal(t, "", result)
	assert.NotNil(t, globalContext.Config)
	assert.Len(t, globalContext.Config.Normalize, 1)
	assert.Equal(t, "\\d+", globalContext.Config.Normalize[0].Pattern)
}

func TestConfigFuncWithEmptyRules(t *testing.T) {
	// Reset context
	globalContext = NewContext()

	fm := FuncMap()
	configFunc := fm["config"].(func(string, ...map[string]any) (string, error))

	// Call with empty rules map
	result, err := configFunc("/tmp/game", map[string]any{})
	require.NoError(t, err)
	assert.Equal(t, "", result)
	assert.NotNil(t, globalContext.Config)
	assert.Len(t, globalContext.Config.Normalize, 0)
}

func TestGocodeFunc(t *testing.T) {
	// Create a test Go file
	tmpDir := t.TempDir()
	goFilePath := filepath.Join(tmpDir, "test.go")

	goContent := `package test

func Example() int {
	return 42
}

type MyStruct struct {
	Name string
}
`
	err := os.WriteFile(goFilePath, []byte(goContent), 0644)
	require.NoError(t, err)

	fm := FuncMap()
	gocodeFunc := fm["gocode"].(func(string, string, ...any) (string, error))

	// Extract function
	result, err := gocodeFunc(goFilePath, "Example")
	require.NoError(t, err)
	assert.Contains(t, result, "func Example()")

	// Extract struct
	result, err = gocodeFunc(goFilePath, "MyStruct")
	require.NoError(t, err)
	assert.Contains(t, result, "type MyStruct struct")
}

func TestGocodeFuncWithTransforms(t *testing.T) {
	// Create a test Go file
	tmpDir := t.TempDir()
	goFilePath := filepath.Join(tmpDir, "test.go")

	goContent := `package test

func OldName() int {
	return 42
}
`
	err := os.WriteFile(goFilePath, []byte(goContent), 0644)
	require.NoError(t, err)

	fm := FuncMap()
	gocodeFunc := fm["gocode"].(func(string, string, ...any) (string, error))

	// Extract with rename transform
	result, err := gocodeFunc(goFilePath, "OldName", []any{"rename:OldName->NewName"})
	require.NoError(t, err)
	assert.Contains(t, result, "NewName")
	assert.NotContains(t, result, "OldName")
}

func TestGocodeFuncNotFound(t *testing.T) {
	// Create a test Go file
	tmpDir := t.TempDir()
	goFilePath := filepath.Join(tmpDir, "test.go")

	goContent := `package test`
	err := os.WriteFile(goFilePath, []byte(goContent), 0644)
	require.NoError(t, err)

	fm := FuncMap()
	gocodeFunc := fm["gocode"].(func(string, string, ...any) (string, error))

	// Extract nonexistent function
	_, err = gocodeFunc(goFilePath, "NonExistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestCommandFuncWithFlags(t *testing.T) {
	// Reset context with no game session
	globalContext = NewContext()

	fm := FuncMap()
	commandFunc := fm["command"].(func(string, ...map[string]any) (string, error))

	// Should still error with no active session, even with flags
	flags := map[string]any{"format": "png"}
	_, err := commandFunc("screenshot", flags)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active game session")
}

func TestConfigFuncWithMapBasedNormalizeRules(t *testing.T) {
	// Reset context
	globalContext = NewContext()

	fm := FuncMap()
	configFunc := fm["config"].(func(string, ...map[string]any) (string, error))

	// Call with map-based normalize rules (from template dict/list)
	rules := map[string]any{
		"rules": []map[string]any{
			{"Pattern": "\\d+", "Replace": "N"},
			{"Pattern": "PID=\\d+", "Replace": "PID=<PID>"},
		},
	}
	result, err := configFunc("/tmp/game", rules)
	require.NoError(t, err)
	assert.Equal(t, "", result)
	assert.NotNil(t, globalContext.Config)
	assert.Len(t, globalContext.Config.Normalize, 2)
	assert.Equal(t, "\\d+", globalContext.Config.Normalize[0].Pattern)
	assert.Equal(t, "N", globalContext.Config.Normalize[0].Replace)
	assert.Equal(t, "PID=\\d+", globalContext.Config.Normalize[1].Pattern)
	assert.Equal(t, "PID=<PID>", globalContext.Config.Normalize[1].Replace)
}

func TestConfigFuncWithSliceAnyNormalizeRules(t *testing.T) {
	// Reset context
	globalContext = NewContext()

	fm := FuncMap()
	configFunc := fm["config"].(func(string, ...map[string]any) (string, error))

	// Call with []any normalize rules (from template list function)
	rules := map[string]any{
		"rules": []any{
			map[string]any{"Pattern": "_addr=\"0x[0-9a-f]+\"", "Replace": "_addr=\"<ADDR>\""},
			map[string]any{"Pattern": "PID=\\d+", "Replace": "PID=<PID>"},
		},
	}
	result, err := configFunc("/tmp/game", rules)
	require.NoError(t, err)
	assert.Equal(t, "", result)
	assert.NotNil(t, globalContext.Config)
	assert.Len(t, globalContext.Config.Normalize, 2)
	assert.Equal(t, "_addr=\"0x[0-9a-f]+\"", globalContext.Config.Normalize[0].Pattern)
	assert.Equal(t, "_addr=\"<ADDR>\"", globalContext.Config.Normalize[0].Replace)
}

func TestEndGameFuncClearsSession(t *testing.T) {
	// Reset context
	globalContext = NewContext()

	// Manually set a mock session to test the clearing logic
	globalContext.GameSession = &GameSession{}

	fm := FuncMap()
	endGameFunc := fm["end_game"].(func() (string, error))

	// End game - will fail since session is not properly initialized
	_, err := endGameFunc()
	// Session should be cleared even if EndGame fails
	assert.Nil(t, globalContext.GameSession)
	// The error from EndGame with nil game is nil
	require.NoError(t, err)
}