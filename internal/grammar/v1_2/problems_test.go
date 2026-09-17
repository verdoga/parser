package v1_2

import "dslparser/internal/grammar"

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
