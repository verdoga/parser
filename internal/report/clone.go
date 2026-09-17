package report

import "dslparser/internal/model"

// cloneDocument возвращает независимую копию документа и его изменяемых значений.
func cloneDocument(document model.Document) model.Document { panic("TODO") }

// cloneProcessing возвращает независимую копию сведений о завершённой попытке.
func cloneProcessing(processing model.Processing) model.Processing { panic("TODO") }

// cloneLines возвращает глубокую копию строк, элементов и ссылок на диагностики.
func cloneLines(lines []model.Line) []model.Line { panic("TODO") }

// cloneDiagnostics возвращает глубокую копию диагностик и связанных диапазонов.
func cloneDiagnostics(diagnostics []model.Diagnostic) []model.Diagnostic { panic("TODO") }

// cloneString возвращает указатель на независимую копию строки либо nil.
func cloneString(value *string) *string { panic("TODO") }

// cloneInt возвращает указатель на независимую копию целого числа либо nil.
func cloneInt(value *int) *int { panic("TODO") }

// cloneBool возвращает указатель на независимую копию логического значения либо nil.
func cloneBool(value *bool) *bool { panic("TODO") }

// cloneRange возвращает указатель на независимую копию диапазона либо nil.
func cloneRange(value *model.Range) *model.Range { panic("TODO") }
