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
	// Path содержит обязательный путь, переданный пользователем.
	Path string
	// Replace равен true, когда целевой JSON разрешено заменить, и false в обычном режиме.
	Replace bool
	// Depth задаёт максимальную глубину обхода; nil означает отсутствие ограничения.
	Depth *int
}

// Run проверяет исходный путь и запускает доступные этапы обработки.
//
// Структурный парсер будет подключён следующим вертикальным срезом. До этого
// момента найденный TXT нельзя ошибочно объявлять успешно обработанным.
func Run(options Options, stdout, stderr io.Writer) int {
	files, scanErrors, err := discovery.Discover(options.Path, options.Depth)
	if err != nil {
		var inputErr *discovery.InputError
		if errors.As(err, &inputErr) {
			fmt.Fprintf(stderr, "ОШИБКА: %v\n", err)
			return 2
		}
		fmt.Fprintf(stderr, "ОШИБКА: %v\n", err)
		return 1
	}

	for _, scanErr := range scanErrors {
		fmt.Fprintf(stderr, "ОШИБКА обхода: %v\n", scanErr)
	}

	// Нулевой набор входов уже является законченным успешным пакетным запуском.
	if len(files) == 0 {
		fmt.Fprintf(stdout, "ИТОГ found=0 parsed=0 created=0 replaced=0 skipped=0 success=0 failed=0 diagnostics=0 scanErrors=%d\n", len(scanErrors))
		if len(scanErrors) > 0 {
			return 1
		}
		return 0
	}

	for _, path := range files {
		fmt.Fprintf(stderr, "ОШИБКА action=failed errors=0 path=%q message=%q\n", path, "обработка DSL ещё не реализована")
	}
	fmt.Fprintf(stdout, "ИТОГ found=%d parsed=0 created=0 replaced=0 skipped=0 success=0 failed=%d diagnostics=0 scanErrors=%d\n", len(files), len(files), len(scanErrors))
	return 1
}
