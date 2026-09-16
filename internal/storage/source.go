package storage

// Source содержит прочитанный целиком источник и его неизменяемые сведения.
type Source struct {
	// Path содержит абсолютный очищенный путь источника.
	Path string
	// FileName содержит последний компонент Path.
	FileName string
	// Bytes содержит отдельную копию точных байтов, включая BOM и переводы строк.
	Bytes []byte
	// ByteLength содержит длину Bytes в байтах.
	ByteLength int
	// SHA256 содержит 64 строчных шестнадцатеричных символа.
	SHA256 string
}

// Reader задаёт необходимую приложению границу чтения файловой системы.
type Reader interface {
	// ReadSource читает точные байты одного обычного файла и вычисляет их SHA-256.
	ReadSource(path string) (Source, error)
}

// FileSystem реализует чтение и запись через стандартную файловую систему.
type FileSystem struct{}

// ReadSource читает источник без нормализации содержимого.
func (FileSystem) ReadSource(path string) (Source, error) { panic("TODO") }

// TargetPath заменяет последнее расширение источника на .json в том же каталоге.
func TargetPath(sourcePath string) (string, error) { panic("TODO") }
