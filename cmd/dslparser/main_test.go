package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseArguments проверяет допустимые параметры и сообщения об ошибках командной строки.
func TestParseArguments(t *testing.T) {
	tests := []struct {
		name        string
		arguments   []string
		wantPath    string
		wantReplace bool
		wantDepth   *int
		wantError   string
	}{
		{name: "path", arguments: []string{"lesson.txt"}, wantPath: "lesson.txt"},
		{name: "replace", arguments: []string{"--replace", "lesson.txt"}, wantPath: "lesson.txt", wantReplace: true},
		{name: "depth zero", arguments: []string{"--depth", "0", "lessons"}, wantPath: "lessons", wantDepth: intPointer(0)},
		{name: "all options", arguments: []string{"--replace", "--depth=2", "lessons"}, wantPath: "lessons", wantReplace: true, wantDepth: intPointer(2)},
		{name: "missing path", wantError: "не указан обязательный путь"},
		{name: "extra path", arguments: []string{"one", "two"}, wantError: "ожидался ровно один путь"},
		{name: "unknown flag", arguments: []string{"--unknown", "one"}, wantError: "недопустимые аргументы"},
		{name: "missing depth", arguments: []string{"--depth"}, wantError: "недопустимые аргументы"},
		{name: "invalid depth", arguments: []string{"--depth", "deep", "one"}, wantError: "глубина должна быть целым числом"},
		{name: "negative depth", arguments: []string{"--depth", "-1", "one"}, wantError: "глубина не может быть отрицательной"},
		{name: "help is not public option", arguments: []string{"--help"}, wantError: "недопустимые аргументы"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			options, err := parseArguments(test.arguments)
			if test.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantError) {
					t.Fatalf("parseArguments() error = %v, want substring %q", err, test.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseArguments() unexpected error: %v", err)
			}
			if options.Path != test.wantPath || options.Replace != test.wantReplace {
				t.Fatalf("parseArguments() = %#v", options)
			}
			assertDepth(t, options.Depth, test.wantDepth)
		})
	}
}

// TestRunRejectsArgumentsBeforeAccessingPath проверяет приоритет валидации аргументов.
func TestRunRejectsArgumentsBeforeAccessingPath(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{"first.txt", "second.txt"}, &stdout, &stderr)

	if exitCode != 2 {
		t.Fatalf("run() exit code = %d, want 2", exitCode)
	}
	if stdout.Len() != 0 {
		t.Fatalf("run() stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "ожидался ровно один путь") {
		t.Fatalf("run() stderr = %q", stderr.String())
	}
}

// TestRunEmptyDirectory проверяет успешное завершение для каталога без входных файлов.
func TestRunEmptyDirectory(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{filepath.Clean(t.TempDir())}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("run() exit code = %d, want 0; stderr: %s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "ИТОГ found=0") {
		t.Fatalf("run() stdout = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("run() stderr = %q, want empty", stderr.String())
	}
}

// intPointer возвращает указатель на отдельную копию целого числа.
func intPointer(value int) *int { return &value }

// assertDepth сравнивает необязательные значения глубины.
func assertDepth(t *testing.T, got, want *int) {
	t.Helper()
	if got == nil || want == nil {
		if got != want {
			t.Fatalf("depth = %v, want %v", got, want)
		}
		return
	}
	if *got != *want {
		t.Fatalf("depth = %d, want %d", *got, *want)
	}
}
