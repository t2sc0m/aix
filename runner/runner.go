package runner

import (
	"context"
	"io"
)

// Result holds the output of a subprocess execution.
type Result struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
}

// Runner abstracts subprocess execution for testability.
type Runner interface {
	Run(ctx context.Context, name string, args []string, stdin io.Reader) (*Result, error)
}
