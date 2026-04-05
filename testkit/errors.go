package testkit

import "errors"

// ErrGameNotRunning is returned when an operation is attempted on a game
// that has not been launched or has already shut down.
var ErrGameNotRunning = errors.New("game is not running")

// ErrTimeout is returned when an operation times out.
var ErrTimeout = errors.New("operation timed out")

// ErrInvalidState is returned when the game is in an invalid state for the
// requested operation.
var ErrInvalidState = errors.New("invalid game state")
