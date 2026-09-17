package v1_2

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
	"reflect"
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

// TestProblemConstructors проверяет код, область, летальность, диапазон и связь с элементом для всех локальных проблем.
func TestProblemConstructors(t *testing.T) {
	location := span{start: 2, end: 8}
	tests := []struct {
		name    string
		make    func() grammar.Problem
		kind    grammar.ProblemKind
		scope   model.DiagnosticScope
		fatal   bool
		element int
	}{
		{name: "unknown tag", make: func() grammar.Problem { return unknownTagProblem(location) }, kind: grammar.ProblemUnknownTag, scope: model.ScopeElement, element: -1},
		{name: "missing separator", make: func() grammar.Problem { return missingSeparatorProblem(location) }, kind: grammar.ProblemMissingSeparator, scope: model.ScopeLine, element: -1},
		{name: "unsupported form", make: func() grammar.Problem { return unsupportedFormProblem(location, 3) }, kind: grammar.ProblemUnsupportedForm, scope: model.ScopeElement, element: 3},
		{name: "missing argument", make: func() grammar.Problem { return missingArgumentProblem(location, 2) }, kind: grammar.ProblemMissingArgument, scope: model.ScopeElement, element: 2},
		{name: "extra content", make: func() grammar.Problem { return extraContentProblem(location, 4) }, kind: grammar.ProblemExtraContent, scope: model.ScopeElement, element: 4},
		{name: "malformed block open", make: func() grammar.Problem { return malformedBlockOpenProblem(location) }, kind: grammar.ProblemMalformedBlockOpen, scope: model.ScopeLine, fatal: true, element: -1},
		{name: "malformed block close", make: func() grammar.Problem { return malformedBlockCloseProblem(location) }, kind: grammar.ProblemMalformedBlockClose, scope: model.ScopeLine, fatal: true, element: -1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.make()
			if got.Kind != test.kind || got.Scope != test.scope || got.Fatal != test.fatal || got.ByteStart != 2 || got.ByteEnd != 8 || got.Element != test.element {
				t.Fatalf("problem = %#v", got)
			}
		})
	}
}

// TestUnescapedBraceProblems проверяет отдельную P012 для каждой незащищённой скобки.
func TestUnescapedBraceProblems(t *testing.T) {
	got := unescapedBraceProblems(`x {y} _____{ok} \{z\}`, []span{{6, 15}})
	if len(got) != 2 {
		t.Fatalf("problems count = %d, want 2", len(got))
	}
	for i, location := range []span{{2, 3}, {4, 5}} {
		problem := got[i]
		if problem.Kind != grammar.ProblemUnescapedBrace || problem.Scope != model.ScopeElement || problem.Fatal || problem.Element != 0 || problem.ByteStart != location.start || problem.ByteEnd != location.end {
			t.Errorf("problem[%d] = %#v", i, problem)
		}
	}
}

// TestRecoveryFor проверяет остановку только при летальной или неоднозначной ошибке.
func TestRecoveryFor(t *testing.T) {
	recoverable := []grammar.Problem{unknownTagProblem(span{0, 3}), extraContentProblem(span{4, 5}, 1)}
	if got := recoveryFor(recoverable, formResult{form: grammar.TagFormLine}); !got.Continue {
		t.Fatal("recoverable problems stopped parsing")
	}
	fatal := append([]grammar.Problem(nil), recoverable...)
	fatal = append(fatal, malformedBlockOpenProblem(span{6, 7}))
	if got := recoveryFor(fatal, formResult{form: grammar.TagFormLine}); got.Continue {
		t.Fatal("fatal problem allowed parsing to continue")
	}
	if got := recoveryFor(nil, formResult{ambiguous: true}); got.Continue {
		t.Fatal("ambiguous form allowed parsing to continue")
	}
	if got := recoveryFor(nil, formResult{form: grammar.TagFormBlock}); !reflect.DeepEqual(got, grammar.RecoveryDecision{Continue: true}) {
		t.Fatalf("clean recovery = %#v", got)
	}
}
