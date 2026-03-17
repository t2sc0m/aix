package prompt

import (
	"fmt"
	"strings"
)

const (
	MaxContextSize   = 100 * 1024 // 100KB
	MaxFileSize      = 100 * 1024 // 100KB per file
	MaxAssembledSize = 500 * 1024 // 500KB total
)

// File represents an attached file with its content.
type File struct {
	Name    string
	Content string
}

// Build assembles a prompt from context, user prompt, and attached files.
func Build(userPrompt, context string, files []File) (string, error) {
	if err := validateContext(context); err != nil {
		return "", err
	}
	if err := validateFiles(files); err != nil {
		return "", err
	}

	var b strings.Builder

	if context != "" {
		b.WriteString(context)
		b.WriteString("\n\n---\n\n")
	}

	b.WriteString(userPrompt)

	for _, f := range files {
		b.WriteString("\n\n---\nFile: ")
		b.WriteString(f.Name)
		b.WriteString("\n```\n")
		b.WriteString(f.Content)
		if !strings.HasSuffix(f.Content, "\n") {
			b.WriteString("\n")
		}
		b.WriteString("```")
	}

	assembled := b.String()
	if len(assembled) > MaxAssembledSize {
		return "", fmt.Errorf("assembled prompt too large (%dKB, max %dKB)", len(assembled)/1024, MaxAssembledSize/1024)
	}

	return assembled, nil
}

func validateContext(context string) error {
	if len(context) > MaxContextSize {
		return fmt.Errorf("context file too large (%dKB, max %dKB)", len(context)/1024, MaxContextSize/1024)
	}
	return nil
}

func validateFiles(files []File) error {
	for _, f := range files {
		if len(f.Content) > MaxFileSize {
			return fmt.Errorf("file too large: %s (%dKB, max %dKB)", f.Name, len(f.Content)/1024, MaxFileSize/1024)
		}
		if isBinary(f.Content) {
			return fmt.Errorf("binary file not supported: %s", f.Name)
		}
	}
	return nil
}

// isBinary checks for null bytes in the first 8KB.
func isBinary(content string) bool {
	check := content
	if len(check) > 8192 {
		check = check[:8192]
	}
	return strings.ContainsRune(check, '\x00')
}
