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
		processAllTemplates()
		return
	}

	switch os.Args[1] {
	case "--generate":
		if len(os.Args) < 3 {
			fatal(fmt.Errorf("usage: docgen --generate <example_dir>"))
		}
		generateExamples(os.Args[2])
	case "--process":
		processAllTemplates()
	default:
		fatal(fmt.Errorf("unknown command: %s", os.Args[1]))
	}
}

func generateExamples(exampleDir string) {
	// Load config
	configPath := filepath.Join(exampleDir, "config.yaml")
	config, err := docgen.LoadConfig(configPath)
	if err != nil {
		fatal(err)
	}

	// Build game binary - use directory name + "_demo" as binary name
	gameBinary := filepath.Base(config.GameDir) + "_demo"
	gameBin := filepath.Join(config.GameDir, gameBinary)
	fmt.Printf("Building: %s (binary: %s)\n", config.GameDir, gameBinary)
	buildCmd := exec.Command("go", "build", "-o", gameBinary, ".")
	buildCmd.Dir = config.GameDir
	if output, err := buildCmd.CombinedOutput(); err != nil {
		fatal(fmt.Errorf("failed to build game: %w\n%s", err, output))
	}

	// Launch game using testkit
	fmt.Printf("Launching game: %s\n", gameBin)
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

	// Find and run all scripts
	scripts := findScripts(exampleDir)
	for _, script := range scripts {
		name := strings.TrimSuffix(filepath.Base(script), ".sh")
		outputFile := filepath.Join(exampleDir, name+"_out.txt")

		fmt.Printf("Running: %s\n", script)
		output := runScript(script)

		// Normalize output
		normalized := docgen.Normalize(output, config.Normalize)

		// Save output
		if err := os.WriteFile(outputFile, []byte(normalized), 0644); err != nil {
			game.Shutdown()
			fatal(fmt.Errorf("failed to write output: %w", err))
		}
		fmt.Printf("Generated: %s\n", outputFile)
	}

	// Shutdown game
	game.Shutdown()
	fmt.Println("Game stopped")
}

func processAllTemplates() {
	generateDir := "docs/generate"

	// Find all .tmpl files recursively
	templates, err := filepath.Glob(filepath.Join(generateDir, "**/*.tmpl"))
	if err != nil {
		fatal(err)
	}

	// Also check for .tmpl files in the root of generateDir
	topLevelTemplates, err := filepath.Glob(filepath.Join(generateDir, "*.tmpl"))
	if err != nil {
		fatal(err)
	}
	templates = append(templates, topLevelTemplates...)

	if len(templates) == 0 {
		fmt.Println("No templates found")
		return
	}

	for _, tmplPath := range templates {
		// Output file: strip .tmpl suffix, place in docs/
		relPath, err := filepath.Rel(generateDir, tmplPath)
		if err != nil {
			fatal(err)
		}
		outName := strings.TrimSuffix(relPath, ".tmpl")
		outPath := filepath.Join("docs", outName)

		fmt.Printf("Processing: %s\n", tmplPath)

		result, err := docgen.ProcessTemplate(tmplPath, generateDir)
		if err != nil {
			fatal(err)
		}

		// Ensure output directory exists
		outDir := filepath.Dir(outPath)
		if err := os.MkdirAll(outDir, 0755); err != nil {
			fatal(fmt.Errorf("failed to create output directory: %w", err))
		}

		if err := os.WriteFile(outPath, []byte(result), 0644); err != nil {
			fatal(fmt.Errorf("failed to write %s: %w", outPath, err))
		}
		fmt.Printf("Generated: %s -> %s\n", tmplPath, outPath)
	}
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
		// Include stderr in output for debugging
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(exitErr.Stderr)
		}
		return fmt.Sprintf("ERROR: %v", err)
	}
	return string(output)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
