package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCLIArgumentErrors проверяет внешний контракт ошибок аргументов командной строки.
func TestCLIArgumentErrors(t *testing.T) {
	binary := buildBinary(t)
	tests := []struct {
		name      string
		arguments []string
		message   string
	}{
		{name: "missing path", message: "не указан обязательный путь"},
		{name: "extra path", arguments: []string{"one.txt", "two.txt"}, message: "ожидался ровно один путь"},
		{name: "unknown flag", arguments: []string{"--unknown", "one.txt"}, message: "недопустимые аргументы"},
		{name: "negative depth", arguments: []string{"--depth", "-1", "one.txt"}, message: "глубина не может быть отрицательной"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			workingDirectory := t.TempDir()
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			command := exec.Command(binary, test.arguments...)
			command.Dir = workingDirectory
			command.Stdout = &stdout
			command.Stderr = &stderr
			err := command.Run()

			var exitError *exec.ExitError
			if !errors.As(err, &exitError) || exitError.ExitCode() != 2 {
				t.Fatalf("exit error = %v, want code 2", err)
			}
			if stdout.Len() != 0 {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
			if !strings.Contains(stderr.String(), test.message) {
				t.Fatalf("stderr = %q, want substring %q", stderr.String(), test.message)
			}
			entries, readErr := os.ReadDir(workingDirectory)
			if readErr != nil {
				t.Fatal(readErr)
			}
			if len(entries) != 0 {
				t.Fatalf("argument error created files: %v", entries)
			}
		})
	}
}

// buildBinary собирает отдельный исполняемый файл для интеграционной проверки.
func buildBinary(t *testing.T) string {
	t.Helper()
	name := "dslparser"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join(t.TempDir(), name)
	command := exec.Command("go", "build", "-o", path, ".")
	command.Dir = "."
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("go build: %v\n%s", err, output)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("built binary: %v", err)
	}
	return path
}
