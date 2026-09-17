package parser

import (
	"dslparser/internal/diagnostics"
	"dslparser/internal/grammar"
)

// requestForProblem преобразует версионную проблему grammar в запрос единого реестра.
func requestForProblem(line physicalLine, problem grammar.Problem) diagnostics.Request { panic("TODO") }

// requestForSourceFailure преобразует первичную ошибку байтового входа в запрос реестра.
func requestForSourceFailure(failure sourceFailure) diagnostics.Request { panic("TODO") }

// unsupportedVersionRequest создаёт P014 с фактической и поддерживаемыми версиями.
func unsupportedVersionRequest(version string, supported []string, line physicalLine) diagnostics.Request {
	panic("TODO")
}
