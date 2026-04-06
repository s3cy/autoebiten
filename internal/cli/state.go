package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/internal/recording"
	"github.com/s3cy/autoebiten/internal/rpc"
	"github.com/s3cy/autoebiten/internal/script"
)

// RunStateCommand queries game state via a registered state exporter.
func (e *CommandExecutor) RunStateCommand(name, path string, shouldRecord bool) error {
	// Build custom command name with state exporter prefix
	customName := autoebiten.StateExporterPathPrefix + name

	params := &rpc.CustomParams{
		Name:    customName,
		Request: path,
	}

	req, err := rpc.BuildRequest("custom", params)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := rpc.SendRequestSocket(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("rpc error: %s", resp.Error.Message)
	}

	var result rpc.CustomResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return fmt.Errorf("invalid response format: %w", err)
	}

	// Record after successful execution
	if shouldRecord {
		recorder := recording.NewRecorderFromSocket(rpc.SocketPath())
		cmd := &script.StateCmd{
			Name: name,
			Path: path,
		}
		if err := recorder.Record(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "autoebiten: recording failed: %v\n", err)
		}
	}

	e.writer.Success(result.Response)
	return nil
}