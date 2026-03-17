package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

// ExecRunner implements Runner using os/exec.
type ExecRunner struct{}

func NewExecRunner() *ExecRunner {
	return &ExecRunner{}
}

func (r *ExecRunner) Run(ctx context.Context, name string, args []string, stdin io.Reader) (*Result, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if stdin != nil {
		cmd.Stdin = stdin
	}

	err := cmd.Run()

	result := &Result{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
	}

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			result.ExitCode = exitErr.ExitCode()
			return result, nil
		}
		if ctx.Err() == context.DeadlineExceeded {
			return result, fmt.Errorf("timeout: %w", ctx.Err())
		}
		return result, fmt.Errorf("exec failed: %w", err)
	}

	result.ExitCode = 0
	return result, nil
}
