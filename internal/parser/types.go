package parser

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
)

// Input содержит точные байты и уже известные сведения одного источника.
type Input struct {
	// Bytes содержит точные байты полного источника, включая BOM и переводы строк.
	Bytes []byte
	// FileName содержит имя исходного файла с расширением.
	FileName string
	// FilePath содержит абсолютный очищенный путь исходного файла.
	FilePath string
	// SHA256 содержит отпечаток точных байтов нижнего регистра.
	SHA256 string
	// ProcessingID содержит источник создаваемых диагностик.
	ProcessingID string
}

// Probe содержит минимальные однозначные сведения для сопоставления с кэшем.
type Probe struct {
	// DSLVersion содержит прочитанную версию, включая неподдерживаемую, либо nil.
	DSLVersion *string
	// DocumentID содержит единственный однозначно извлечённый идентификатор либо nil.
	DocumentID *string
}

// Result содержит построенную parser часть результата без processing и сериализации.
type Result struct {
	// Document содержит сведения об источнике и извлечённые метаданные.
	Document model.Document
	// Lines содержит все доступные физические строки в исходном порядке.
	Lines []model.Line
	// Diagnostics содержит диагностики parser в нормативном порядке.
	Diagnostics []model.Diagnostic
	// Fatal равен true при фатальной диагностике этого запуска и false иначе.
	Fatal bool
}

// ProbeSource выполняет лёгкое извлечение версии и document-id без полного разбора.
// Синтаксическая проблема даёт отсутствующее поле Probe, но не ошибку: диагностики
// создаются только последующим полным Parse после решения не пропускать источник.
// ProbeSource не изменяет data и не сохраняет ссылку на переданный срез.
func ProbeSource(data []byte, registry grammar.Registry) (Probe, error) { panic("TODO") }

// Parse выполняет декодирование, выбор грамматики и полный построчный структурный разбор.
// Parse создаёт единственный diagnostics.Builder с Input.ProcessingID, не изменяет
// Input.Bytes и не сохраняет ссылку на переданный срез. Input.SHA256 должен описывать
// точные Input.Bytes; несоответствие является ошибкой контракта вызова.
func Parse(input Input, registry grammar.Registry) (Result, error) {
	panic("TODO")
}
