package v1_2

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
	"testing"
)

// request создаёт запрос с согласованными Raw и Trimmed для тестовой строки.
func request(line string, context grammar.Context, mode grammar.ContentMode) grammar.GrammarRequest {
	return grammar.GrammarRequest{Raw: line, Trimmed: line, Line: 7, Context: context, ContentMode: mode}
}

// value возвращает указатель на отдельную строку для ожидаемого элемента.
func value(text string) *string { return &text }

// requireSpan проверяет байтовые границы диапазона.
func requireSpan(t *testing.T, got span, start, end int) {
	t.Helper()
	if got.start != start || got.end != end {
		t.Fatalf("span = {%d, %d}, want {%d, %d}", got.start, got.end, start, end)
	}
}

// requireElement проверяет тип, значение и исключающую конечную границу элемента.
func requireElement(t *testing.T, got grammar.ElementDecision, typ model.ElementType, want *string, start, end int) {
	t.Helper()
	if got.Type != typ || got.ByteStart != start || got.ByteEnd != end {
		t.Fatalf("element = %#v, want type %q range [%d,%d)", got, typ, start, end)
	}
	if (got.Value == nil) != (want == nil) || got.Value != nil && *got.Value != *want {
		t.Fatalf("element value = %v, want %v", got.Value, want)
	}
}

// requireProblem проверяет полный общий контракт локальной проблемы.
func requireProblem(t *testing.T, got grammar.Problem, kind grammar.ProblemKind, scope model.DiagnosticScope, fatal bool, start, end, element int) {
	t.Helper()
	if got.Kind != kind || got.Scope != scope || got.Fatal != fatal || got.ByteStart != start || got.ByteEnd != end || got.Element != element {
		t.Fatalf("problem = %#v, want kind=%q scope=%q fatal=%t range=[%d,%d) element=%d", got, kind, scope, fatal, start, end, element)
	}
}
