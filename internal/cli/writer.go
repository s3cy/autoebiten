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

// PrintDiff prints a diff output if present.
func (w *Writer) PrintDiff(diff string) {
	if diff != "" {
		fmt.Println(diff)
	}
}

// SuccessWithDiff prints a success message with optional diff output.
func (w *Writer) SuccessWithDiff(message, diff string) {
	w.PrintDiff(diff)
	w.Success(message)
}
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
