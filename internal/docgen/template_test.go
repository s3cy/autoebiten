package docgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutputFunc(t *testing.T) {
	// Create temp output file
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test_out.txt")

	content := `<UI>
  <Button id="btn"/>
</UI>`
	err := os.WriteFile(outputPath, []byte(content), 0644)
	require.NoError(t, err)

	// Create output function with base dir
	outputFn := OutputFunc(tmpDir)

	// Test reading file
	result, err := outputFn("test_out.txt", "xml")
	require.NoError(t, err)

	expected := "```xml\n<UI>\n  <Button id=\"btn\"/>\n</UI>\n```"
	assert.Equal(t, expected, result)
}

func TestOutputFuncNotFound(t *testing.T) {
	outputFn := OutputFunc("/nonexistent")
	_, err := outputFn("missing.txt", "text")
	assert.Error(t, err)
}
