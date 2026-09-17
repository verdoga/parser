package v1_2

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
	"testing"
)

// Компиляционные проверки фиксируют конструкторы локальных проблем DSL v1.2.
var (
	_ func(span) grammar.Problem             = unknownTagProblem
	_ func(span) grammar.Problem             = missingSeparatorProblem
	_ func(span, int) grammar.Problem        = unsupportedFormProblem
	_ func(span, int) grammar.Problem        = missingArgumentProblem
	_ func(span, int) grammar.Problem        = extraContentProblem
	_ func(span) grammar.Problem             = malformedBlockOpenProblem
	_ func(span) grammar.Problem             = malformedBlockCloseProblem
	_ func(string, []span) []grammar.Problem = unescapedBraceProblems
)

// TestGrammarProblemsAndRecovery проверяет P003–P008, P010, P012 и восстановление.
func TestGrammarProblemsAndRecovery(t *testing.T) {
	location := span{start: 2, end: 5}
	tests := []struct {
		name    string
		got     grammar.Problem
		kind    grammar.ProblemKind
		scope   model.DiagnosticScope
		fatal   bool
		element int
	}{
		{name: "unknown", got: unknownTagProblem(location), kind: grammar.ProblemUnknownTag, scope: model.ScopeElement, element: 0},
		{name: "separator", got: missingSeparatorProblem(location), kind: grammar.ProblemMissingSeparator, scope: model.ScopeLine, element: -1},
		{name: "form", got: unsupportedFormProblem(location, 1), kind: grammar.ProblemUnsupportedForm, scope: model.ScopeElement, element: 1},
		{name: "argument", got: missingArgumentProblem(location, 1), kind: grammar.ProblemMissingArgument, scope: model.ScopeElement, element: 1},
		{name: "extra", got: extraContentProblem(location, 2), kind: grammar.ProblemExtraContent, scope: model.ScopeElement, element: 2},
		{name: "block open", got: malformedBlockOpenProblem(location), kind: grammar.ProblemMalformedBlockOpen, scope: model.ScopeLine, fatal: true, element: -1},
		{name: "block close", got: malformedBlockCloseProblem(location), kind: grammar.ProblemMalformedBlockClose, scope: model.ScopeLine, element: -1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requireProblem(t, test.got, test.kind, test.scope, test.fatal, 2, 5, test.element)
		})
	}

	braces := unescapedBraceProblems(`x{y\}z}`, nil)
	if len(braces) != 2 {
		t.Fatalf("unescapedBraceProblems() = %#v", braces)
	}
	for index, bounds := range [][2]int{{1, 2}, {6, 7}} {
		requireProblem(t, braces[index], grammar.ProblemUnescapedBrace, model.ScopeElement, false, bounds[0], bounds[1], 0)
	}

	recoverable := recoveryFor([]grammar.Problem{unknownTagProblem(location)}, formResult{form: grammar.TagFormLine})
	if !recoverable.Continue {
		t.Fatal("recoverable decision stopped classification")
	}
	fatal := recoveryFor([]grammar.Problem{malformedBlockOpenProblem(location)}, formResult{ambiguous: true})
	if fatal.Continue {
		t.Fatal("ambiguous block-open decision continued classification")
	}
}
