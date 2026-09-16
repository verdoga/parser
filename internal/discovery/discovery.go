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

// InputError означает ошибку исходного пути, при которой обработка не начинается.
type InputError struct {
	// err содержит исходную ошибку с контекстом пути.
	err error
}

// Error возвращает описание ошибки исходного пути.
func (e *InputError) Error() string { return e.err.Error() }

// Unwrap возвращает исходную ошибку для стандартной проверки цепочки ошибок.
func (e *InputError) Unwrap() error { return e.err }

// directoryWalker накапливает результат одного последовательного обхода каталога.
type directoryWalker struct {
	// root содержит абсолютный очищенный путь корневого каталога.
	root string
	// depth задаёт максимальную глубину; nil означает отсутствие ограничения.
	depth *int
	// files содержит найденные пути в порядке обхода до окончательной сортировки.
	files []string
	// scanErrors содержит независимые ошибки чтения элементов каталога.
	scanErrors []error
}

// Discover возвращает отсортированные TXT-файлы и отдельные ошибки обхода.
func Discover(path string, depth *int) ([]string, []error, error) {
	absolute, info, err := inspectInput(path)
	if err != nil {
		return nil, nil, err
	}
	if info.Mode().IsRegular() {
		if !isTXT(absolute) {
			return nil, nil, inputError("исходный файл %q должен иметь расширение .txt", absolute)
		}
		return []string{absolute}, []error{}, nil
	}

	return discoverDirectory(absolute, depth)
}

// inspectInput нормализует исходный путь и проверяет допустимость его типа.
func inspectInput(path string) (string, fs.FileInfo, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", nil, inputError("не удалось получить абсолютный путь %q: %v", path, err)
	}
	absolute = filepath.Clean(absolute)

	info, err := os.Lstat(absolute)
	if err != nil {
		return "", nil, inputError("недоступен исходный путь %q: %v", absolute, err)
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		return "", nil, inputError("исходный путь %q является символической ссылкой", absolute)
	}
	if !info.Mode().IsRegular() && !info.IsDir() {
		return "", nil, inputError("исходный путь %q не является обычным файлом или каталогом", absolute)
	}
	return absolute, info, nil
}

// discoverDirectory обходит каталог и сортирует найденные пути.
func discoverDirectory(root string, depth *int) ([]string, []error, error) {
	walker := directoryWalker{
		root:       root,
		depth:      depth,
		files:      []string{},
		scanErrors: []error{},
	}
	walkErr := filepath.WalkDir(root, walker.visit)
	if walkErr != nil && len(walker.scanErrors) == 0 {
		return nil, nil, inputError("не удалось прочитать каталог %q: %v", root, walkErr)
	}
	slices.Sort(walker.files)
	return walker.files, walker.scanErrors, nil
}

// visit обрабатывает один элемент дерева и управляет ограничением глубины.
func (walker *directoryWalker) visit(path string, entry fs.DirEntry, walkErr error) error {
	if walkErr != nil {
		walker.scanErrors = append(walker.scanErrors, fmt.Errorf("не удалось прочитать %q: %w", path, walkErr))
		if path == walker.root {
			return walkErr
		}
		return fs.SkipDir
	}
	if path == walker.root {
		return nil
	}

	relative, err := filepath.Rel(walker.root, path)
	if err != nil {
		walker.scanErrors = append(walker.scanErrors, fmt.Errorf("не удалось определить глубину %q: %w", path, err))
		if entry.IsDir() {
			return fs.SkipDir
		}
		return nil
	}
	depth := pathDepth(relative)
	if entry.Type()&fs.ModeSymlink != 0 {
		return nil
	}
	if entry.IsDir() {
		if walker.depth != nil && depth > *walker.depth {
			return fs.SkipDir
		}
		return nil
	}
	if walker.depth != nil && depth > *walker.depth {
		return nil
	}

	return walker.addRegularTXT(path, entry)
}

// addRegularTXT добавляет обычный TXT-файл либо сохраняет ошибку получения сведений.
func (walker *directoryWalker) addRegularTXT(path string, entry fs.DirEntry) error {
	info, err := entry.Info()
	if err != nil {
		walker.scanErrors = append(walker.scanErrors, fmt.Errorf("не удалось проверить %q: %w", path, err))
		return nil
	}
	if info.Mode().IsRegular() && isTXT(path) {
		walker.files = append(walker.files, filepath.Clean(path))
	}
	return nil
}

// inputError создаёт типизированную ошибку исходного пути.
func inputError(format string, args ...any) error {
	return &InputError{err: fmt.Errorf(format, args...)}
}

// isTXT сообщает, имеет ли путь расширение .txt без учёта регистра.
func isTXT(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".txt")
}

// pathDepth возвращает глубину каталога, содержащего относительный путь.
func pathDepth(relative string) int {
	directory := filepath.Dir(relative)
	if directory == "." {
		return 0
	}
	return len(strings.Split(directory, string(filepath.Separator)))
}
