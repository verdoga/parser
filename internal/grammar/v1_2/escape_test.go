package v1_2

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
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

// TestEscaping проверяет контекстные escape-последовательности и отсутствие повторного разбора.
func TestEscaping(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		policy contentPolicy
		want   string
	}{
		{name: "plain structural characters", input: `\@tag \#head \{x\} \\`, policy: contentPlain, want: `@tag #head {x} \`},
		{name: "unknown escape remains", input: `a\qb`, policy: contentPlain, want: `a\qb`},
		{name: "single pass", input: `\\\{`, policy: contentPlain, want: `\{`},
		{name: "wordlist semicolon", input: `one\;two`, policy: contentWordlist, want: `one;two`},
		{name: "example separator", input: `\---`, policy: contentExample, want: `---`},
		{name: "table cell", input: `a\|b`, policy: contentTable, want: `a|b`},
		{name: "html text", input: `\<b\>`, policy: contentHTMLText, want: `<b>`},
		{name: "editor tag start", input: `\@task`, policy: contentEditor, want: `@task`},
		{name: "resource path keeps slash", input: `dir\name`, policy: contentResourcePath, want: `dir\name`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := unescape(test.input, test.policy); got != test.want {
				t.Fatalf("unescape(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
	if got := unescapeText(`\@x\{y\}`); got != `@x{y}` {
		t.Fatalf("unescapeText() = %q", got)
	}
	if got := unescapeMediaSource(`a\"b\\c\q`); got != `a"b\c\q` {
		t.Fatalf("unescapeMediaSource() = %q", got)
	}
}

// TestUnescapedBraceRanges проверяет экранированные и защищённые байтовые диапазоны скобок.
func TestUnescapedBraceRanges(t *testing.T) {
	got := unescapedBraceRanges(`я{a\}b}_____{x}z{`, []span{{start: 8, end: 16}})
	want := []span{{start: 2, end: 3}, {start: 7, end: 8}, {start: 17, end: 18}}
	if len(got) != len(want) {
		t.Fatalf("ranges = %#v, want %#v", got, want)
	}
	for index := range want {
		requireSpan(t, got[index], want[index].start, want[index].end)
	}
}

// TestMultifillPlaceholders проверяет только минимально защищённую специальную форму.
func TestMultifillPlaceholders(t *testing.T) {
	tests := []struct {
		name, input string
		want        []span
	}{
		{name: "one", input: "x_____{yes}y", want: []span{{start: 1, end: 11}}},
		{name: "two", input: "_____{a} _____{б}", want: []span{{start: 0, end: 8}, {start: 9, end: 18}}},
		{name: "ordinary placeholder", input: "_____"},
		{name: "too few underscores", input: "____{x}"},
		{name: "escaped opening brace", input: `_____\{x}`},
		{name: "missing close", input: "_____{x"},
		{name: "empty answer", input: "_____{}"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := placeholderRanges(test.input)
			if len(got) != len(test.want) {
				t.Fatalf("placeholderRanges(%q) = %#v, want %#v", test.input, got, test.want)
			}
			for index := range test.want {
				requireSpan(t, got[index], test.want[index].start, test.want[index].end)
			}
		})
	}
	decision := New().Classify(request("text _____{answer}", contextMultifill, grammar.ContentMultifill))
	for _, element := range decision.Elements {
		if element.Type == model.ElementPlaceholder {
			t.Fatal("multifill classification created forbidden placeholder element")
		}
	}
	if answer := New().Classify(request("@answer wrong", contextMultifill, grammar.ContentMultifill)); len(answer.Problems) == 0 {
		t.Fatal("@answer in multifill has no problem")
	}
}
