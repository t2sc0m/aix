package adapter

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/t2sc0m/aix/prompt"
	"github.com/t2sc0m/aix/runner"
)

// CodexAdapter implements Adapter for OpenAI Codex CLI.
type CodexAdapter struct {
	runner runner.Runner
}

func NewCodexAdapter(r runner.Runner) *CodexAdapter {
	return &CodexAdapter{runner: r}
}

func (a *CodexAdapter) Name() string {
	return "codex"
}

func (a *CodexAdapter) IsInstalled() bool {
	_, err := exec.LookPath("codex")
	return err == nil
}

func (a *CodexAdapter) AuthStatus() AuthInfo {
	home, err := os.UserHomeDir()
	if err != nil {
		return AuthInfo{Authenticated: false, Detail: "unknown"}
	}
	configPath := filepath.Join(home, ".codex", "config.toml")
	if _, err := os.Stat(configPath); err != nil {
		return AuthInfo{Authenticated: false, Detail: "not authenticated"}
	}
	return AuthInfo{Authenticated: true, Detail: "config found"}
}

func (a *CodexAdapter) Send(ctx context.Context, req *Request) (*Response, error) {
	assembled, err := prompt.Build(req.Prompt, req.Context, req.Files)
	if err != nil {
		return nil, fmt.Errorf("prompt build: %w", err)
	}

	if req.Raw {
		return a.sendRaw(ctx, req, assembled)
	}
	return a.sendWithOutput(ctx, req, assembled)
}

func (a *CodexAdapter) sendWithOutput(ctx context.Context, req *Request, assembled string) (*Response, error) {
	tmpFile, err := os.CreateTemp("", "aix-*.out")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	args := a.buildArgs(req, assembled)
	args = append(args, "-o", tmpPath)

	result, err := a.runner.Run(ctx, "codex", args, nil)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return nil, fmt.Errorf("timeout after %ds", req.Timeout)
		}
		return nil, fmt.Errorf("codex exec: %w", err)
	}

	if result.ExitCode != 0 {
		stderrMsg := strings.TrimSpace(string(result.Stderr))
		if stderrMsg == "" {
			stderrMsg = fmt.Sprintf("exit code %d", result.ExitCode)
		}
		return &Response{Content: stderrMsg, ExitCode: result.ExitCode}, nil
	}

	content, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("read output: %w", err)
	}

	return &Response{
		Content:  strings.TrimSpace(string(content)),
		ExitCode: 0,
	}, nil
}

func (a *CodexAdapter) sendRaw(ctx context.Context, req *Request, assembled string) (*Response, error) {
	args := a.buildArgs(req, assembled)

	result, err := a.runner.Run(ctx, "codex", args, nil)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return nil, fmt.Errorf("timeout after %ds", req.Timeout)
		}
		return nil, fmt.Errorf("codex exec: %w", err)
	}

	output := string(result.Stdout)
	if len(result.Stderr) > 0 {
		output += string(result.Stderr)
	}

	return &Response{
		Content:  output,
		ExitCode: result.ExitCode,
	}, nil
}

func (a *CodexAdapter) buildArgs(req *Request, assembled string) []string {
	args := []string{"exec", "--ephemeral"}

	if req.Model != "" {
		args = append(args, "-m", req.Model)
	}

	sandbox := req.Sandbox
	if sandbox == "" {
		sandbox = "read-only"
	}
	args = append(args, "-s", sandbox)

	if req.Cwd != "" {
		args = append(args, "-C", req.Cwd)
	}

	args = append(args, assembled)
	return args
}
