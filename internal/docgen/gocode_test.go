package docgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractFunction(t *testing.T) {
	// Create a temp Go file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")

	content := `package test

import "fmt"

func HelloWorld() {
    fmt.Println("hello")
}
`
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)

	// Extract function
	result, err := ExtractGoCode(tmpFile, "HelloWorld", nil)
	require.NoError(t, err)

	assert.Contains(t, result, "func HelloWorld()")
	assert.Contains(t, result, `fmt.Println("hello")`)
}

func TestExtractFunctionStripImports(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")

	content := `package test

import "fmt"

func HelloWorld() {
    fmt.Println("hello")
}
`
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)

	result, err := ExtractGoCode(tmpFile, "HelloWorld", []string{"stripImports"})
	require.NoError(t, err)

	assert.Contains(t, result, "func HelloWorld()")
	assert.NotContains(t, result, `import "fmt"`)
}

func TestExtractStruct(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")

	content := `package test

type Person struct {
    Name string
    Age  int
}
`
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)

	result, err := ExtractGoCode(tmpFile, "Person", nil)
	require.NoError(t, err)

	assert.Contains(t, result, "type Person struct")
	assert.Contains(t, result, "Name string")
	assert.Contains(t, result, "Age  int")
}

func TestExtractFunctionWithRename(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")

	content := `package test

func OldName() string {
    return "test"
}
`
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)

	result, err := ExtractGoCode(tmpFile, "OldName", []string{"rename:OldName→NewName"})
	require.NoError(t, err)

	assert.Contains(t, result, "func NewName()")
	assert.NotContains(t, result, "OldName")
}

func TestExtractFunctionWithPackageRename(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")

	content := `package test

func MyFunc() string {
    return "test"
}
`
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)

	result, err := ExtractGoCode(tmpFile, "MyFunc", []string{"package:example"})
	require.NoError(t, err)

	assert.Contains(t, result, "package example")
	assert.NotContains(t, result, "package test")
}

func TestExtractFunctionNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")

	content := `package test

func Other() {}
`
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)

	_, err = ExtractGoCode(tmpFile, "NotFound", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}