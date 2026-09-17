package parser

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
)

// parseDecoded последовательно строит строки, состояние, метаданные и диагностики выбранной версии.
func parseDecoded(input Input, source physicalSource, selected grammar.Grammar, builder *diagnosticBuilder) Result {
	panic("TODO")
}

// buildDocument собирает доступные сведения источника без processing и сериализации.
func buildDocument(input Input, source physicalSource, version *string, metadata model.Metadata, diagnostics []model.Diagnostic) model.Document {
	panic("TODO")
}

// markLineErrors вычисляет признаки строк только по основным диапазонам и element-ссылкам.
func markLineErrors(lines []model.Line, diagnostics []model.Diagnostic) { panic("TODO") }

// hasFatal сообщает true при наличии фатальной диагностики и false при её отсутствии.
func hasFatal(diagnostics []model.Diagnostic) bool { panic("TODO") }

// hasErrors сообщает true при наличии error-диагностики и false при её отсутствии.
func hasErrors(diagnostics []model.Diagnostic) bool { panic("TODO") }
