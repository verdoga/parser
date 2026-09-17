package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// errInjected обозначает преднамеренный сбой операции в тестовом адаптере.
var errInjected = errors.New("преднамеренная ошибка")

// operationCall описывает один вызов файловой операции.
type operationCall struct {
	// name содержит стабильное тестовое имя операции.
	name string
	// paths содержит пути, переданные операции.
	paths []string
}

// recordingOperations выполняет реальные операции, записывает вызовы и умеет однократно отказать.
type recordingOperations struct {
	// calls сохраняет порядок наблюдавшихся операций.
	calls []operationCall
	// fail задаёт номер вызова с указанным именем, на котором возвращается errInjected.
	fail map[string]int
	// counts считает вызовы каждого имени.
	counts map[string]int
	// shortWrite ограничивает объём одной успешной записи при положительном значении.
	shortWrite int
	// temporaryPaths содержит имена созданных временных файлов.
	temporaryPaths map[string]bool
}

// record записывает вызов и сообщает, должен ли он завершиться преднамеренной ошибкой.
func (o *recordingOperations) record(name string, paths ...string) bool {
	o.calls = append(o.calls, operationCall{name: name, paths: append([]string(nil), paths...)})
	o.counts[name]++
	return o.fail[name] == o.counts[name]
}

// CreateTemp создаёт отслеживаемый временный файл.
func (o *recordingOperations) CreateTemp(directory, pattern string) (syncFile, error) {
	if o.record("create-temp", directory, pattern) {
		return nil, errInjected
	}
	file, err := os.CreateTemp(directory, pattern)
	if err != nil {
		return nil, err
	}
	o.temporaryPaths[file.Name()] = true
	return &recordingFile{File: file, operations: o, kind: "temp"}, nil
}

// Link устанавливает жёсткую ссылку либо возвращает преднамеренную ошибку.
func (o *recordingOperations) Link(oldPath, newPath string) error {
	if o.record("link", oldPath, newPath) {
		return errInjected
	}
	return os.Link(oldPath, newPath)
}

// Rename переименовывает путь либо возвращает преднамеренную ошибку.
func (o *recordingOperations) Rename(oldPath, newPath string) error {
	if o.record("rename", oldPath, newPath) {
		return errInjected
	}
	return os.Rename(oldPath, newPath)
}

// Remove удаляет путь либо возвращает преднамеренную ошибку.
func (o *recordingOperations) Remove(path string) error {
	name := "remove-backup"
	if o.temporaryPaths[path] {
		name = "remove-temp"
	}
	if o.record(name, path) {
		return errInjected
	}
	return os.Remove(path)
}

// Open открывает каталог для последующей синхронизации.
func (o *recordingOperations) Open(path string) (syncFile, error) {
	if o.record("open-directory", path) {
		return nil, errInjected
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &recordingFile{File: file, operations: o, kind: "directory"}, nil
}

// recordingFile добавляет регистрацию операций к настоящему файлу.
type recordingFile struct {
	// File выполняет настоящие файловые операции.
	*os.File
	// operations принимает сведения о вызовах и правила сбоя.
	operations *recordingOperations
	// kind различает временный файл и каталог.
	kind string
}

// Write записывает не больше shortWrite байт за вызов и регистрирует каждый вызов.
func (f *recordingFile) Write(data []byte) (int, error) {
	name := "write-" + f.kind
	if f.operations.record(name, f.Name()) {
		return 0, errInjected
	}
	if limit := f.operations.shortWrite; limit > 0 && len(data) > limit {
		data = data[:limit]
	}
	return f.File.Write(data)
}

// Sync синхронизирует файл либо возвращает преднамеренную ошибку.
func (f *recordingFile) Sync() error {
	name := "sync-" + f.kind
	if f.operations.record(name, f.Name()) {
		return errInjected
	}
	return f.File.Sync()
}

// Close закрывает файл и после этого при необходимости возвращает преднамеренную ошибку.
func (f *recordingFile) Close() error {
	name := "close-" + f.kind
	failed := f.operations.record(name, f.Name())
	err := f.File.Close()
	if failed {
		return errInjected
	}
	return err
}

// TestWriteJSONCreateOnly проверяет атомарное создание и отказ от замены результата.
func TestWriteJSONCreateOnly(t *testing.T) {
	t.Run("создаёт полный синхронизированный результат", func(t *testing.T) {
		directory := t.TempDir()
		target := filepath.Join(directory, "lesson.json")
		data := []byte(`{"document":"данные длиннее одного write"}`)
		original := append([]byte(nil), data...)
		operations := newRecordingOperations()
		operations.shortWrite = 7

		err := fileSystemWithOperations(operations).WriteJSON(target, data, CreateOnly)
		if err != nil {
			t.Fatalf("WriteJSON() error: %v", err)
		}
		assertFileBytes(t, target, original)
		if !reflect.DeepEqual(data, original) {
			t.Fatalf("WriteJSON() changed data: %q", data)
		}
		assertCallsContain(t, operations.calls, "create-temp", "write-temp", "sync-temp", "close-temp", "link", "remove-temp", "open-directory", "sync-directory", "close-directory")
		assertOnlyTargetRemains(t, directory, target)
	})

	t.Run("не изменяет существующий результат", func(t *testing.T) {
		directory := t.TempDir()
		target := filepath.Join(directory, "lesson.json")
		oldData := []byte("old")
		writeTestFile(t, target, oldData)
		operations := newRecordingOperations()

		err := fileSystemWithOperations(operations).WriteJSON(target, []byte("new"), CreateOnly)
		if err == nil {
			t.Fatal("WriteJSON() error = nil")
		}
		if !errors.Is(err, os.ErrExist) {
			t.Fatalf("WriteJSON() error = %v, want os.ErrExist", err)
		}
		assertWriteError(t, err)
		assertFileBytes(t, target, oldData)
		assertOnlyTargetRemains(t, directory, target)
	})
}

// TestWriteJSONReplaceExisting проверяет успешную замену и rollback каждого внешнего этапа.
func TestWriteJSONReplaceExisting(t *testing.T) {
	t.Run("заменяет через резервную копию и удаляет служебные файлы", func(t *testing.T) {
		directory := t.TempDir()
		target := filepath.Join(directory, "lesson.json")
		writeTestFile(t, target, []byte("old"))
		operations := newRecordingOperations()

		if err := fileSystemWithOperations(operations).WriteJSON(target, []byte("new"), ReplaceExisting); err != nil {
			t.Fatalf("WriteJSON() error: %v", err)
		}
		assertFileBytes(t, target, []byte("new"))
		if operations.counts["rename"] < 2 {
			t.Fatalf("rename calls = %d, want at least 2", operations.counts["rename"])
		}
		assertOnlyTargetRemains(t, directory, target)
	})

	failures := []struct {
		name      string
		operation string
		call      int
	}{
		{name: "создание временного файла", operation: "create-temp", call: 1},
		{name: "запись", operation: "write-temp", call: 1},
		{name: "синхронизация временного файла", operation: "sync-temp", call: 1},
		{name: "закрытие временного файла", operation: "close-temp", call: 1},
		{name: "перенос результата в backup", operation: "rename", call: 1},
		{name: "установка нового результата", operation: "rename", call: 2},
		{name: "удаление backup", operation: "remove-backup", call: 1},
		{name: "открытие каталога", operation: "open-directory", call: 1},
		{name: "синхронизация каталога", operation: "sync-directory", call: 1},
		{name: "закрытие каталога", operation: "close-directory", call: 1},
	}
	for _, test := range failures {
		t.Run("rollback при сбое "+test.name, func(t *testing.T) {
			directory := t.TempDir()
			target := filepath.Join(directory, "lesson.json")
			oldData := []byte("old bytes")
			writeTestFile(t, target, oldData)
			operations := newRecordingOperations()
			operations.fail[test.operation] = test.call

			err := fileSystemWithOperations(operations).WriteJSON(target, []byte("new bytes"), ReplaceExisting)
			if err == nil {
				t.Fatal("WriteJSON() error = nil")
			}
			if !errors.Is(err, errInjected) {
				t.Fatalf("WriteJSON() error = %v, want errInjected", err)
			}
			assertWriteError(t, err)
			assertFileBytes(t, target, oldData)
			assertOnlyTargetRemains(t, directory, target)
		})
	}
}

// TestWriteJSONValidation проверяет ошибку режима и недопустимого целевого пути.
func TestWriteJSONValidation(t *testing.T) {
	tests := []struct {
		name   string
		target string
		mode   WriteMode
	}{
		{name: "неизвестный режим", target: filepath.Join(t.TempDir(), "result.json"), mode: WriteMode(99)},
		{name: "пустой путь", target: "", mode: CreateOnly},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := (FileSystem{}).WriteJSON(test.target, []byte("data"), test.mode)
			if err == nil {
				t.Fatal("WriteJSON() error = nil")
			}
			assertWriteError(t, err)
		})
	}
}

// TestWriteError проверяет поля, текст и стандартную цепочку причин WriteError.
func TestWriteError(t *testing.T) {
	cause := io.ErrClosedPipe
	err := &WriteError{Operation: "sync temporary", Path: "/tmp/result", Err: cause}
	if err.Operation != "sync temporary" || err.Path != "/tmp/result" || err.Err != cause {
		t.Fatalf("WriteError fields = %#v", err)
	}
	text := err.Error()
	for _, fragment := range []string{"sync temporary", "/tmp/result", cause.Error()} {
		if !strings.Contains(text, fragment) {
			t.Errorf("WriteError.Error() = %q, missing %q", text, fragment)
		}
	}
	if !errors.Is(err, cause) {
		t.Fatalf("errors.Is(%v, cause) = false", err)
	}
	var got *WriteError
	if !errors.As(err, &got) || got != err {
		t.Fatalf("errors.As() = %#v, want original error", got)
	}
}

// newRecordingOperations создаёт готовый регистрирующий адаптер.
func newRecordingOperations() *recordingOperations {
	return &recordingOperations{
		fail:           make(map[string]int),
		counts:         make(map[string]int),
		temporaryPaths: make(map[string]bool),
	}
}

// assertFileBytes проверяет точное содержимое файла.
func assertFileBytes(t *testing.T, path string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", path, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("file %q = %q, want %q", path, got, want)
	}
}

// assertOnlyTargetRemains проверяет отсутствие временных файлов и backup.
func assertOnlyTargetRemains(t *testing.T, directory, target string) {
	t.Helper()
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || filepath.Join(directory, entries[0].Name()) != target {
		names := make([]string, 0, len(entries))
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		t.Fatalf("directory entries = %v, want only %q", names, filepath.Base(target))
	}
}

// assertCallsContain проверяет наличие обязательных этапов независимо от дополнительных rollback-вызовов.
func assertCallsContain(t *testing.T, calls []operationCall, names ...string) {
	t.Helper()
	seen := make(map[string]bool)
	for _, call := range calls {
		seen[call.name] = true
	}
	for _, name := range names {
		if !seen[name] {
			t.Errorf("operation %q was not called; calls: %s", name, formatCalls(calls))
		}
	}
}

// formatCalls форматирует вызовы для сообщения об ошибке теста.
func formatCalls(calls []operationCall) string {
	parts := make([]string, len(calls))
	for index, call := range calls {
		parts[index] = fmt.Sprintf("%s(%s)", call.name, strings.Join(call.paths, ", "))
	}
	return strings.Join(parts, "; ")
}

// assertWriteError проверяет типизированную оболочку ошибки записи.
func assertWriteError(t *testing.T, err error) *WriteError {
	t.Helper()
	var writeError *WriteError
	if !errors.As(err, &writeError) {
		t.Fatalf("error type = %T, want *WriteError: %v", err, err)
	}
	if writeError.Operation == "" || writeError.Path == "" || writeError.Err == nil {
		t.Fatalf("incomplete WriteError = %#v", writeError)
	}
	return writeError
}
