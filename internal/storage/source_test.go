package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestReadSource проверяет точные байты и производные сведения источника.
func TestReadSource(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "nested", "..", "lesson.txt")
	data := []byte("\xef\xbb\xbfстрока\r\n")
	if err := os.WriteFile(filepath.Join(root, "lesson.txt"), data, 0o644); err != nil {
		t.Fatal(err)
	}
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
}

// TestReadSourceRejectsNonRegularFile проверяет отказ от каталогов.
func TestReadSourceRejectsNonRegularFile(t *testing.T) {
	_, err := (FileSystem{}).ReadSource(t.TempDir())
	if err == nil {
		t.Fatal("ReadSource(directory) error = nil")
	}
}

// TestTargetPath проверяет замену только последнего расширения.
func TestTargetPath(t *testing.T) {
	tests := []struct{ input, suffix string }{
		{input: "lesson.txt", suffix: "lesson.json"},
		{input: "lesson.TxT", suffix: "lesson.json"},
		{input: "lesson.part.txt", suffix: "lesson.part.json"},
		{input: "lesson", suffix: "lesson.json"},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := TargetPath(test.input)
			if err != nil {
				t.Fatalf("TargetPath() error: %v", err)
			}
			if !strings.HasSuffix(got, test.suffix) || !filepath.IsAbs(got) {
				t.Fatalf("TargetPath(%q) = %q", test.input, got)
			}
		})
	}
	if _, err := TargetPath(""); err == nil {
		t.Fatal("TargetPath(empty) error = nil")
	}
}
