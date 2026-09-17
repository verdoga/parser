package storage

// syncFile задаёт минимальные операции временного или каталожного файла.
type syncFile interface {
	// Write записывает очередную часть содержимого.
	Write(data []byte) (int, error)
	// Name возвращает путь открытого файла.
	Name() string
	// Sync синхронизирует содержимое или запись каталога.
	Sync() error
	// Close закрывает файл.
	Close() error
}

// fileOperations задаёт внешние операции, на сбоях которых проверяется rollback.
type fileOperations interface {
	// CreateTemp создаёт временный файл в целевом каталоге.
	CreateTemp(directory, pattern string) (syncFile, error)
	// Link атомарно устанавливает новый файл без замены существующего target.
	Link(oldPath, newPath string) error
	// Rename атомарно переименовывает объект в пределах файловой системы.
	Rename(oldPath, newPath string) error
	// Remove удаляет временный файл или завершённую резервную копию.
	Remove(path string) error
	// Open открывает каталог для синхронизации его записи.
	Open(path string) (syncFile, error)
}

// fileSystemWithOperations создаёт файловый адаптер с тестовой границей операций.
func fileSystemWithOperations(operations fileOperations) FileSystem { panic("TODO") }

// WriteMode задаёт допустимое отношение к существующему целевому файлу.
type WriteMode int

const (
	// CreateOnly запрещает изменение существующего результата.
	CreateOnly WriteMode = iota
	// ReplaceExisting разрешает безопасную замену вычисленного результата.
	ReplaceExisting
)

// Writer задаёт границу безопасной установки полностью сериализованного JSON.
type Writer interface {
	// WriteJSON записывает data через временный файл и соблюдает выбранный режим.
	WriteJSON(target string, data []byte, mode WriteMode) error
}

// WriteJSON полностью записывает, синхронизирует и закрывает временный файл до установки.
// WriteJSON не изменяет data. Неизвестный WriteMode и любой незавершённый этап возвращают
// WriteError; существующий target при ошибке должен сохранить прежнее содержимое.
func (FileSystem) WriteJSON(target string, data []byte, mode WriteMode) error { panic("TODO") }

// WriteError описывает этап неудачной безопасной записи.
type WriteError struct {
	// Operation содержит стабильное имя неудачного этапа.
	Operation string
	// Path содержит путь объекта файловой системы на этом этапе.
	Path string
	// Err содержит исходную причину.
	Err error
}

// Error возвращает контекст безопасной записи.
func (e *WriteError) Error() string { panic("TODO") }

// Unwrap возвращает исходную причину для errors.Is и errors.As.
func (e *WriteError) Unwrap() error { panic("TODO") }
