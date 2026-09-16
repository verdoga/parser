package discovery

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

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
			result, err := Discover(Request{Path: root, Depth: test.depth})
			if err != nil {
				t.Fatalf("Discover() error: %v", err)
			}
			if !reflect.DeepEqual(result.Files, test.want) {
				t.Fatalf("Discover() files = %q, want %q", result.Files, test.want)
			}
		})
	}
}

func TestDiscoverExplicitFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "LESSON.TxT")
	writeFile(t, path)
	depth := 0

	result, err := Discover(Request{Path: path, Depth: &depth})
	if err != nil {
		t.Fatalf("Discover() error: %v", err)
	}
	if !result.IsFile || !reflect.DeepEqual(result.Files, []string{path}) {
		t.Fatalf("Discover() = %#v", result)
	}
}

func TestDiscoverRejectsExplicitNonTXT(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lesson.md")
	writeFile(t, path)

	_, err := Discover(Request{Path: path})
	if err == nil {
		t.Fatal("Discover() error = nil")
	}
	if _, ok := err.(*InputError); !ok {
		t.Fatalf("Discover() error type = %T, want *InputError", err)
	}
}

func TestDiscoverIgnoresSymbolicLinks(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.txt")
	writeFile(t, target)
	link := filepath.Join(root, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symbolic links are unavailable: %v", err)
	}

	result, err := Discover(Request{Path: root})
	if err != nil {
		t.Fatalf("Discover() error: %v", err)
	}
	if !reflect.DeepEqual(result.Files, []string{target}) {
		t.Fatalf("Discover() files = %q, want [%q]", result.Files, target)
	}
}

func writeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}
}

func absolutePaths(root string, paths ...string) []string {
	result := make([]string, len(paths))
	for index, path := range paths {
		result[index] = filepath.Join(root, path)
	}
	return result
}

func intPointer(value int) *int { return &value }
