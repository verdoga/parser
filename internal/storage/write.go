package storage

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
