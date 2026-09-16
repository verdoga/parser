package discovery

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// TestDiscoverDirectoryDepthAndOrder проверяет глубину обхода и порядок результатов.
func TestDiscoverDirectoryDepthAndOrder(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "z.TXT"))
	writeFile(t, filepath.Join(root, "a.txt"))
	writeFile(t, filepath.Join(root, "ignore.md"))
	writeFile(t, filepath.Join(root, "first", "b.txt"))
	writeFile(t, filepath.Join(root, "first", "second", "c.txt"))

	tests := []struct {
		name  string
		depth *int
		want  []string
	}{
		{name: "unlimited", want: absolutePaths(root, "a.txt", filepath.Join("first", "b.txt"), filepath.Join("first", "second", "c.txt"), "z.TXT")},
		{name: "zero", depth: intPointer(0), want: absolutePaths(root, "a.txt", "z.TXT")},
		{name: "one", depth: intPointer(1), want: absolutePaths(root, "a.txt", filepath.Join("first", "b.txt"), "z.TXT")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			files, scanErrors, err := Discover(root, test.depth)
			if err != nil {
				t.Fatalf("Discover() error: %v", err)
			}
			if len(scanErrors) != 0 {
				t.Fatalf("Discover() scan errors: %v", scanErrors)
			}
			if !slices.Equal(files, test.want) {
				t.Fatalf("Discover() files = %q, want %q", files, test.want)
			}
		})
	}
}

// TestDiscoverExplicitFile проверяет регистронезависимое расширение явного файла.
func TestDiscoverExplicitFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "LESSON.TxT")
	writeFile(t, path)
	depth := 0

	files, scanErrors, err := Discover(path, &depth)
	if err != nil {
		t.Fatalf("Discover() error: %v", err)
	}
	if len(scanErrors) != 0 {
		t.Fatalf("Discover() scan errors: %v", scanErrors)
	}
	if !slices.Equal(files, []string{path}) {
		t.Fatalf("Discover() files = %q, want [%q]", files, path)
	}
}

// TestDiscoverRejectsExplicitNonTXT проверяет отказ для явного файла другого типа.
func TestDiscoverRejectsExplicitNonTXT(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lesson.md")
	writeFile(t, path)

	_, _, err := Discover(path, nil)
	if err == nil {
		t.Fatal("Discover() error = nil")
	}
	if _, ok := err.(*InputError); !ok {
		t.Fatalf("Discover() error type = %T, want *InputError", err)
	}
}

// TestDiscoverIgnoresSymbolicLinks проверяет исключение ссылок из результатов.
func TestDiscoverIgnoresSymbolicLinks(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.txt")
	writeFile(t, target)
	link := filepath.Join(root, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symbolic links are unavailable: %v", err)
	}

	files, scanErrors, err := Discover(root, nil)
	if err != nil {
		t.Fatalf("Discover() error: %v", err)
	}
	if len(scanErrors) != 0 {
		t.Fatalf("Discover() scan errors: %v", scanErrors)
	}
	if !slices.Equal(files, []string{target}) {
		t.Fatalf("Discover() files = %q, want [%q]", files, target)
	}
}

// writeFile создаёт пустой файл вместе с родительскими каталогами.
func writeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}
}

// absolutePaths строит ожидаемые абсолютные пути от общего корня.
func absolutePaths(root string, paths ...string) []string {
	result := make([]string, len(paths))
	for index, path := range paths {
		result[index] = filepath.Join(root, path)
	}
	return result
}

// intPointer возвращает указатель на отдельную копию целого числа.
func intPointer(value int) *int { return &value }
