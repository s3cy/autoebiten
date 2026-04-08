package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

// Writer handles CLI output formatting.
type Writer struct {
	encoder *json.Encoder
}

// NewWriter creates a new Writer.
func NewWriter() *Writer {
	return &Writer{
		encoder: json.NewEncoder(os.Stdout),
	}
}

// PrintJSON prints a value as JSON.
func (w *Writer) PrintJSON(v any) error {
	return w.encoder.Encode(v)
}

// Success prints a success message.
func (w *Writer) Success(message string) {
	fmt.Fprintln(os.Stdout, "OK:", message)
}
