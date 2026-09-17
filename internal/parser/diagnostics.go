package parser

import "dslparser/internal/model"

// diagnosticDraft содержит проблему до назначения уникального ID и нормативного сообщения.
type diagnosticDraft struct {
	// code содержит стабильный код parser.
	code string
	// scope содержит область диагностики.
	scope model.DiagnosticScope
	// fatal сообщает невозможность штатно продолжить разбор.
	fatal bool
	// location содержит основной диапазон либо nil.
	location *model.Range
	// related содержит дополнительные диапазоны.
	related []model.Range
}

// diagnosticBuilder создаёт диагностики одного запуска и сохраняет порядок обнаружения.
type diagnosticBuilder struct {
	// source содержит processing ID создаваемых диагностик.
	source string
	// items содержит уже созданные диагностики.
	items []model.Diagnostic
}

// newDiagnosticBuilder создаёт пустой накопитель одного запуска parser.
func newDiagnosticBuilder(source string) *diagnosticBuilder { panic("TODO") }

// add создаёт уникальную диагностику или возвращает существующую для того же кода и диапазона.
func (b *diagnosticBuilder) add(draft diagnosticDraft) model.Diagnostic { panic("TODO") }

// diagnostics возвращает отдельную копию диагностик в нормативном порядке.
func (b *diagnosticBuilder) diagnostics() []model.Diagnostic { panic("TODO") }
