package docgen

import (
	"fmt"
)

// VerifyOutputs checks that all provided outputs are identical.
// Returns an error if any outputs differ.
func VerifyOutputs(outputs ...string) error {
	if len(outputs) == 0 {
		return nil
	}

	first := outputs[0]
	for i, output := range outputs {
		if output != first {
			return fmt.Errorf("output mismatch: output[0]=%q != output[%d]=%q", first, i, output)
		}
	}

	return nil
}