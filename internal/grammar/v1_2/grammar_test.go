package v1_2

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
	"reflect"
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

// TestGrammarIdentity проверяет версию и независимост экземпляров грамматики.
func TestGrammarIdentity(t *testing.T) {
	first := New()
	second := New()
	if first == nil || second == nil {
		t.Fatal("New() returned nil")
	}
	if got := first.Version(); got != "1.2" {
		t.Fatalf("Version() = %q, want 1.2", got)
	}
	if got := second.Version(); got != "1.2" {
		t.Fatalf("second Version() = %q, want 1.2", got)
	}
}

// TestClassifyRootSyntax проверяет все виды строк корневого контекста.
func TestClassifyRootSyntax(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		lineType   model.LineType
		elements   []elementWant
		problem    grammar.ProblemKind
		block      grammar.BlockAction
		transition grammar.Transition
	}{
		{name: "blank", raw: " \t", lineType: model.LineBlank},
		{name: "content", raw: "ordinary \\{text\\}", lineType: model.LineContent, elements: []elementWant{{model.ElementContent, "ordinary {text}", 0, 17}}},
		{name: "heading", raw: "### Lesson", lineType: model.LineHeading, elements: []elementWant{{model.ElementHeadingLevel, "3", 0, 3}, {model.ElementTitle, "Lesson", 4, 10}}, transition: grammar.Transition{SetStep: true}},
		{name: "block close", raw: "}", lineType: model.LineBlockEnd, elements: []elementWant{{model.ElementBlockClose, "", 0, 1}}, block: grammar.BlockClose},
		{name: "malformed close", raw: "} tail", lineType: model.LineInvalid, problem: grammar.ProblemMalformedBlockClose},
		{name: "unknown tag", raw: "@missing value", lineType: model.LineInvalid, problem: grammar.ProblemUnknownTag},
		{name: "block start", raw: "@text {", lineType: model.LineBlockStart, block: grammar.BlockOpen},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decision := New().Classify(testRequest(test.raw, contextRoot, grammar.ContentDSL))
			if decision.LineType != test.lineType || decision.Block.Action != test.block {
				t.Fatalf("Classify() = type %q, block %v; want %q, %v", decision.LineType, decision.Block.Action, test.lineType, test.block)
			}
			if !reflect.DeepEqual(decision.Transition, test.transition) {
				t.Fatalf("Transition = %#v, want %#v", decision.Transition, test.transition)
			}
			assertElements(t, decision.Elements, test.elements)
			if test.problem != "" {
				assertProblemKinds(t, decision.Problems, test.problem)
			}
		})
	}
}

// TestClassifyContentContexts проверяет контекстную классификацию непрозрачного содержимого.
func TestClassifyContentContexts(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		context  grammar.Context
		mode     grammar.ContentMode
		lineType model.LineType
		value    string
		problem  grammar.ProblemKind
	}{
		{name: "opaque tag-like text", raw: "@task not-a-tag", context: contextOpaque, mode: grammar.ContentOpaque, lineType: model.LineContent, value: "@task not-a-tag"},
		{name: "example separator", raw: "---", context: contextExample, mode: grammar.ContentOpaque, lineType: model.LineSeparator},
		{name: "table separator remains content", raw: "---", context: contextTable, mode: grammar.ContentOpaque, lineType: model.LineContent, value: "---"},
		{name: "text blank", raw: "", context: contextText, mode: grammar.ContentOpaque, lineType: model.LineBlank},
		{name: "table blank", raw: "", context: contextTable, mode: grammar.ContentOpaque, lineType: model.LineBlank},
		{name: "multifill blank", raw: "", context: contextMultifill, mode: grammar.ContentMultifill, lineType: model.LineContent, value: ""},
		{name: "editor braces are literal", raw: "{literal}", context: contextEditor, mode: grammar.ContentEditor, lineType: model.LineContent, value: "{literal}"},
		{name: "multifill placeholder protects braces", raw: "Use _____{answer}", context: contextMultifill, mode: grammar.ContentMultifill, lineType: model.LineContent, value: "Use _____{answer}"},
		{name: "unprotected multifill brace", raw: "Use {answer}", context: contextMultifill, mode: grammar.ContentMultifill, lineType: model.LineContent, value: "Use {answer}", problem: grammar.ProblemUnescapedBrace},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decision := New().Classify(testRequest(test.raw, test.context, test.mode))
			if decision.LineType != test.lineType {
				t.Fatalf("LineType = %q, want %q", decision.LineType, test.lineType)
			}
			if test.lineType == model.LineContent {
				assertElements(t, decision.Elements, []elementWant{{model.ElementContent, test.value, 0, len(test.raw)}})
			}
			if test.problem != "" {
				assertProblemKinds(t, decision.Problems, test.problem)
			}
		})
	}
}

// elementWant задаёт ожидаемый элемент решения.
type elementWant struct {
	// elementType содержит тип элемента.
	elementType model.ElementType
	// value содержит ожидаемое семантическое значение.
	value string
	// start содержит включительную байтовую границу.
	start int
	// end содержит исключающую байтовую границу.
	end int
}

// testRequest создаёт запрос с нормативной обрезкой ASCII-пробелов.
func testRequest(raw string, context grammar.Context, mode grammar.ContentMode) grammar.GrammarRequest {
	trimmed := raw
	for len(trimmed) > 0 && (trimmed[0] == ' ' || trimmed[0] == '\t') {
		trimmed = trimmed[1:]
	}
	for len(trimmed) > 0 && (trimmed[len(trimmed)-1] == ' ' || trimmed[len(trimmed)-1] == '\t') {
		trimmed = trimmed[:len(trimmed)-1]
	}
	return grammar.GrammarRequest{Raw: raw, Trimmed: trimmed, Line: 7, Context: context, ContentMode: mode}
}

// assertElements сравнивает типы, значения и байтовые границы элементов.
func assertElements(t *testing.T, got []grammar.ElementDecision, want []elementWant) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("elements count = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i].Type != want[i].elementType || got[i].Value == nil || *got[i].Value != want[i].value || got[i].ByteStart != want[i].start || got[i].ByteEnd != want[i].end {
			t.Errorf("element[%d] = %#v, want type %q, value %q, range [%d,%d)", i, got[i], want[i].elementType, want[i].value, want[i].start, want[i].end)
		}
	}
}

// assertProblemKinds сравнивает порядок кодов грамматических проблем.
func assertProblemKinds(t *testing.T, got []grammar.Problem, want ...grammar.ProblemKind) {
	t.Helper()
	kinds := make([]grammar.ProblemKind, len(got))
	for i := range got {
		kinds[i] = got[i].Kind
	}
	if !reflect.DeepEqual(kinds, want) {
		t.Fatalf("problem kinds = %v, want %v", kinds, want)
	}
}
