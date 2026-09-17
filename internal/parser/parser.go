package parser

import (
	"dslparser/internal/diagnostics"
	"dslparser/internal/grammar"
	"dslparser/internal/model"
)

// parseDecoded последовательно строит строки, состояние, метаданные и диагностики выбранной версии.
func parseDecoded(input Input, source physicalSource, selected grammar.Grammar, builder *diagnostics.Builder) (Result, error) {
	panic("TODO")
}

// buildDocument собирает доступные сведения источника без processing и сериализации.
func buildDocument(input Input, source physicalSource, version *string, metadata model.Metadata) model.Document {
	panic("TODO")
}

// hasFatal сообщает true при наличии фатальной диагностики и false при её отсутствии.
func hasFatal(diagnostics []model.Diagnostic) bool { panic("TODO") }
