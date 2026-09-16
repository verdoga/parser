package console

import "io"

// Status задаёт уже вычисленный статус строки файла.
type Status string

const (
	// Success обозначает успешный итог файла.
	Success Status = "success"
	// Failed обозначает неуспешный итог файла.
	Failed Status = "failed"
)

// FileLine содержит данные одной итоговой строки TXT-файла.
type FileLine struct {
	// Status содержит готовый статус и не вычисляется форматтером.
	Status Status
	// Action содержит машинное имя выполненного действия.
	Action string
	// Errors содержит число диагностик уровня error.
	Errors int
	// Path содержит абсолютный путь источника.
	Path string
	// Message содержит необязательное краткое сообщение.
	Message string
}

// SummaryLine содержит данные итоговой строки пакетного запуска.
type SummaryLine struct {
	// Found содержит число найденных TXT-файлов.
	Found int
	// Parsed содержит число полностью разобранных источников.
	Parsed int
	// Created содержит число созданных JSON.
	Created int
	// Replaced содержит число заменённых JSON.
	Replaced int
	// Skipped содержит число пропущенных актуальных JSON.
	Skipped int
	// Success содержит число успешных файлов.
	Success int
	// Failed содержит число неуспешных файлов.
	Failed int
	// Diagnostics содержит число диагностик уровня error.
	Diagnostics int
	// ScanErrors содержит число ошибок обхода.
	ScanErrors int
}

// FormatFile формирует одну строку результата без завершающего LF.
func FormatFile(line FileLine) string { panic("TODO") }

// FormatSummary формирует итоговую строку с нормативным порядком полей без завершающего LF.
func FormatSummary(line SummaryLine) string { panic("TODO") }

// Writer направляет итоговые и операционные строки в заданные потоки.
type Writer struct {
	// stdout принимает файловые и итоговые строки.
	stdout io.Writer
	// stderr принимает ошибки запуска и обхода.
	stderr io.Writer
}

// NewWriter создаёт консольный вывод с явными потоками.
func NewWriter(stdout, stderr io.Writer) Writer { panic("TODO") }

// File записывает ровно одну итоговую строку файла с LF.
func (w Writer) File(line FileLine) error { panic("TODO") }

// Summary записывает ровно одну итоговую строку запуска с LF.
func (w Writer) Summary(line SummaryLine) error { panic("TODO") }

// OperationalError записывает операционную ошибку в stderr с LF.
func (w Writer) OperationalError(err error) error { panic("TODO") }
