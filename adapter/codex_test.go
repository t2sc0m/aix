package adapter

import (
	"context"
	"io"
	"testing"

	"github.com/t2sc0m/aix/runner"
)

// mockRunner implements runner.Runner for testing.
type mockRunner struct {
	result *runner.Result
	err    error
	// captured args for verification
	capturedName string
	capturedArgs []string
}

func (m *mockRunner) Run(_ context.Context, name string, args []string, _ io.Reader) (*runner.Result, error) {
	m.capturedName = name
	m.capturedArgs = args
	return m.result, m.err
}

func TestCodexAdapter_Name(t *testing.T) {
	a := NewCodexAdapter(&mockRunner{})
	if a.Name() != "codex" {
		t.Errorf("expected 'codex', got %q", a.Name())
	}
}

func TestCodexAdapter_Send_RawSuccess(t *testing.T) {
	mock := &mockRunner{
		result: &runner.Result{
			Stdout:   []byte("response content"),
			Stderr:   []byte(""),
			ExitCode: 0,
		},
	}

	a := NewCodexAdapter(mock)
	req := &Request{
		Prompt:  "hello",
		Sandbox: "read-only",
		Raw:     true,
	}

	resp, err := a.Send(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Errorf("expected exit 0, got %d", resp.ExitCode)
	}
	if resp.Content != "response content" {
		t.Errorf("expected 'response content', got %q", resp.Content)
	}

	if mock.capturedName != "codex" {
		t.Errorf("expected 'codex', got %q", mock.capturedName)
	}
	if mock.capturedArgs[0] != "exec" {
		t.Errorf("expected 'exec', got %q", mock.capturedArgs[0])
	}
	if !containsArg(mock.capturedArgs, "--ephemeral") {
		t.Error("should include --ephemeral")
	}
	if !containsArg(mock.capturedArgs, "-s") {
		t.Error("should include -s sandbox")
	}
	if !containsArg(mock.capturedArgs, "--") {
		t.Error("should include -- separator before prompt")
	}
}

func TestCodexAdapter_Send_WithModel(t *testing.T) {
	mock := &mockRunner{
		result: &runner.Result{ExitCode: 0},
	}

	a := NewCodexAdapter(mock)
	req := &Request{
		Prompt: "hello",
		Model:  "o3",
		Raw:    true,
	}

	_, err := a.Send(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsArgPair(mock.capturedArgs, "-m", "o3") {
		t.Error("should include -m o3")
	}
}

func TestCodexAdapter_Send_WithCwd(t *testing.T) {
	mock := &mockRunner{
		result: &runner.Result{ExitCode: 0},
	}

	a := NewCodexAdapter(mock)
	req := &Request{
		Prompt: "hello",
		Cwd:    "/tmp/work",
		Raw:    true,
	}

	_, err := a.Send(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsArgPair(mock.capturedArgs, "-C", "/tmp/work") {
		t.Error("should include -C /tmp/work")
	}
}

func TestCodexAdapter_Send_NonZeroExit(t *testing.T) {
	mock := &mockRunner{
		result: &runner.Result{
			Stdout:   []byte(""),
			Stderr:   []byte("some error"),
			ExitCode: 1,
		},
	}

	a := NewCodexAdapter(mock)
	req := &Request{
		Prompt: "hello",
		Raw:    true,
	}

	resp, err := a.Send(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ExitCode != 1 {
		t.Errorf("expected exit 1, got %d", resp.ExitCode)
	}
}

func TestCodexAdapter_Send_Timeout(t *testing.T) {
	mock := &mockRunner{
		err: context.DeadlineExceeded,
	}

	a := NewCodexAdapter(mock)
	req := &Request{
		Prompt:  "hello",
		Timeout: 5,
		Raw:     true,
	}

	_, err := a.Send(context.Background(), req)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func containsArg(args []string, target string) bool {
	for _, a := range args {
		if a == target {
			return true
		}
	}
	return false
}

func containsArgPair(args []string, key, value string) bool {
	for i, a := range args {
		if a == key && i+1 < len(args) && args[i+1] == value {
			return true
		}
	}
	return false
}
