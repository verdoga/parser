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
	Path  string
	Depth *int
}

// Result содержит нормализованный исходный путь и найденные TXT-файлы.
type Result struct {
	Root       string
	IsFile     bool
	Files      []string
	ScanErrors []error
}

// InputError означает ошибку исходного пути, при которой обработка не начинается.
type InputError struct {
	err error
}

func (e *InputError) Error() string { return e.err.Error() }
func (e *InputError) Unwrap() error { return e.err }

// Discover определяет тип пути и находит обычные TXT-файлы в пределах глубины.
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

	result := Result{Root: absolute, Files: []string{}, ScanErrors: []error{}}
	walkErr := filepath.WalkDir(absolute, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			result.ScanErrors = append(result.ScanErrors, fmt.Errorf("не удалось прочитать %q: %w", path, walkErr))
			if path == absolute {
				return walkErr
			}
			return fs.SkipDir
		}

		if path == absolute {
			return nil
		}
		relative, err := filepath.Rel(absolute, path)
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
			if request.Depth != nil && depth > *request.Depth {
				return fs.SkipDir
			}
			return nil
		}
		if request.Depth != nil && depth > *request.Depth {
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
		return Result{}, inputError("не удалось прочитать каталог %q: %v", absolute, walkErr)
	}
	slices.Sort(result.Files)
	return result, nil
}

func inputError(format string, args ...any) error {
	return &InputError{err: fmt.Errorf(format, args...)}
}

func isTXT(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".txt")
}

func pathDepth(relative string) int {
	directory := filepath.Dir(relative)
	if directory == "." {
		return 0
	}
	return len(strings.Split(directory, string(filepath.Separator)))
}
