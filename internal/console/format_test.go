package console

import (
	"bytes"
	"errors"
	"testing"
)

// TestFormatFile проверяет статусы, все поля, экранирование строк и отсутствие завершающего LF.
func TestFormatFile(t *testing.T) {
	tests := []struct {
		// name задаёт имя сценария форматирования.
		name string
		// line содержит входные данные файловой строки.
		line FileLine
		// want содержит ожидаемую строку целиком.
		want string
	}{
		{
			name: "успех без сообщения",
			line: FileLine{Status: Success, Action: "created", Path: "/a b.txt"},
			want: `УСПЕХ action=created errors=0 path="/a b.txt"`,
		},
		{
			name: "успех с диагностиками и сообщением",
			line: FileLine{Status: Success, Action: "replaced", Errors: 2, Path: "/урок.txt", Message: "есть замечания"},
			want: `УСПЕХ action=replaced errors=2 path="/урок.txt" message="есть замечания"`,
		},
		{
			name: "ошибка",
			line: FileLine{Status: Failed, Action: "failed", Errors: 1, Path: "x.txt", Message: "нет доступа"},
			want: `ОШИБКА action=failed errors=1 path="x.txt" message="нет доступа"`,
		},
		{
			name: "экранирование пути и сообщения",
			line: FileLine{Status: Failed, Action: "failed", Path: "C:\\data\\\"lesson\".txt", Message: "первая\nвторая\tстрока"},
			want: `ОШИБКА action=failed errors=0 path="C:\\data\\\"lesson\".txt" message="первая\nвторая\tстрока"`,
		},
		{
			name: "пустые значения",
			line: FileLine{Status: Success},
			want: `УСПЕХ action= errors=0 path=""`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := FormatFile(test.line); got != test.want {
				t.Fatalf("FormatFile() = %q, want %q", got, test.want)
			}
		})
	}
}

// TestFormatSummary проверяет нормативный порядок всех девяти счётчиков и отсутствие завершающего LF.
func TestFormatSummary(t *testing.T) {
	tests := []struct {
		// name задаёт имя сценария форматирования.
		name string
		// line содержит входные данные итоговой строки.
		line SummaryLine
		// want содержит ожидаемую строку целиком.
		want string
	}{
		{
			name: "все счётчики",
			line: SummaryLine{Found: 1, Parsed: 2, Created: 3, Replaced: 4, Skipped: 5, Success: 6, Failed: 7, Diagnostics: 8, ScanErrors: 9},
			want: "ИТОГ found=1 parsed=2 created=3 replaced=4 skipped=5 success=6 failed=7 diagnostics=8 scanErrors=9",
		},
		{
			name: "нулевые счётчики",
			line: SummaryLine{},
			want: "ИТОГ found=0 parsed=0 created=0 replaced=0 skipped=0 success=0 failed=0 diagnostics=0 scanErrors=0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := FormatSummary(test.line); got != test.want {
				t.Fatalf("FormatSummary() = %q, want %q", got, test.want)
			}
		})
	}
}

// TestWriterRoutesCompleteLines проверяет назначение потоков, порядок строк и ровно один LF на строку.
func TestWriterRoutesCompleteLines(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	writer := NewWriter(&stdout, &stderr)
	file := FileLine{Status: Success, Action: "skipped", Path: "a.txt"}
	summary := SummaryLine{Found: 1, Skipped: 1, Success: 1}

	if err := writer.File(file); err != nil {
		t.Fatalf("File() error: %v", err)
	}
	if err := writer.Summary(summary); err != nil {
		t.Fatalf("Summary() error: %v", err)
	}
	if err := writer.OperationalError(errors.New("сбой запуска")); err != nil {
		t.Fatalf("OperationalError() error: %v", err)
	}

	wantStdout := FormatFile(file) + "\n" + FormatSummary(summary) + "\n"
	if got := stdout.String(); got != wantStdout {
		t.Fatalf("stdout = %q, want %q", got, wantStdout)
	}
	if got, want := stderr.String(), "сбой запуска\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}

// TestWriterReturnsDestinationErrors проверяет передачу ошибок каждого целевого потока вызывающему коду.
func TestWriterReturnsDestinationErrors(t *testing.T) {
	errStdout := errors.New("stdout write")
	errStderr := errors.New("stderr write")
	writer := NewWriter(failingWriter{err: errStdout}, failingWriter{err: errStderr})
	tests := []struct {
		// name задаёт имя проверяемой операции.
		name string
		// write запускает одну операцию Writer.
		write func() error
		// want содержит ошибку соответствующего потока.
		want error
	}{
		{name: "файловая строка", write: func() error { return writer.File(FileLine{Status: Success}) }, want: errStdout},
		{name: "итоговая строка", write: func() error { return writer.Summary(SummaryLine{}) }, want: errStdout},
		{name: "операционная ошибка", write: func() error { return writer.OperationalError(errors.New("operation")) }, want: errStderr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.write(); !errors.Is(err, test.want) {
				t.Fatalf("error = %v, want an error wrapping %v", err, test.want)
			}
		})
	}
}

// failingWriter всегда возвращает настроенную ошибку.
type failingWriter struct {
	// err задаёт ошибку каждой попытки записи.
	err error
}

// Write реализует io.Writer для проверки передачи ошибки.
func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }
