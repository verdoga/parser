package model

// Document содержит доступные сведения об исходнике и вычисленные признаки результата.
type Document struct {
	// DSLVersion содержит прочитанную версию DSL либо nil при её отсутствии.
	DSLVersion *string `json:"dslVersion"`
	// FileName содержит имя исходного файла с расширением либо nil, если оно неизвестно.
	FileName *string `json:"fileName"`
	// FilePath содержит абсолютный путь исходного файла либо nil, если он неизвестен.
	FilePath *string `json:"filePath"`
	// Encoding содержит UTF-8 для принятого входа либо nil, если кодировка не установлена.
	Encoding *string `json:"encoding"`
	// HasBOM сообщает наличие BOM; nil означает, что байтовый вход проверить нельзя.
	HasBOM *bool `json:"hasBom"`
	// LineCount содержит число физических строк; nil означает, что полное число неизвестно.
	LineCount *int `json:"lineCount"`
	// ByteLength содержит размер точных входных байтов; nil означает неизвестный размер.
	ByteLength *int `json:"byteLength"`
	// SHA256 содержит отпечаток точных входных байтов либо nil при недоступном источнике.
	SHA256 *string `json:"sha256"`
	// Metadata содержит однозначно извлечённые метаданные.
	Metadata Metadata `json:"metadata"`
	// HasErrors равен true при наличии error-диагностики и false при её отсутствии.
	HasErrors bool `json:"hasErrors"`
}

// Metadata содержит однозначно извлечённые объявления документа.
type Metadata struct {
	// DocumentID содержит единственный document-id с исходным регистром либо nil.
	DocumentID *string `json:"documentId"`
	// Title содержит текст единственного заголовка первого уровня либо nil.
	Title *string `json:"title"`
	// Subtitle содержит текст единственного заголовка второго уровня либо nil.
	Subtitle *string `json:"subtitle"`
	// Section содержит единственное значение section либо nil.
	Section *string `json:"section"`
	// Order содержит исходную запись положительного целого либо nil.
	Order *string `json:"order"`
	// ResourceDirs содержит пути в порядке объявлений; nil означает недоступное извлечение.
	ResourceDirs []string `json:"resourceDirs"`
}
