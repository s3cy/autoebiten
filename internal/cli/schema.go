package cli

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/s3cy/autoebiten/internal/script"
)

// GenerateSchema generates a JSON Schema from the Go structs.
func GenerateSchema() ([]byte, error) {
	reflector := &jsonschema.Reflector{
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
	}

	schema := reflector.Reflect(&script.ScriptSchema{})

	// Add schema version and title
	schema.Version = jsonschema.Version
	schema.Title = "AutoEbiten Script"
	schema.Description = "JSON script format for automating Ebitengine games"

	// Enhance definitions with better metadata
	if schema.Definitions != nil {
		// Enhance CommandSchema with better descriptions
		if cw, ok := schema.Definitions["CommandSchema"]; ok {
			cw.Title = "Command"
			cw.Description = "A single command to execute"
		}

		// Note: InputCmd, MouseCmd, WheelCmd, ScreenshotCmd, DelayCmd
		// are enhanced via jsonschema struct tags in ast.go
	}

	return json.MarshalIndent(schema, "", "  ")
}

// PrintSchema prints the JSON schema to stdout.
func PrintSchema() error {
	data, err := GenerateSchema()
	if err != nil {
		return fmt.Errorf("failed to generate schema: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
