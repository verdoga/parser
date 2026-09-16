// Package app оркестрирует один пакетный запуск dslparser.
package app

import (
	"errors"
	"fmt"
	"io"

	"dslparser/internal/discovery"
)

// Options содержит параметры, полученные из командной строки.
type Options struct {
	// Path задаёт исходный файл или каталог.
	Path string
	// Replace разрешает замену целевого результата.
	Replace bool
	// Depth ограничивает глубину поиска; nil снимает ограничение.
	Depth *int
	// ToolVersion задаёт версию исполняемого файла для результата обработки.
	ToolVersion string
}

// Run проверяет исходный путь и запускает доступные этапы обработки.
//
// Структурный парсер будет подключён следующим вертикальным срезом. До этого
// момента найденный TXT нельзя ошибочно объявлять успешно обработанным.
func Run(options Options, stdout, stderr io.Writer) int {
	result, err := discovery.Discover(discovery.Request{Path: options.Path, Depth: options.Depth})
	if err != nil {
		var inputErr *discovery.InputError
		if errors.As(err, &inputErr) {
			fmt.Fprintf(stderr, "ОШИБКА: %v\n", err)
			return 2
		}
		fmt.Fprintf(stderr, "ОШИБКА: %v\n", err)
		return 1
	}

	for _, scanErr := range result.ScanErrors {
		fmt.Fprintf(stderr, "ОШИБКА обхода: %v\n", scanErr)
	}

	// Нулевой набор входов уже является законченным успешным пакетным запуском.
	if len(result.Files) == 0 {
		fmt.Fprintf(stdout, "ИТОГ found=0 parsed=0 created=0 replaced=0 skipped=0 success=0 failed=0 diagnostics=0 scanErrors=%d\n", len(result.ScanErrors))
		if len(result.ScanErrors) > 0 {
			return 1
		}
		return 0
	}

	for _, path := range result.Files {
		fmt.Fprintf(stderr, "ОШИБКА action=failed errors=0 path=%q message=%q\n", path, "обработка DSL ещё не реализована")
	}
	fmt.Fprintf(stdout, "ИТОГ found=%d parsed=0 created=0 replaced=0 skipped=0 success=0 failed=%d diagnostics=0 scanErrors=%d\n", len(result.Files), len(result.Files), len(result.ScanErrors))
	return 1
}
