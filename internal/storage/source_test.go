package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestReadSource проверяет точные байты и производные сведения источника.
func TestReadSource(t *testing.T) {
	t.Run("точные байты и производные сведения", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "nested", "..", "lesson.txt")
		data := []byte("\xef\xbb\xbfстрока\r\n")
		writeTestFile(t, filepath.Join(root, "lesson.txt"), data)

		source, err := (FileSystem{}).ReadSource(path)
		if err != nil {
			t.Fatalf("ReadSource() error: %v", err)
		}
		sum := sha256.Sum256(data)
		if source.Path != filepath.Join(root, "lesson.txt") || source.FileName != "lesson.txt" || source.ByteLength != len(data) || source.SHA256 != hex.EncodeToString(sum[:]) {
			t.Fatalf("ReadSource() = %#v", source)
		}
		if !reflect.DeepEqual(source.Bytes, data) {
			t.Fatalf("Bytes = %q, want %q", source.Bytes, data)
		}

		source.Bytes[0] = 0
		readBack, err := os.ReadFile(filepath.Join(root, "lesson.txt"))
		if err != nil || !reflect.DeepEqual(readBack, data) {
			t.Fatalf("returned bytes alias file data: %q, %v", readBack, err)
		}
	})

	t.Run("пустой файл", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "empty.txt")
		writeTestFile(t, path, nil)
		source, err := (FileSystem{}).ReadSource(path)
		if err != nil {
			t.Fatalf("ReadSource() error: %v", err)
		}
		if source.ByteLength != 0 || source.SHA256 != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" || len(source.Bytes) != 0 {
			t.Fatalf("ReadSource(empty) = %#v", source)
		}
	})

	for _, test := range []struct {
		name string
		path func(*testing.T) string
	}{
		{name: "каталог", path: func(t *testing.T) string { return t.TempDir() }},
		{name: "отсутствующий файл", path: func(t *testing.T) string { return filepath.Join(t.TempDir(), "missing.txt") }},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := (FileSystem{}).ReadSource(test.path(t)); err == nil {
				t.Fatal("ReadSource() error = nil")
			}
		})
	}
}

// TestTargetPath проверяет замену только последнего расширения.
func TestTargetPath(t *testing.T) {
	root := t.TempDir()
	tests := []struct{ input, want string }{
		{input: filepath.Join(root, "lesson.txt"), want: filepath.Join(root, "lesson.json")},
		{input: filepath.Join(root, "lesson.TxT"), want: filepath.Join(root, "lesson.json")},
		{input: filepath.Join(root, "lesson.part.txt"), want: filepath.Join(root, "lesson.part.json")},
		{input: filepath.Join(root, "lesson"), want: filepath.Join(root, "lesson.json")},
		{input: filepath.Join(root, "child", "..", "lesson.txt"), want: filepath.Join(root, "lesson.json")},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := TargetPath(test.input)
			if err != nil {
				t.Fatalf("TargetPath() error: %v", err)
			}
			if got != test.want || !filepath.IsAbs(got) {
				t.Fatalf("TargetPath(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
	for _, input := range []string{"", string([]byte{0})} {
		if _, err := TargetPath(input); err == nil {
			t.Fatalf("TargetPath(%q) error = nil", input)
		}
	}
}

// writeTestFile записывает файл и завершает тест при ошибке.
func writeTestFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestSourceErrorsRemainInspectable проверяет сохранение системной причины ошибки чтения.
func TestSourceErrorsRemainInspectable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.txt")
	_, err := (FileSystem{}).ReadSource(path)
	if err == nil {
		t.Fatal("ReadSource() error = nil")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ReadSource() error = %v, want os.ErrNotExist", err)
	}
}
