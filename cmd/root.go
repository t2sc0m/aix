package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/t2sc0m/aix/adapter"
	"github.com/t2sc0m/aix/config"
	"github.com/t2sc0m/aix/prompt"
	"github.com/t2sc0m/aix/runner"
)

const (
	ExitOK           = 0
	ExitGeneralError = 1
	ExitNotInstalled = 2
	ExitAuthFailure  = 3
	ExitExecFailure  = 4
	ExitTimeout      = 5
)

var Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:     "aix [prompt]",
	Short:   "AI eXchange - delegate tasks to AI CLIs",
	Long:    "aix delegates prompts to AI CLI tools (Codex, Claude, Gemini, Kiro) with context file injection.",
	Version: Version,
	Args:    cobra.ArbitraryArgs,
	RunE:    runAsk,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if exitErr, ok := err.(*exitError); ok {
			fmt.Fprintln(os.Stderr, exitErr.Error())
			os.Exit(exitErr.code)
		}
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(ExitGeneralError)
	}
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("aix v%s\n", Version))

	rootCmd.Flags().StringP("context", "c", "", "context file path (developer-instructions replacement)")
	rootCmd.Flags().StringArrayP("file", "f", nil, "attach file(s) to prompt (repeatable)")
	rootCmd.Flags().StringP("model", "m", "", "model override")
	rootCmd.Flags().StringP("sandbox", "s", "", "sandbox mode (read-only, workspace-write, danger-full-access)")
	rootCmd.Flags().String("cwd", "", "working directory for codex")
	rootCmd.Flags().Bool("raw", false, "passthrough codex stdout/stderr without processing")
	rootCmd.Flags().IntP("timeout", "t", 0, "timeout in seconds")
}

type exitError struct {
	code int
	msg  string
}

func (e *exitError) Error() string { return e.msg }

func runAsk(cmd *cobra.Command, args []string) error {
	// Build prompt from args or stdin
	userPrompt, err := resolvePrompt(args)
	if err != nil {
		return &exitError{ExitGeneralError, err.Error()}
	}

	// Read flags
	contextFile, _ := cmd.Flags().GetString("context")
	filePaths, _ := cmd.Flags().GetStringArray("file")
	flagModel, _ := cmd.Flags().GetString("model")
	flagSandbox, _ := cmd.Flags().GetString("sandbox")
	cwd, _ := cmd.Flags().GetString("cwd")
	raw, _ := cmd.Flags().GetBool("raw")
	flagTimeout, _ := cmd.Flags().GetInt("timeout")

	// Load config: defaults < config.yaml < CLI flags
	cfg := config.Load()
	model, sandbox, timeout := config.Resolve(cfg, flagModel, flagSandbox, flagTimeout)

	// Validate sandbox
	validSandboxes := map[string]bool{
		"read-only": true, "workspace-write": true, "danger-full-access": true,
	}
	if !validSandboxes[sandbox] {
		return &exitError{ExitGeneralError, fmt.Sprintf("invalid sandbox mode: %q (allowed: read-only, workspace-write, danger-full-access)", sandbox)}
	}

	// Read context file
	var contextContent string
	if contextFile != "" {
		data, err := os.ReadFile(contextFile)
		if err != nil {
			return &exitError{ExitGeneralError, fmt.Sprintf("cannot read context file: %v", err)}
		}
		contextContent = string(data)
	}

	// Read attached files
	var files []prompt.File
	for _, fp := range filePaths {
		data, err := os.ReadFile(fp)
		if err != nil {
			return &exitError{ExitGeneralError, fmt.Sprintf("cannot read file %s: %v", fp, err)}
		}
		files = append(files, prompt.File{Name: fp, Content: string(data)})
	}

	// Create adapter
	a := adapter.NewCodexAdapter(runner.NewExecRunner())

	// Check installation
	if !a.IsInstalled() {
		return &exitError{ExitNotInstalled, "codex is not installed. Install: npm i -g @openai/codex"}
	}

	// Check auth
	auth := a.AuthStatus()
	if !auth.Authenticated {
		return &exitError{ExitAuthFailure, fmt.Sprintf("codex auth failed: %s", auth.Detail)}
	}

	// Build request
	req := &adapter.Request{
		Prompt:  userPrompt,
		Context: contextContent,
		Files:   files,
		Model:   model,
		Sandbox: sandbox,
		Cwd:     cwd,
		Raw:     raw,
		Timeout: timeout,
	}

	// Execute with timeout
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		defer cancel()
	}

	resp, err := a.Send(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return &exitError{ExitTimeout, fmt.Sprintf("timeout after %ds", timeout)}
		}
		return &exitError{ExitExecFailure, fmt.Sprintf("codex exec failed: %v", err)}
	}

	// Output result
	fmt.Print(resp.Content)

	if resp.ExitCode != 0 {
		return &exitError{resp.ExitCode, ""}
	}
	return nil
}

// resolvePrompt reads prompt from args or stdin (pipe).
func resolvePrompt(args []string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	// Check if stdin has data (pipe)
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read stdin: %w", err)
		}
		text := strings.TrimSpace(string(data))
		if text != "" {
			return text, nil
		}
	}

	return "", fmt.Errorf("prompt required\n\nUsage: aix \"your prompt here\" [flags]")
}
