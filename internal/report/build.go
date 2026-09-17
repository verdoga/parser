package report

import "dslparser/internal/model"

// Input содержит подготовленные части одной завершённой попытки обработки.
type Input struct {
	// Document содержит сведения об источнике и извлечённые метаданные.
	Document model.Document
	// Processing содержит завершённую попытку, для которой приложение уже зафиксировало длительность.
	Processing model.Processing
	// Lines содержит все доступные физические строки в исходном порядке.
	Lines []model.Line
	// Diagnostics содержит диагностики попытки в нормативном порядке.
	Diagnostics []model.Diagnostic
}

// Build собирает независимый полный результат и вычисляет document.hasErrors и line.hasErrors.
// Build сохраняет нормативный порядок входных данных, копирует изменяемые агрегаты
// и представляет обязательные массивы пустыми срезами, а не nil.
func Build(input Input) model.Result { panic("TODO") }

// documentHasErrors сообщает, содержит ли набор хотя бы одну error-диагностику.
// documentHasErrors возвращает true при наличии error и false для пустого набора
// либо набора только из warning и recommendation.
func documentHasErrors(diagnostics []model.Diagnostic) bool { panic("TODO") }

// lineHasErrors сообщает, относится ли error-диагностика непосредственно к строке.
// lineHasErrors возвращает true при начале основной location на строке либо при
// ссылке её элемента на error-диагностику; ошибки дочерних строк не наследуются.
// При отсутствии таких условий lineHasErrors возвращает false.
func lineHasErrors(line model.Line, diagnostics []model.Diagnostic) bool { panic("TODO") }

// errorDiagnosticIDs возвращает множество идентификаторов error-диагностик.
func errorDiagnosticIDs(diagnostics []model.Diagnostic) map[string]struct{} { panic("TODO") }
