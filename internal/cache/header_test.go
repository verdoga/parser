package cache

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestReadHeader проверяет минимальное декодирование заголовка, null и игнорирование остальных полей.
func TestReadHeader(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "nested", "..", "result.json")
	writtenPath := filepath.Join(root, "result.json")
	data := `{
		"formatVersion":"1.0",
		"unknown":{"must":"be ignored"},
		"document":{
			"dslVersion":"1.2",
			"sha256":"abc123",
			"metadata":{"documentId":"Course-A","title":"ignored"},
			"lines":[{"arbitrary":true}]
		},
		"processing":[{"also":"ignored"}]
	}`
	if err := os.WriteFile(writtenPath, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	header, err := (JSONReader{}).ReadHeader(path)
	if err != nil {
		t.Fatalf("ReadHeader() error: %v", err)
	}
	if header.Path != writtenPath || header.FormatVersion != "1.0" {
		t.Fatalf("ReadHeader() = %#v", header)
	}
	assertStringPointer(t, "DSLVersion", header.DSLVersion, "1.2")
	assertStringPointer(t, "DocumentID", header.DocumentID, "Course-A")
	assertStringPointer(t, "SHA256", header.SHA256, "abc123")

	if err := os.WriteFile(writtenPath, []byte(`{"formatVersion":"changed"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if *header.DSLVersion != "1.2" || *header.DocumentID != "Course-A" || *header.SHA256 != "abc123" {
		t.Fatalf("ReadHeader() returned values changed after rewriting the file: %#v", header)
	}
}

// TestReadHeaderNullAndMissingFields проверяет представление отсутствующих и явно пустых полей.
func TestReadHeaderNullAndMissingFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "result.json")
	if err := os.WriteFile(path, []byte(`{
		"formatVersion":"1.0",
		"document":{"dslVersion":null,"sha256":null,"metadata":{"documentId":null}}
	}`), 0o644); err != nil {
		t.Fatal(err)
	}

	header, err := (JSONReader{}).ReadHeader(path)
	if err != nil {
		t.Fatalf("ReadHeader() error: %v", err)
	}
	if header.DSLVersion != nil || header.DocumentID != nil || header.SHA256 != nil {
		t.Fatalf("ReadHeader() nullable fields = %#v, want nil", header)
	}
}

// TestReadHeaderErrors проверяет тип, контекст и исходную причину ошибок чтения и JSON.
func TestReadHeaderErrors(t *testing.T) {
	root := t.TempDir()
	missing := filepath.Join(root, "missing.json")
	broken := filepath.Join(root, "broken.json")
	if err := os.WriteFile(broken, []byte(`{"formatVersion":`), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		// name задаёт имя проверяемого сценария.
		name string
		// path задаёт путь, передаваемый читателю.
		path string
		// wantNotExist отмечает ожидаемую ошибку отсутствующего файла.
		wantNotExist bool
	}{
		{name: "read", path: missing, wantNotExist: true},
		{name: "decode", path: broken},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := (JSONReader{}).ReadHeader(test.path)
			if err == nil {
				t.Fatal("ReadHeader() error = nil")
			}
			var readErr *ReadError
			if !errors.As(err, &readErr) {
				t.Fatalf("ReadHeader() error type = %T, want *ReadError", err)
			}
			if readErr.Path != filepath.Clean(test.path) || readErr.Err == nil {
				t.Fatalf("ReadError = %#v", readErr)
			}
			if err.Error() == readErr.Err.Error() {
				t.Fatalf("ReadError.Error() = %q, want path context", err)
			}
			if !errors.Is(err, readErr.Err) {
				t.Fatal("errors.Is(ReadError, cause) = false")
			}
			if test.wantNotExist && !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("errors.Is(error, os.ErrNotExist) = false: %v", err)
			}
		})
	}
}

// assertStringPointer проверяет обязательное строковое значение указателя.
func assertStringPointer(t *testing.T, field string, got *string, want string) {
	t.Helper()
	if got == nil || *got != want {
		t.Fatalf("%s = %v, want %q", field, got, want)
	}
}
