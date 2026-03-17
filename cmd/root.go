package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:     "aix [prompt]",
	Short:   "AI eXchange - delegate tasks to AI CLIs",
	Long:    "aix delegates prompts to AI CLI tools (Codex, Claude, Gemini, Kiro) with context file injection.",
	Version: Version,
	Args:    cobra.ArbitraryArgs,
	RunE:    runAsk,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("aix v%s\n", Version))

	// ask flags on root command (aix "prompt" is shorthand for aix ask "prompt")
	rootCmd.Flags().StringP("context", "c", "", "context file path (developer-instructions replacement)")
	rootCmd.Flags().StringArrayP("file", "f", nil, "attach file(s) to prompt (repeatable)")
	rootCmd.Flags().StringP("model", "m", "", "model override")
	rootCmd.Flags().StringP("sandbox", "s", "read-only", "sandbox mode (read-only, workspace-write, danger-full-access)")
	rootCmd.Flags().String("cwd", "", "working directory for codex")
	rootCmd.Flags().Bool("raw", false, "passthrough codex stdout/stderr without processing")
	rootCmd.Flags().IntP("timeout", "t", 300, "timeout in seconds")
}

// runAsk is a placeholder - will be implemented in Phase 4
func runAsk(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("prompt required\n\nUsage: aix \"your prompt here\" [flags]")
	}
	fmt.Fprintf(os.Stderr, "aix v%s - not yet implemented (Phase 4)\n", Version)
	return nil
}
