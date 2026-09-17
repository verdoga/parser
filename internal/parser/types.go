package parser

import "dslparser/internal/model"

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
func ProbeSource(input Input, registry Registry) (Probe, []model.Diagnostic) { panic("TODO") }

// Parse выполняет декодирование, выбор грамматики и полный построчный структурный разбор.
func Parse(input Input, registry Registry) Result { panic("TODO") }
