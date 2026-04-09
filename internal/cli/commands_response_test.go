package cli

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/s3cy/autoebiten/internal/rpc"
)

// captureOutput captures stdout and stderr during test execution.
func captureOutput(f func()) (stdout, stderr string) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rStdout, wStdout, _ := os.Pipe()
	rStderr, wStderr, _ := os.Pipe()
	os.Stdout = wStdout
	os.Stderr = wStderr

	f()

	wStdout.Close()
	wStderr.Close()
	outBytes, _ := io.ReadAll(rStdout)
	errBytes, _ := io.ReadAll(rStderr)
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	return string(outBytes), string(errBytes)
}

// TestHandleResponseWithDiff verifies diff output and success callback.
func TestHandleResponseWithDiff(t *testing.T) {
	executor := NewCommandExecutor()

	resp := &rpc.RPCResponse{
		Result: json.RawMessage(`{"status":"ok"}`),
		Extra: map[string]interface{}{
			"diff": "line1\nline2",
		},
	}

	successCalled := false
	onSuccess := func() {
		successCalled = true
	}

	stdout, stderr := captureOutput(func() {
		executor.handleResponse(resp, onSuccess)
	})

	// Verify diff output
	expectedDiff := "<log_diff>\nline1\nline2</log_diff>\n"
	if stdout != expectedDiff {
		t.Errorf("stdout = %q, want %q", stdout, expectedDiff)
	}

	// Verify no stderr
	if stderr != "" {
		t.Errorf("stderr = %q, want empty", stderr)
	}

	// Verify success callback was called
	if !successCalled {
		t.Error("onSuccess callback was not called")
	}
}

// TestHandleResponseWithProxyError verifies proxy_error output and error handling.
func TestHandleResponseWithProxyError(t *testing.T) {
	executor := NewCommandExecutor()

	resp := &rpc.RPCResponse{
		Result: json.RawMessage(`{"status":"ok"}`),
		Extra: map[string]interface{}{
			"proxy_error": "crash detected: segmentation fault",
		},
	}

	successCalled := false
	onSuccess := func() {
		successCalled = true
	}

	stdout, stderr := captureOutput(func() {
		executor.handleResponse(resp, onSuccess)
	})

	// Verify proxy_error output
	expectedProxyError := "<proxy_error>\ncrash detected: segmentation fault\n</proxy_error>\n"
	if stdout != expectedProxyError {
		t.Errorf("stdout = %q, want %q", stdout, expectedProxyError)
	}

	// Verify no stderr (no RPC error)
	if stderr != "" {
		t.Errorf("stderr = %q, want empty", stderr)
	}

	// Verify success callback was called
	if !successCalled {
		t.Error("onSuccess callback was not called")
	}
}

// TestHandleResponseEmptyDiff verifies empty diff produces no output.
func TestHandleResponseEmptyDiff(t *testing.T) {
	executor := NewCommandExecutor()

	resp := &rpc.RPCResponse{
		Result: json.RawMessage(`{"status":"ok"}`),
		Extra: map[string]interface{}{
			"diff": "",
		},
	}

	successCalled := false
	onSuccess := func() {
		successCalled = true
	}

	stdout, _ := captureOutput(func() {
		executor.handleResponse(resp, onSuccess)
	})

	// Verify empty diff produces no output
	expectedDiff := ""
	if stdout != expectedDiff {
		t.Errorf("stdout = %q, want %q", stdout, expectedDiff)
	}

	if !successCalled {
		t.Error("onSuccess callback was not called")
	}
}

// TestHandleResponseErrorOnly verifies error without diff/proxy_error.
func TestHandleResponseErrorOnly(t *testing.T) {
	executor := NewCommandExecutor()

	resp := &rpc.RPCResponse{
		Error: &rpc.RPCError{
			Code:    -32600,
			Message: "Invalid Request",
		},
	}

	successCalled := false
	onSuccess := func() {
		successCalled = true
	}

	stdout, stderr := captureOutput(func() {
		executor.handleResponse(resp, onSuccess)
	})

	// Verify no stdout (no diff or proxy_error)
	if stdout != "" {
		t.Errorf("stdout = %q, want empty", stdout)
	}

	// Verify error to stderr
	expectedError := "Error: Invalid Request\n"
	if stderr != expectedError {
		t.Errorf("stderr = %q, want %q", stderr, expectedError)
	}

	// Verify success callback was NOT called
	if successCalled {
		t.Error("onSuccess callback was called when it should not have been")
	}
}

// TestHandleResponseBothDiffAndProxyError verifies both outputs when present.
func TestHandleResponseBothDiffAndProxyError(t *testing.T) {
	executor := NewCommandExecutor()

	resp := &rpc.RPCResponse{
		Result: json.RawMessage(`{"status":"ok"}`),
		Extra: map[string]interface{}{
			"diff":        "some log output",
			"proxy_error": "game crashed",
		},
	}

	successCalled := false
	onSuccess := func() {
		successCalled = true
	}

	stdout, stderr := captureOutput(func() {
		executor.handleResponse(resp, onSuccess)
	})

	// Verify both outputs
	expected := "<log_diff>\nsome log output</log_diff>\n<proxy_error>\ngame crashed\n</proxy_error>\n"
	if stdout != expected {
		t.Errorf("stdout = %q, want %q", stdout, expected)
	}

	if stderr != "" {
		t.Errorf("stderr = %q, want empty", stderr)
	}

	if !successCalled {
		t.Error("onSuccess callback was not called")
	}
}

// TestHandleResponseErrorWithDiffAndProxyError verifies error case still outputs diff and proxy_error.
func TestHandleResponseErrorWithDiffAndProxyError(t *testing.T) {
	executor := NewCommandExecutor()

	resp := &rpc.RPCResponse{
		Error: &rpc.RPCError{
			Code:    -32603,
			Message: "Internal error",
		},
		Extra: map[string]interface{}{
			"diff":        "log before error",
			"proxy_error": "crash during processing",
		},
	}

	successCalled := false
	onSuccess := func() {
		successCalled = true
	}

	stdout, stderr := captureOutput(func() {
		executor.handleResponse(resp, onSuccess)
	})

	// Verify diff and proxy_error still output even on error
	expected := "<log_diff>\nlog before error</log_diff>\n<proxy_error>\ncrash during processing\n</proxy_error>\n"
	if stdout != expected {
		t.Errorf("stdout = %q, want %q", stdout, expected)
	}

	// Verify error to stderr
	expectedError := "Error: Internal error\n"
	if stderr != expectedError {
		t.Errorf("stderr = %q, want %q", stderr, expectedError)
	}

	// Verify success callback was NOT called
	if successCalled {
		t.Error("onSuccess callback was called when it should not have been")
	}
}

// TestHandleResponseNilExtra verifies handling when Extra is nil.
func TestHandleResponseNilExtra(t *testing.T) {
	executor := NewCommandExecutor()

	resp := &rpc.RPCResponse{
		Result: json.RawMessage(`{"status":"ok"}`),
		Extra:  nil,
	}

	successCalled := false
	onSuccess := func() {
		successCalled = true
	}

	stdout, stderr := captureOutput(func() {
		executor.handleResponse(resp, onSuccess)
	})

	// Verify no output when Extra is nil
	if stdout != "" {
		t.Errorf("stdout = %q, want empty", stdout)
	}

	if stderr != "" {
		t.Errorf("stderr = %q, want empty", stderr)
	}

	if !successCalled {
		t.Error("onSuccess callback was not called")
	}
}
