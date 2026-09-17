package v1_2

import (
	"reflect"
	"testing"
)

// Компиляционные проверки фиксируют специальные лексические границы DSL v1.2.
var (
	_ func(string) string                = unescapeText
	_ func(string) string                = unescapeMediaSource
	_ func(string, contentPolicy) string = unescape
	_ func(string) resourcePathResult    = scanResourcePaths
	_ func(string) []span                = placeholderRanges
	_ func(string, []span) []span        = unescapedBraceRanges
)

// TestUnescapePolicies проверяет снятие только тех экранов, которые значимы в выбранном контексте.
func TestUnescapePolicies(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		policy contentPolicy
		want   string
	}{
		{name: "plain braces and slash", value: `a\{b\}c\\d\;e\|f\<g\>`, policy: contentPlain, want: `a{b}c\d\;e\|f\<g\>`},
		{name: "wordlist semicolon", value: `one\;two\{x\}`, policy: contentWordlist, want: `one;two{x}`},
		{name: "table pipe", value: `a\|b\{c\}`, policy: contentTable, want: `a|b{c}`},
		{name: "html angles", value: `\<b\>x\</b\>`, policy: contentHTMLText, want: `<b>x</b>`},
		{name: "media quote and slash", value: `a\"b\\c\{d`, policy: contentMediaSource, want: `a"b\c\{d`},
		{name: "resource path is literal", value: `C:\dir\file`, policy: contentResourcePath, want: `C:\dir\file`},
		{name: "editor is literal", value: `\{x\}\|`, policy: contentEditor, want: `\{x\}\|`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := unescape(test.value, test.policy); got != test.want {
				t.Fatalf("unescape(%q, %d) = %q, want %q", test.value, test.policy, got, test.want)
			}
		})
	}
}

// TestSpecializedUnescape проверяет специализированные функции текста и media SOURCE.
func TestSpecializedUnescape(t *testing.T) {
	if got := unescapeText(`\\ \{x\} \; \q`); got != `\ {x} \; \q` {
		t.Fatalf("unescapeText() = %q", got)
	}
	if got := unescapeMediaSource(`say \"yes\" at C:\\tmp\{x`); got != `say "yes" at C:\tmp\{x` {
		t.Fatalf("unescapeMediaSource() = %q", got)
	}
}

// TestPlaceholderRanges проверяет минимальное распознавание защищённых multifill-диапазонов.
func TestPlaceholderRanges(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  []span
	}{
		{name: "none", value: "plain {text}"},
		{name: "one", value: "_____{answer}", want: []span{{0, 13}}},
		{name: "multiple and unicode", value: "я _____{да} + _____{no}", want: []span{{3, 15}, {18, 27}}},
		{name: "requires five underscores", value: "____{no}"},
		{name: "requires closing brace", value: "_____{open"},
		{name: "empty answer", value: "_____{}", want: []span{{0, 7}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := placeholderRanges(test.value); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("placeholderRanges(%q) = %#v, want %#v", test.value, got, test.want)
			}
		})
	}
}

// TestUnescapedBraceRanges проверяет байтовые диапазоны неэкранированных и незащищённых скобок.
func TestUnescapedBraceRanges(t *testing.T) {
	value := `я {a} \{b\} _____{ok} }`
	protected := placeholderRanges(value)
	want := []span{{3, 4}, {5, 6}, {25, 26}}
	if got := unescapedBraceRanges(value, protected); !reflect.DeepEqual(got, want) {
		t.Fatalf("unescapedBraceRanges() = %#v, want %#v", got, want)
	}
}

// TestScanResourcePaths проверяет кавычки, запятые, пробелы и ошибочные остатки списка путей.
func TestScanResourcePaths(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		wantValues    []string
		wantRemainder *span
	}{
		{name: "mixed", value: ` Audio, "../Shared resources", "../Video, additional" `, wantValues: []string{"Audio", "../Shared resources", "../Video, additional"}},
		{name: "windows path literal", value: `C:\audio\unit 1`, wantValues: []string{`C:\audio\unit 1`}},
		{name: "empty", value: "", wantRemainder: spanPointer(span{0, 0})},
		{name: "empty item", value: "one,,two", wantValues: []string{"one"}, wantRemainder: spanPointer(span{4, 5})},
		{name: "trailing comma", value: "one,", wantValues: []string{"one"}, wantRemainder: spanPointer(span{3, 4})},
		{name: "unterminated quote", value: `"one`, wantRemainder: spanPointer(span{0, 4})},
		{name: "junk after quote", value: `"one"junk`, wantValues: []string{"one"}, wantRemainder: spanPointer(span{5, 9})},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := scanResourcePaths(test.value)
			values := make([]string, len(got.paths))
			for i := range got.paths {
				values[i] = got.paths[i].value
			}
			if !reflect.DeepEqual(values, test.wantValues) || !reflect.DeepEqual(got.remainder, test.wantRemainder) {
				t.Fatalf("scanResourcePaths(%q) = values %q, remainder %#v; want %q, %#v", test.value, values, got.remainder, test.wantValues, test.wantRemainder)
			}
		})
	}
}

// spanPointer создаёт указатель на отдельную копию диапазона.
func spanPointer(value span) *span { return &value }
