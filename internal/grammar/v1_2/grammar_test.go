package v1_2

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
	"testing"
)

// Компиляционные проверки фиксируют внешнюю границу грамматики DSL v1.2.
var (
	_ grammar.Grammar                                                = implementation{}
	_ func() grammar.Grammar                                         = New
	_ func(string) (span, tagName, bool)                             = scanTag
	_ func(string) formResult                                        = detectForm
	_ func(string) resourcePathResult                                = scanResourcePaths
	_ func(string) []span                                            = placeholderRanges
	_ func(string, []span) []span                                    = unescapedBraceRanges
	_ func(string) string                                            = unescapeText
	_ func(string) string                                            = unescapeMediaSource
	_ func(string, contentPolicy) string                             = unescape
	_ func([]grammar.Problem, formResult) grammar.RecoveryDecision   = recoveryFor
	_ func(grammar.GrammarRequest) grammar.GrammarDecision           = classifyBlank
	_ func(string) heading                                           = scanHeading
	_ model.LineType                                                 = model.LineBlank
	_ func(model.ElementType, span, *string) grammar.ElementDecision = elementDecision
	_ func(model.ElementType, int) grammar.ElementDecision           = missingElementDecision
	_ func(span) grammar.ElementDecision                             = unparsedElementDecision
	_ func(span, string) grammar.ElementDecision                     = contentRemainderDecision
)

// TestGrammarIdentity проверяет создание независимых реализаций и точную версию.
func TestGrammarIdentity(t *testing.T) {
	first := New()
	second := New()
	if first == nil || second == nil {
		t.Fatal("New() returned nil")
	}
	if first.Version() != "1.2" || second.Version() != "1.2" {
		t.Fatalf("versions = %q, %q; want 1.2", first.Version(), second.Version())
	}
}

// TestDeclarationScanning проверяет имя тега, регистр, разделитель и определение формы.
func TestDeclarationScanning(t *testing.T) {
	tagTests := []struct {
		name  string
		input string
		want  tagName
		end   int
		ok    bool
	}{
		{name: "plain", input: "@task id", want: tagTask, end: 5, ok: true},
		{name: "hyphen", input: "@dsl-version 1.2", want: tagDSLVersion, end: 12, ok: true},
		{name: "uppercase is scanned", input: "@Task id", want: "Task", end: 5, ok: true},
		{name: "at only", input: "@", ok: false},
		{name: "not at start", input: "x@task", ok: false},
		{name: "invalid name", input: "@ task", ok: false},
	}
	for _, test := range tagTests {
		t.Run(test.name, func(t *testing.T) {
			location, got, ok := scanTag(test.input)
			if ok != test.ok || got != test.want {
				t.Fatalf("scanTag(%q) = (%#v, %q, %t)", test.input, location, got, ok)
			}
			if ok {
				requireSpan(t, location, 0, test.end)
			}
		})
	}

	formTests := []struct {
		name      string
		input     string
		form      grammar.TagForm
		open      *span
		ambiguous bool
	}{
		{name: "line", input: "@note text", form: grammar.TagFormLine},
		{name: "block", input: "@note {", form: grammar.TagFormBlock, open: &span{start: 6, end: 7}},
		{name: "escaped brace", input: `@note \{`, form: grammar.TagFormLine},
		{name: "brace followed by content", input: "@note { tail", ambiguous: true},
		{name: "multiple braces", input: "@note { {", ambiguous: true},
	}
	for _, test := range formTests {
		t.Run(test.name, func(t *testing.T) {
			got := detectForm(test.input)
			if got.form != test.form || got.ambiguous != test.ambiguous || (got.blockOpen == nil) != (test.open == nil) {
				t.Fatalf("detectForm(%q) = %#v", test.input, got)
			}
			if test.open != nil {
				requireSpan(t, *got.blockOpen, test.open.start, test.open.end)
			}
		})
	}
}

// TestDecisionConstructors проверяет диапазоны и значения всех видов элементов.
func TestDecisionConstructors(t *testing.T) {
	text := "значение"
	requireElement(t, elementDecision(model.ElementContent, span{2, 6}, &text), model.ElementContent, &text, 2, 6)
	requireElement(t, missingElementDecision(model.ElementIdentifier, 7), model.ElementIdentifier, nil, 7, 7)
	requireElement(t, unparsedElementDecision(span{3, 9}), model.ElementUnparsed, nil, 3, 9)
	requireElement(t, contentRemainderDecision(span{4, 10}, text), model.ElementContent, &text, 4, 10)
}
