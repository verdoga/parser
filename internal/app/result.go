package app

// FileStatus задаёт итоговый статус обработки одного исходного файла.
type FileStatus int

const (
	// FileSuccess означает успешную обработку или пропуск актуального результата.
	FileSuccess FileStatus = iota
	// FileFailed означает фатальную или операционную ошибку обработки файла.
	FileFailed
)

// FileAction задаёт выполненное над результатом файла действие.
type FileAction int

const (
	// ActionCreated означает создание нового результата.
	ActionCreated FileAction = iota
	// ActionReplaced означает безопасную замену существующего результата.
	ActionReplaced
	// ActionSkipped означает пропуск актуального результата.
	ActionSkipped
	// ActionFailed означает отсутствие успешно завершённого действия.
	ActionFailed
)

// String возвращает стабильное машинное имя статуса.
func (s FileStatus) String() string { panic("TODO") }

// String возвращает стабильное машинное имя действия.
func (a FileAction) String() string { panic("TODO") }

// FileResult описывает уже принятое приложением решение по одному TXT-файлу.
type FileResult struct {
	// Path содержит абсолютный очищенный путь источника.
	Path string
	// Status содержит итоговый статус файла.
	Status FileStatus
	// Action содержит итоговое действие над результатом.
	Action FileAction
	// ErrorCount содержит число диагностик уровня error.
	ErrorCount int
	// Message содержит краткое сообщение для неуспешного результата.
	Message string
}

// Summary содержит счетчики завершённого пакетного запуска.
type Summary struct {
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
	// Success содержит число файлов со статусом успеха.
	Success int
	// Failed содержит число файлов со статусом ошибки.
	Failed int
	// Diagnostics содержит суммарное число диагностик уровня error.
	Diagnostics int
	// ScanErrors содержит число ошибок обхода каталогов.
	ScanErrors int
}

// AddFile учитывает завершённый результат файла в статистике.
func (s *Summary) AddFile(result FileResult) { panic("TODO") }

// AddScanErrors учитывает ошибки обхода каталогов.
func (s *Summary) AddScanErrors(count int) { panic("TODO") }

// Code возвращает 0 при полном успехе и 1 при наличии файловых или обходных ошибок.
func (s Summary) Code() int { panic("TODO") }

// RunResult содержит упорядоченные результаты и итог пакетного запуска.
type RunResult struct {
	// Files содержит результаты в порядке обработки источников.
	Files []FileResult
	// Summary содержит итоговые счетчики.
	Summary Summary
	// ScanErrors содержит ошибки отдельных участков обхода в порядке обнаружения.
	ScanErrors []error
	// ExitCode содержит вычисленный приложением код завершения.
	ExitCode int
}
