package cache

// Entry связывает путь существующего результата с заголовком либо ошибкой чтения.
type Entry struct {
	// Path содержит абсолютный очищенный путь JSON.
	Path string
	// Header содержит прочитанные поля; при Err они могут быть недоступны.
	Header Header
	// Err содержит ошибку чтения или декодирования заголовка.
	Err error
}

// Index хранит упорядоченные записи и регистронезависимый индекс document-id.
type Index struct {
	// entries хранит записи в лексикографическом порядке путей.
	entries []Entry
	// identifiers сопоставляет нормализованный ID с индексами entries.
	identifiers map[string][]int
}

// Build читает все переданные пути и строит неизменяемый индекс.
// Ошибка отдельного JSON сохраняется в Entry.Err и не прерывает чтение остальных путей.
// Возвращаемая ошибка означает недопустимую конфигурацию вызова, при которой индекс не построен.
func Build(paths []string, reader HeaderReader) (Index, error) { panic("TODO") }

// Entries возвращает копию записей в детерминированном порядке.
func (i Index) Entries() []Entry { panic("TODO") }

// SourceIdentity содержит однозначно извлечённые сведения проверяемого источника.
type SourceIdentity struct {
	// TargetPath содержит единственный вычисленный путь результата.
	TargetPath string
	// DocumentID содержит однозначный ID; пустая строка запрещает повторное использование.
	DocumentID string
	// DSLVersion содержит прочитанную версию DSL.
	DSLVersion string
	// SHA256 содержит отпечаток точных исходных байтов.
	SHA256 string
}

// MatchKind задаёт результат сопоставления источника с индексом.
type MatchKind int

const (
	// MatchMissing означает отсутствие целевого результата и конфликтов ID.
	MatchMissing MatchKind = iota
	// MatchCurrent означает полностью актуальный и однозначный целевой результат.
	MatchCurrent
	// MatchTargetConflict означает существующий, но непригодный целевой результат.
	MatchTargetConflict
	// MatchDuplicateDocumentID означает наличие другого JSON с тем же ID.
	MatchDuplicateDocumentID
)

// Match содержит классификацию и упорядоченные участвующие записи.
type Match struct {
	// Kind содержит итог сопоставления.
	Kind MatchKind
	// Target содержит целевую запись либо nil при её отсутствии.
	Target *Entry
	// Conflicts содержит копию конфликтующих записей в порядке путей.
	Conflicts []Entry
}

// MatchSource сопоставляет источник по пути, ID, версии и SHA-256.
func (i Index) MatchSource(source SourceIdentity) Match { panic("TODO") }

// CanSkip сообщает, разрешено ли пропустить полный разбор.
//
// CanSkip возвращает true только для MatchCurrent и false для любого иного результата.
func (m Match) CanSkip() bool { panic("TODO") }
