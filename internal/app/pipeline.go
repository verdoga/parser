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

// ClockFunc возвращает текущее время для одной точки измерения.
// Возвращённое значение приложение приводит к UTC.
type ClockFunc func() time.Time

// ProcessingIDFunc создаёт непустой идентификатор, уникальный внутри результата.
type ProcessingIDFunc func() string

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

// ProcessFileFunc выполняет последовательный конвейер одного источника без форматирования консоли.
type ProcessFileFunc func(request ProcessRequest) FileResult

// Dependencies содержит явные подменяемые границы пакетного оркестратора.
type Dependencies struct {
	// Discover выполняет нормализацию пути и поиск файлов.
	Discover DiscoverFunc
	// BuildCache строит индекс до последовательной обработки TXT.
	BuildCache BuildCacheFunc
	// HeaderReader читает минимальные заголовки найденных JSON.
	HeaderReader cache.HeaderReader
	// Clock предоставляет время для новых попыток обработки.
	Clock ClockFunc
	// Processor выполняет конвейер одного файла.
	Processor ProcessFileFunc
}

// DefaultDependencies создаёт production-композицию стандартной файловой системы,
// grammar v1.2, кэша, parser, report и безопасной записи.
func DefaultDependencies() (Dependencies, error) { panic("TODO") }

// RunWithDependencies выполняет пакетный запуск с явно переданными зависимостями.
// Ошибка означает, что обработка не начиналась; завершённые файловые и обходные
// ошибки представлены в RunResult и его ExitCode.
func RunWithDependencies(options Options, dependencies Dependencies) (RunResult, error) {
	panic("TODO")
}
