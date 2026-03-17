package adapter

import (
	"context"

	"github.com/t2sc0m/aix/prompt"
)

// AuthInfo holds authentication status for an adapter.
type AuthInfo struct {
	Authenticated bool
	Detail        string // "config found", "not authenticated", "unknown"
}

// Request holds all parameters for an AI delegation call.
type Request struct {
	Prompt  string
	Context string       // -c file content (developer-instructions replacement)
	Files   []prompt.File // -f attached files
	Model   string
	Sandbox string // default: read-only
	Cwd     string
	Raw     bool
	Timeout int // seconds
}

// Response holds the result of an AI delegation call.
type Response struct {
	Content  string
	ExitCode int
}

// Adapter defines the interface for AI CLI backends.
type Adapter interface {
	Name() string
	IsInstalled() bool
	AuthStatus() AuthInfo
	Send(ctx context.Context, req *Request) (*Response, error)
}
