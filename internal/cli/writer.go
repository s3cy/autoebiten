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

// Println prints a line.
func (w *Writer) Println(a ...any) {
	fmt.Fprintln(os.Stdout, a...)
}

// Errorln prints an error line to stderr.
func (w *Writer) Errorln(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

// Success prints a success message.
func (w *Writer) Success(message string) {
	w.Println("OK:", message)
}

// Failure prints a failure message.
func (w *Writer) Failure(message string) {
	w.Errorln("ERROR:", message)
}

// ExitWithError prints an error and exits with code 1.
func (w *Writer) ExitWithError(message string) {
	w.Failure(message)
	os.Exit(1)
}
