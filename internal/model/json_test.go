package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestJSONContract проверяет имена полей, значения перечислений и различие null и пустых массивов.
func TestJSONContract(t *testing.T) {
	fixture, err := os.ReadFile(filepath.Join("testdata", "minimal_result.json"))
	if err != nil {
		t.Fatal(err)
	}
	result := Result{
		FormatVersion: FormatVersion,
		Document:      Document{Metadata: Metadata{ResourceDirs: []string{}}, HasErrors: false},
		Processing:    []Processing{}, Lines: []Line{}, Diagnostics: []Diagnostic{},
	}
	got, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got = append(got, '\n')
	if string(got) != string(fixture) {
		t.Fatalf("JSON модели:\n%s\nожидался:\n%s", got, fixture)
	}
}
