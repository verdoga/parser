package app

import (
	"time"

	"dslparser/internal/cache"
	"dslparser/internal/discovery"
)

// DiscoverFunc находит входные TXT и область существующих JSON.
type DiscoverFunc func(request discovery.Request) (discovery.Result, error)

// BuildCacheFunc строит неизменяемый индекс минимальных заголовков JSON.
type BuildCacheFunc func(paths []string, reader cache.HeaderReader) (cache.Index, error)

// Clock предоставляет время начала попытки обработки.
type Clock interface {
	// Now возвращает текущее время; приложение приводит его к UTC.
	Now() time.Time
}

// ProcessingIDGenerator создаёт идентификатор нового запуска обработки.
type ProcessingIDGenerator interface {
	// NewProcessingID возвращает непустой идентификатор, уникальный внутри результата.
	NewProcessingID() string
}

// SourceProbe содержит минимальные сведения до запуска полного структурного разбора.
type SourceProbe struct {
	// Path содержит абсолютный очищенный путь источника.
	Path string
	// TargetPath содержит вычисленный путь целевого JSON.
	TargetPath string
	// DSLVersion содержит прочитанную версию либо пустую строку при её отсутствии.
	DSLVersion string
	// DocumentID содержит однозначный идентификатор либо пустую строку.
	DocumentID string
	// SHA256 содержит отпечаток точных исходных байтов.
	SHA256 string
}

// ProcessRequest задаёт параметры обработки одного файла.
type ProcessRequest struct {
	// Path содержит абсолютный очищенный путь TXT-файла.
	Path string
	// Replace разрешает замену только вычисленного целевого JSON.
	Replace bool
	// ToolVersion содержит версию исполняемого файла.
	ToolVersion string
	// StartedAt содержит зафиксированное начало попытки.
	StartedAt time.Time
	// Cache содержит общий неизменяемый индекс области пакетного запуска.
	Cache cache.Index
}

// FileProcessor выполняет последовательный конвейер одного источника.
type FileProcessor interface {
	// Process обрабатывает один TXT и возвращает принятое решение без форматирования консоли.
	Process(request ProcessRequest) FileResult
}

// Dependencies содержит явные подменяемые границы пакетного оркестратора.
type Dependencies struct {
	// Discover выполняет нормализацию пути и поиск файлов.
	Discover DiscoverFunc
	// BuildCache строит индекс до последовательной обработки TXT.
	BuildCache BuildCacheFunc
	// HeaderReader читает минимальные заголовки найденных JSON.
	HeaderReader cache.HeaderReader
	// Clock предоставляет время для новых попыток обработки.
	Clock Clock
	// IDs создаёт processing ID только для формируемых результатов.
	IDs ProcessingIDGenerator
	// Processor выполняет конвейер одного файла.
	Processor FileProcessor
}

// RunWithDependencies выполняет пакетный запуск с явно переданными зависимостями.
func RunWithDependencies(options Options, dependencies Dependencies) RunResult { panic("TODO") }
