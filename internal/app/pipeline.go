package app

import (
	"time"

	"dslparser/internal/cache"
	"dslparser/internal/discovery"
)

// discoverFunc находит входные TXT и область существующих JSON.
type discoverFunc func(request discovery.Request) (discovery.Result, error)

// buildCacheFunc строит неизменяемый индекс минимальных заголовков JSON.
type buildCacheFunc func(paths []string, reader cache.HeaderReader) (cache.Index, error)

// clockFunc возвращает текущее время для одной точки измерения.
// Возвращённое значение приложение приводит к UTC.
type clockFunc func() time.Time

// processingIDFunc создаёт непустой идентификатор, уникальный внутри результата.
type processingIDFunc func() string

// processRequest задаёт параметры обработки одного файла.
type processRequest struct {
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

// processFileFunc выполняет последовательный конвейер одного источника без форматирования консоли.
type processFileFunc func(request processRequest) FileResult

// dependencies содержит явные подменяемые границы пакетного оркестратора.
type dependencies struct {
	// Discover выполняет нормализацию пути и поиск файлов.
	Discover discoverFunc
	// BuildCache строит индекс до последовательной обработки TXT.
	BuildCache buildCacheFunc
	// HeaderReader читает минимальные заголовки найденных JSON.
	HeaderReader cache.HeaderReader
	// Clock предоставляет время для новых попыток обработки.
	Clock clockFunc
	// processor выполняет конвейер одного файла.
	processor processFileFunc
}

// defaultDependencies создаёт production-композицию стандартной файловой системы,
// grammar v1.2, кэша, parser, report и безопасной записи.
func defaultDependencies() (dependencies, error) { panic("TODO") }

// runWithDependencies выполняет пакетный запуск с явно переданными зависимостями.
// Ошибка означает, что обработка не начиналась; завершённые файловые и обходные
// ошибки представлены в RunResult и его ExitCode.
func runWithDependencies(options Options, dependencies dependencies) (RunResult, error) {
	panic("TODO")
}
