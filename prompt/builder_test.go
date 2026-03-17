package prompt

import (
	"strings"
	"testing"
)

func TestBuild_PromptOnly(t *testing.T) {
	result, err := Build("hello world", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello world" {
		t.Errorf("expected 'hello world', got %q", result)
	}
}

func TestBuild_WithContext(t *testing.T) {
	result, err := Build("review this", "You are an expert reviewer.", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(result, "You are an expert reviewer.") {
		t.Error("context should be at the beginning")
	}
	if !strings.Contains(result, "---") {
		t.Error("should contain separator")
	}
	if !strings.Contains(result, "review this") {
		t.Error("should contain user prompt")
	}
}

func TestBuild_WithFiles(t *testing.T) {
	files := []File{
		{Name: "plan.md", Content: "# Plan\nStep 1\n"},
	}
	result, err := Build("review this", "", files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "File: plan.md") {
		t.Error("should contain filename")
	}
	if !strings.Contains(result, "```\n# Plan") {
		t.Error("should contain file content in code block")
	}
}

func TestBuild_FullCombination(t *testing.T) {
	files := []File{
		{Name: "a.go", Content: "package main\n"},
		{Name: "b.go", Content: "package util\n"},
	}
	result, err := Build("review these files", "You are a code reviewer.", files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(result, "You are a code reviewer.") {
		t.Error("context first")
	}
	if !strings.Contains(result, "review these files") {
		t.Error("prompt present")
	}
	if !strings.Contains(result, "File: a.go") || !strings.Contains(result, "File: b.go") {
		t.Error("both files present")
	}
}

func TestBuild_ContextTooLarge(t *testing.T) {
	bigContext := strings.Repeat("x", MaxContextSize+1)
	_, err := Build("prompt", bigContext, nil)
	if err == nil {
		t.Fatal("expected error for large context")
	}
	if !strings.Contains(err.Error(), "context file too large") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuild_FileTooLarge(t *testing.T) {
	files := []File{
		{Name: "big.txt", Content: strings.Repeat("x", MaxFileSize+1)},
	}
	_, err := Build("prompt", "", files)
	if err == nil {
		t.Fatal("expected error for large file")
	}
	if !strings.Contains(err.Error(), "file too large") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuild_BinaryFile(t *testing.T) {
	files := []File{
		{Name: "image.png", Content: "PNG\x00\x00data"},
	}
	_, err := Build("prompt", "", files)
	if err == nil {
		t.Fatal("expected error for binary file")
	}
	if !strings.Contains(err.Error(), "binary file not supported") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuild_AssembledTooLarge(t *testing.T) {
	fileContent := strings.Repeat("x", MaxFileSize-10)
	var files []File
	for i := range 6 {
		files = append(files, File{
			Name:    strings.Repeat("f", i+1) + ".txt",
			Content: fileContent,
		})
	}
	_, err := Build("prompt", "", files)
	if err == nil {
		t.Fatal("expected error for assembled too large")
	}
	if !strings.Contains(err.Error(), "assembled prompt too large") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuild_EmptyFile(t *testing.T) {
	files := []File{
		{Name: "empty.txt", Content: ""},
	}
	result, err := Build("check this", "", files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "File: empty.txt") {
		t.Error("empty file should still be included")
	}
}
