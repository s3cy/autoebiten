package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/s3cy/autoebiten/internal/docgen"
	"github.com/s3cy/autoebiten/testkit"
)

func main() {
	if len(os.Args) < 2 {
		fatal(fmt.Errorf("usage: docverify <example_dir>"))
	}
	verifyExamples(os.Args[1])
}

func verifyExamples(exampleDir string) {
	// Load config
	configPath := filepath.Join(exampleDir, "config.yaml")
	config, err := docgen.LoadConfig(configPath)
	if err != nil {
		fatal(err)
	}

	// Build game binary
	gameBin := filepath.Join(config.GameDir, "autoui_demo")
	fmt.Printf("Building: %s\n", config.GameDir)
	buildCmd := exec.Command("go", "build", "-o", "autoui_demo", ".")
	buildCmd.Dir = config.GameDir
	if err := buildCmd.Run(); err != nil {
		fatal(fmt.Errorf("failed to build game: %w", err))
	}

	// Launch game using testkit
	fmt.Println("Launching game...")
	t := &testing.T{}
	game := testkit.Launch(t, gameBin, testkit.WithTimeout(30*time.Second))

	// Wait for game to be ready
	ready := game.WaitFor(func() bool {
		return game.Ping() == nil
	}, 5*time.Second)
	if !ready {
		game.Shutdown()
		fatal(fmt.Errorf("game failed to start"))
	}
	fmt.Println("Game ready")

	// Verify each script
	hasMismatch := false
	scripts := findScripts(exampleDir)

	for _, script := range scripts {
		name := strings.TrimSuffix(filepath.Base(script), ".sh")
		outputFile := filepath.Join(exampleDir, name+"_out.txt")

		fmt.Printf("Verifying: %s\n", name)

		// Run script
		liveOutput := runScript(script)

		// Normalize live output
		liveNorm := docgen.Normalize(liveOutput, config.Normalize)

		// Load and normalize expected
		expected, err := os.ReadFile(outputFile)
		if err != nil {
			fmt.Printf("  MISSING: %s\n", outputFile)
			hasMismatch = true
			continue
		}
		expectedNorm := docgen.Normalize(string(expected), config.Normalize)

		// Compare
		if liveNorm == expectedNorm {
			fmt.Println("  OK: matches")
		} else {
			fmt.Println("  MISMATCH: output differs")
			printDiff(expectedNorm, liveNorm)
			hasMismatch = true
		}
	}

	// Cleanup
	game.Shutdown()
	fmt.Println("Game stopped")

	if hasMismatch {
		fmt.Println("\nFAILED: Some outputs do not match")
		fmt.Println("Run: make docs-generate to update outputs")
		os.Exit(1)
	}
	fmt.Println("\nSUCCESS: All outputs match")
}

func findScripts(dir string) []string {
	files, err := filepath.Glob(filepath.Join(dir, "*.sh"))
	if err != nil {
		fatal(err)
	}
	return files
}

func runScript(script string) string {
	cmd := exec.Command("bash", script)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(exitErr.Stderr)
		}
		return fmt.Sprintf("ERROR: %v", err)
	}
	return string(output)
}

func printDiff(expected, live string) {
	fmt.Println("=== Full diff (normalized) ===")

	expectedLines := strings.Split(expected, "\n")
	liveLines := strings.Split(live, "\n")

	maxLen := max(len(expectedLines), len(liveLines))
	for i := 0; i < maxLen; i++ {
		if i < len(expectedLines) && i < len(liveLines) {
			if expectedLines[i] != liveLines[i] {
				fmt.Printf("  @@ line %d @@\n", i+1)
				fmt.Printf("  -%s\n", expectedLines[i])
				fmt.Printf("  +%s\n", liveLines[i])
			}
		} else if i < len(expectedLines) {
			fmt.Printf("  -%s\n", expectedLines[i])
		} else {
			fmt.Printf("  +%s\n", liveLines[i])
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
