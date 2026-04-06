package cli

import (
	"encoding/json"
	"fmt"

	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/internal/rpc"
)

// RunStateCommand queries game state via a registered state exporter.
func (e *CommandExecutor) RunStateCommand(name, path string) error {
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

	e.writer.Success(result.Response)
	return nil
}