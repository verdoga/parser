package console

import (
	"bytes"
	"errors"
	"testing"
)

// TestFormatFile проверяет все отображаемые поля и отсутствие перевода строки.
func TestFormatFile(t *testing.T) {
	tests := []struct {
		name string
		line FileLine
		want string
	}{
		{name: "success", line: FileLine{Status: Success, Action: "created", Path: "/a b.txt"}, want: `УСПЕХ action=created errors=0 path="/a b.txt"`},
		{name: "failure message", line: FileLine{Status: Failed, Action: "failed", Errors: 2, Path: "x.txt", Message: "нет доступа"}, want: `ОШИБКА action=failed errors=2 path="x.txt" message="нет доступа"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := FormatFile(test.line); got != test.want {
				t.Fatalf("FormatFile() = %q, want %q", got, test.want)
			}
		})
	}
}

// TestFormatSummary проверяет нормативный порядок девяти счётчиков.
func TestFormatSummary(t *testing.T) {
	line := SummaryLine{Found: 1, Parsed: 2, Created: 3, Replaced: 4, Skipped: 5, Success: 6, Failed: 7, Diagnostics: 8, ScanErrors: 9}
	want := "ИТОГ found=1 parsed=2 created=3 replaced=4 skipped=5 success=6 failed=7 diagnostics=8 scanErrors=9"
	if got := FormatSummary(line); got != want {
		t.Fatalf("FormatSummary() = %q, want %q", got, want)
	}
}

// TestWriterRouting проверяет назначение потоков и единственный LF.
func TestWriterRouting(t *testing.T) {
	var stdout, stderr bytes.Buffer
	writer := NewWriter(&stdout, &stderr)
	if err := writer.File(FileLine{Status: Success, Action: "skipped", Path: "a.txt"}); err != nil {
		t.Fatal(err)
	}
	if err := writer.Summary(SummaryLine{Found: 1}); err != nil {
		t.Fatal(err)
	}
	if err := writer.OperationalError(errors.New("сбой")); err != nil {
		t.Fatal(err)
	}
	wantOut := FormatFile(FileLine{Status: Success, Action: "skipped", Path: "a.txt"}) + "\n" + FormatSummary(SummaryLine{Found: 1}) + "\n"
	if stdout.String() != wantOut || stderr.String() != "сбой\n" {
		t.Fatalf("stdout = %q, stderr = %q", stdout.String(), stderr.String())
	}
}

// TestWriterReturnsUnderlyingError проверяет ошибку целевого writer.
func TestWriterReturnsUnderlyingError(t *testing.T) {
	want := errors.New("write")
	writer := NewWriter(failingWriter{err: want}, failingWriter{err: want})
	if err := writer.File(FileLine{}); !errors.Is(err, want) {
		t.Fatalf("File() error = %v", err)
	}
	if err := writer.OperationalError(errors.New("operation")); !errors.Is(err, want) {
		t.Fatalf("OperationalError() error = %v", err)
	}
}

// failingWriter всегда возвращает настроенную ошибку.
type failingWriter struct{ err error }

// Write реализует io.Writer для проверки передачи ошибки.
func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }
