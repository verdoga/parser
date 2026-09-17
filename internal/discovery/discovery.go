// Package discovery находит входные TXT-файлы и нормализует исходный путь.
package discovery

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Request задаёт параметры поиска входных файлов.
type Request struct {
	// Path задаёт исходный файл или каталог.
	Path string
	// Depth ограничивает глубину каталогов; nil снимает ограничение.
	Depth *int
}

// Result содержит нормализованный исходный путь и найденные TXT-файлы.
type Result struct {
	// Root содержит абсолютный очищенный исходный путь.
	Root string
	// IsFile равен true для явно заданного файла и false для каталога.
	IsFile bool
	// Files содержит абсолютные очищенные пути в лексикографическом порядке.
	Files []string
	// JSONFiles содержит абсолютные очищенные пути обычных JSON-файлов из каталогов с найденными TXT.
	JSONFiles []string
	// ScanErrors содержит ошибки отдельных участков обхода в порядке обнаружения.
	ScanErrors []error
}

// InputError означает ошибку исходного пути, при которой обработка не начинается.
type InputError struct {
	// err хранит причину невозможности начать обработку.
	err error
}

// Error возвращает описание ошибки исходного пути.
func (e *InputError) Error() string { return e.err.Error() }

// Unwrap возвращает исходную ошибку для errors.Is и errors.As.
func (e *InputError) Unwrap() error { return e.err }

// Discover определяет тип пути, находит TXT в пределах глубины и собирает JSON
// только из каталогов, в которых найден хотя бы один TXT. Для явно переданного
// TXT область JSON ограничена его родительским каталогом. Ошибки просмотра
// отдельных каталогов возвращаются в Result.ScanErrors и не скрывают найденные пути.
func Discover(request Request) (Result, error) {
	absolute, err := filepath.Abs(request.Path)
	if err != nil {
		return Result{}, inputError("не удалось получить абсолютный путь %q: %v", request.Path, err)
	}
	absolute = filepath.Clean(absolute)

	info, err := os.Lstat(absolute)
	if err != nil {
		return Result{}, inputError("недоступен исходный путь %q: %v", absolute, err)
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		return Result{}, inputError("исходный путь %q является символической ссылкой", absolute)
	}

	if info.Mode().IsRegular() {
		if !isTXT(absolute) {
			return Result{}, inputError("исходный файл %q должен иметь расширение .txt", absolute)
		}
		return Result{Root: absolute, IsFile: true, Files: []string{absolute}}, nil
	}
	if !info.IsDir() {
		return Result{}, inputError("исходный путь %q не является обычным файлом или каталогом", absolute)
	}

	return discoverDirectory(absolute, request.Depth)
}

// discoverDirectory собирает обычные TXT-файлы внутри корневого каталога.
func discoverDirectory(root string, depthLimit *int) (Result, error) {
	result := Result{Root: root, Files: []string{}, ScanErrors: []error{}}
	walkErr := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			result.ScanErrors = append(result.ScanErrors, fmt.Errorf("не удалось прочитать %q: %w", path, walkErr))
			if path == root {
				return walkErr
			}
			return fs.SkipDir
		}

		if path == root {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			result.ScanErrors = append(result.ScanErrors, fmt.Errorf("не удалось определить глубину %q: %w", path, err))
			if entry.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		depth := pathDepth(relative)
		if entry.Type()&fs.ModeSymlink != 0 {
			if entry.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			if depthLimit != nil && depth > *depthLimit {
				return fs.SkipDir
			}
			return nil
		}
		if depthLimit != nil && depth > *depthLimit {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			result.ScanErrors = append(result.ScanErrors, fmt.Errorf("не удалось проверить %q: %w", path, err))
			return nil
		}
		if info.Mode().IsRegular() && isTXT(path) {
			result.Files = append(result.Files, filepath.Clean(path))
		}
		return nil
	})
	if walkErr != nil && len(result.ScanErrors) == 0 {
		return Result{}, inputError("не удалось прочитать каталог %q: %v", root, walkErr)
	}
	slices.Sort(result.Files)
	return result, nil
}

// discoverJSONFiles возвращает обычные JSON-файлы только из каталогов найденных TXT.
// Результат не содержит дубликатов и отсортирован по абсолютному очищенному пути.
func discoverJSONFiles(txtFiles []string) ([]string, []error) { panic("TODO") }

// inputError создаёт типизированную ошибку исходного пути с форматированным сообщением.
func inputError(format string, args ...any) error {
	return &InputError{err: fmt.Errorf(format, args...)}
}

// isTXT сообщает, имеет ли путь расширение .txt без учёта регистра.
func isTXT(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".txt")
}

// pathDepth возвращает число каталогов между корнем и объектом.
func pathDepth(relative string) int {
	directory := filepath.Dir(relative)
	if directory == "." {
		return 0
	}
	return len(strings.Split(directory, string(filepath.Separator)))
}
