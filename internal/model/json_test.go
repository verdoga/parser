package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
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

// TestJSONEnumerations проверяет полный закрытый словарь строковых перечислений контракта.
func TestJSONEnumerations(t *testing.T) {
	tests := []struct {
		name string
		got  []string
		want []string
	}{
		{
			name: "severity",
			got:  []string{string(SeverityError), string(SeverityWarning), string(SeverityRecommendation)},
			want: []string{"error", "warning", "recommendation"},
		},
		{
			name: "diagnostic scope",
			got:  []string{string(ScopeElement), string(ScopeLine), string(ScopeBlock), string(ScopeDocument)},
			want: []string{"element", "line", "block", "document"},
		},
		{
			name: "line type",
			got: []string{
				string(LineTag), string(LineHeading), string(LineContent), string(LineBlockStart),
				string(LineBlockEnd), string(LineSeparator), string(LineBlank), string(LineInvalid),
			},
			want: []string{"tag", "heading", "content", "block-start", "block-end", "separator", "blank", "invalid"},
		},
		{
			name: "line ending",
			got:  []string{string(EOLNone), string(EOLLF), string(EOLCRLF)},
			want: []string{"", "\n", "\r\n"},
		},
		{
			name: "element type",
			got: []string{
				string(ElementTag), string(ElementHeadingLevel), string(ElementTitle), string(ElementContent),
				string(ElementIdentifier), string(ElementName), string(ElementVersion), string(ElementNumber),
				string(ElementMediaType), string(ElementSource), string(ElementResourcePath), string(ElementPlaceholder),
				string(ElementUnparsed), string(ElementBlockOpen), string(ElementBlockClose),
			},
			want: []string{
				"tag", "heading-level", "title", "content", "identifier", "name", "version", "number",
				"media-type", "source", "resource-path", "placeholder", "unparsed", "block-open", "block-close",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !reflect.DeepEqual(test.got, test.want) {
				t.Fatalf("enumeration = %q, want %q", test.got, test.want)
			}
		})
	}
}
