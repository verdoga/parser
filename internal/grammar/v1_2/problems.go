package v1_2

import "dslparser/internal/grammar"

// unknownTagProblem создаёт P003 для неизвестного имени тега.
func unknownTagProblem(location span) grammar.Problem { panic("TODO") }

// missingSeparatorProblem создаёт P004 для нарушенного обязательного разделителя.
func missingSeparatorProblem(location span) grammar.Problem { panic("TODO") }

// unsupportedFormProblem создаёт P005 для неподдерживаемой формы тега.
func unsupportedFormProblem(location span, element int) grammar.Problem { panic("TODO") }

// missingArgumentProblem создаёт P006 для отсутствующего структурного параметра.
func missingArgumentProblem(location span, element int) grammar.Problem { panic("TODO") }

// extraContentProblem создаёт P007 для лишнего остатка объявления.
func extraContentProblem(location span, element int) grammar.Problem { panic("TODO") }

// malformedBlockOpenProblem создаёт P008 для неправильной строки открытия блока.
func malformedBlockOpenProblem(location span) grammar.Problem { panic("TODO") }

// malformedBlockCloseProblem создаёт P010 для неправильной строки закрытия блока.
func malformedBlockCloseProblem(location span) grammar.Problem { panic("TODO") }

// unescapedBraceProblems создаёт P012 для каждой незащищённой текстовой скобки.
func unescapedBraceProblems(value string, protected []span) []grammar.Problem { panic("TODO") }

// recoveryFor сообщает возможность однозначного продолжения после обнаруженных проблем.
func recoveryFor(problems []grammar.Problem, form formResult) grammar.RecoveryDecision {
	panic("TODO")
}
