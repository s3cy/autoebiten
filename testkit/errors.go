package testkit

import "errors"

// ErrPathNotFound is returned when a state query path cannot be resolved.
// This occurs when:
//   - The path references a non-existent field
//   - An array index is out of bounds
//   - A map key does not exist
//   - The path traverses a nil pointer
var ErrPathNotFound = errors.New("path not found")

// ErrGameNotRunning is returned when an operation is attempted on a game
// that has not been launched or has already shut down.
var ErrGameNotRunning = errors.New("game is not running")

// ErrTimeout is returned when an operation times out.
var ErrTimeout = errors.New("operation timed out")

// ErrInvalidState is returned when the game is in an invalid state for the
// requested operation.
var ErrInvalidState = errors.New("invalid game state")
