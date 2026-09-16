package cache

// Header содержит только поля существующего JSON, необходимые для сопоставления.
type Header struct {
	// Path содержит абсолютный очищенный путь JSON.
	Path string
	// FormatVersion содержит версию JSON-контракта.
	FormatVersion string
	// DSLVersion содержит версию DSL либо nil при отсутствии.
	DSLVersion *string
	// DocumentID содержит исходное написание идентификатора либо nil.
	DocumentID *string
	// SHA256 содержит отпечаток точных байтов источника либо nil.
	SHA256 *string
}

// HeaderReader задаёт границу чтения минимального заголовка JSON.
type HeaderReader interface {
	// ReadHeader читает только поля, необходимые для сопоставления результата.
	ReadHeader(path string) (Header, error)
}

// JSONReader читает минимальные заголовки средствами стандартной библиотеки.
type JSONReader struct{}

// ReadHeader читает заголовок, игнорируя неизвестные поля JSON.
func (JSONReader) ReadHeader(path string) (Header, error) { panic("TODO") }

// ReadError описывает невозможность прочитать или декодировать заголовок.
type ReadError struct {
	// Path содержит путь существующего JSON.
	Path string
	// Err содержит исходную причину.
	Err error
}

// Error возвращает контекст ошибки заголовка.
func (e *ReadError) Error() string { panic("TODO") }

// Unwrap возвращает исходную причину для errors.Is и errors.As.
func (e *ReadError) Unwrap() error { panic("TODO") }
