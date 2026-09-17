package cache

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

// stubReader возвращает заранее заданные заголовки и ошибки и сохраняет порядок обращений.
type stubReader struct {
	// headers сопоставляет пути с возвращаемыми заголовками.
	headers map[string]Header
	// errors сопоставляет пути с возвращаемыми ошибками.
	errors map[string]error
	// calls хранит пути в порядке вызовов ReadHeader.
	calls []string
}

// ReadHeader сохраняет путь и возвращает настроенный для него результат.
func (r *stubReader) ReadHeader(path string) (Header, error) {
	r.calls = append(r.calls, path)
	return r.headers[path], r.errors[path]
}

// TestBuildIndex проверяет нормализацию, сортировку, сохранение частных ошибок и индексирование ID.
func TestBuildIndex(t *testing.T) {
	root := t.TempDir()
	a := filepath.Join(root, "a.json")
	b := filepath.Join(root, "b.json")
	c := filepath.Join(root, "c.json")
	readFailure := errors.New("cannot decode")
	reader := &stubReader{
		headers: map[string]Header{
			a: header(a, "1.0", "Course", "1.2", "sum-a"),
			b: header(b, "1.0", "course", "1.2", "sum-b"),
		},
		errors: map[string]error{c: readFailure},
	}
	paths := []string{c, filepath.Join(root, "nested", "..", "b.json"), a}
	original := append([]string(nil), paths...)

	index, err := Build(paths, reader)
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if !reflect.DeepEqual(paths, original) {
		t.Fatalf("Build() mutated paths: got %q, want %q", paths, original)
	}
	if !reflect.DeepEqual(reader.calls, []string{a, b, c}) {
		t.Fatalf("ReadHeader calls = %q, want sorted normalized paths", reader.calls)
	}
	entries := index.Entries()
	if len(entries) != 3 || entries[0].Path != a || entries[1].Path != b || entries[2].Path != c {
		t.Fatalf("Entries() = %#v", entries)
	}
	if !errors.Is(entries[2].Err, readFailure) {
		t.Fatalf("entry error = %v, want %v", entries[2].Err, readFailure)
	}

	match := index.MatchSource(SourceIdentity{TargetPath: filepath.Join(root, "missing.json"), DocumentID: "COURSE", DSLVersion: "1.2", SHA256: "other"})
	if match.Kind != MatchDuplicateDocumentID || len(match.Conflicts) != 2 || match.Conflicts[0].Path != a || match.Conflicts[1].Path != b {
		t.Fatalf("case-insensitive indexed match = %#v", match)
	}
}

// TestBuildIndexRejectsInvalidReader проверяет ошибку недопустимой конфигурации.
func TestBuildIndexRejectsInvalidReader(t *testing.T) {
	if _, err := Build([]string{"result.json"}, nil); err == nil {
		t.Fatal("Build(paths, nil) error = nil")
	}
}

// TestEntriesReturnsIndependentCopies проверяет отсутствие доступа к внутреннему срезу записей.
func TestEntriesReturnsIndependentCopies(t *testing.T) {
	path := filepath.Join(t.TempDir(), "result.json")
	reader := &stubReader{headers: map[string]Header{path: header(path, "1.0", "id", "1.2", "sum")}}
	index, err := Build([]string{path}, reader)
	if err != nil {
		t.Fatal(err)
	}

	first := index.Entries()
	first[0].Path = "changed"
	*first[0].Header.DocumentID = "changed"
	second := index.Entries()
	if second[0].Path != path || second[0].Header.DocumentID == nil || *second[0].Header.DocumentID != "id" {
		t.Fatalf("Entries() aliases returned data: %#v", second)
	}
}

// TestMatchSource проверяет все классификации и точные правила актуальности.
func TestMatchSource(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.json")
	other := filepath.Join(root, "other.json")
	readFailure := errors.New("broken JSON")
	tests := []struct {
		// name задаёт имя сценария сопоставления.
		name string
		// entries задаёт содержимое строящегося индекса.
		entries []Entry
		// source задаёт проверяемую идентичность источника.
		source SourceIdentity
		// wantKind задаёт ожидаемую классификацию.
		wantKind MatchKind
		// wantTarget отмечает ожидаемое наличие целевой записи.
		wantTarget bool
		// wantConflicts задаёт ожидаемые пути конфликтов.
		wantConflicts []string
	}{
		{name: "missing", source: identity(target, "id", "1.2", "sum"), wantKind: MatchMissing},
		{name: "current ignores ID case", entries: []Entry{{Path: target, Header: header(target, "1.0", "ID", "1.2", "sum")}}, source: identity(target, "id", "1.2", "sum"), wantKind: MatchCurrent, wantTarget: true},
		{name: "empty source ID", entries: []Entry{{Path: target, Header: header(target, "1.0", "id", "1.2", "sum")}}, source: identity(target, "", "1.2", "sum"), wantKind: MatchTargetConflict, wantTarget: true},
		{name: "target read error", entries: []Entry{{Path: target, Err: readFailure}}, source: identity(target, "id", "1.2", "sum"), wantKind: MatchTargetConflict, wantTarget: true},
		{name: "unsupported format", entries: []Entry{{Path: target, Header: header(target, "2.0", "id", "1.2", "sum")}}, source: identity(target, "id", "1.2", "sum"), wantKind: MatchTargetConflict, wantTarget: true},
		{name: "different ID", entries: []Entry{{Path: target, Header: header(target, "1.0", "other", "1.2", "sum")}}, source: identity(target, "id", "1.2", "sum"), wantKind: MatchTargetConflict, wantTarget: true},
		{name: "different DSL version", entries: []Entry{{Path: target, Header: header(target, "1.0", "id", "1.1", "sum")}}, source: identity(target, "id", "1.2", "sum"), wantKind: MatchTargetConflict, wantTarget: true},
		{name: "different SHA-256", entries: []Entry{{Path: target, Header: header(target, "1.0", "id", "1.2", "other")}}, source: identity(target, "id", "1.2", "sum"), wantKind: MatchTargetConflict, wantTarget: true},
		{name: "missing nullable field", entries: []Entry{{Path: target, Header: Header{Path: target, FormatVersion: "1.0"}}}, source: identity(target, "id", "1.2", "sum"), wantKind: MatchTargetConflict, wantTarget: true},
		{name: "foreign duplicate", entries: []Entry{{Path: other, Header: header(other, "1.0", "ID", "1.2", "other")}}, source: identity(target, "id", "1.2", "sum"), wantKind: MatchDuplicateDocumentID, wantConflicts: []string{other}},
		{name: "duplicate overrides current", entries: []Entry{{Path: other, Header: header(other, "1.0", "id", "1.2", "other")}, {Path: target, Header: header(target, "1.0", "ID", "1.2", "sum")}}, source: identity(target, "id", "1.2", "sum"), wantKind: MatchDuplicateDocumentID, wantTarget: true, wantConflicts: []string{other}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			index := buildIndexFromEntries(t, test.entries)
			match := index.MatchSource(test.source)
			if match.Kind != test.wantKind || (match.Target != nil) != test.wantTarget {
				t.Fatalf("MatchSource() = %#v, want kind %v, target %v", match, test.wantKind, test.wantTarget)
			}
			gotConflicts := make([]string, len(match.Conflicts))
			for position, conflict := range match.Conflicts {
				gotConflicts[position] = conflict.Path
			}
			if !reflect.DeepEqual(gotConflicts, test.wantConflicts) {
				t.Fatalf("MatchSource() conflicts = %q, want %q", gotConflicts, test.wantConflicts)
			}
			if match.Target != nil {
				match.Target.Path = "changed"
				if match.Target.Header.DocumentID != nil {
					*match.Target.Header.DocumentID = "changed"
				}
			}
			if len(match.Conflicts) != 0 {
				match.Conflicts[0].Path = "changed"
				if match.Conflicts[0].Header.DocumentID != nil {
					*match.Conflicts[0].Header.DocumentID = "changed"
				}
			}
			for _, entry := range index.Entries() {
				if entry.Path == "changed" || entry.Header.DocumentID != nil && *entry.Header.DocumentID == "changed" {
					t.Fatal("MatchSource() exposes mutable index entries")
				}
			}
		})
	}
}

// TestMatchCanSkip проверяет разрешение пропуска только для актуального результата.
func TestMatchCanSkip(t *testing.T) {
	for _, kind := range []MatchKind{MatchMissing, MatchCurrent, MatchTargetConflict, MatchDuplicateDocumentID, MatchKind(99)} {
		if got, want := (Match{Kind: kind}).CanSkip(), kind == MatchCurrent; got != want {
			t.Errorf("Match{Kind: %v}.CanSkip() = %v, want %v", kind, got, want)
		}
	}
}

// buildIndexFromEntries строит индекс через публичную границу чтения.
func buildIndexFromEntries(t *testing.T, entries []Entry) Index {
	t.Helper()
	reader := &stubReader{headers: make(map[string]Header), errors: make(map[string]error)}
	paths := make([]string, len(entries))
	for position, entry := range entries {
		paths[position] = entry.Path
		reader.headers[entry.Path] = entry.Header
		reader.errors[entry.Path] = entry.Err
	}
	index, err := Build(paths, reader)
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	return index
}

// header создаёт полный заголовок для теста сопоставления.
func header(path, formatVersion, documentID, dslVersion, sha256 string) Header {
	return Header{Path: path, FormatVersion: formatVersion, DocumentID: stringPointer(documentID), DSLVersion: stringPointer(dslVersion), SHA256: stringPointer(sha256)}
}

// identity создаёт сведения источника для теста сопоставления.
func identity(targetPath, documentID, dslVersion, sha256 string) SourceIdentity {
	return SourceIdentity{TargetPath: targetPath, DocumentID: documentID, DSLVersion: dslVersion, SHA256: sha256}
}

// stringPointer возвращает указатель на отдельную копию строки.
func stringPointer(value string) *string { return &value }
