// Package app оркестрирует один пакетный запуск dslparser.
package app

import (
	"errors"
	"fmt"

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

// validate проверяет параметры внутренней точки входа до обращения к файловой системе.
func (o Options) validate() error { panic("TODO") }

// Run проверяет исходный путь и возвращает решения пакетного обработчика без консольного вывода.
//
// Структурный парсер будет подключён следующим вертикальным срезом. До этого
// момента найденный TXT нельзя ошибочно объявлять успешно обработанным.
func Run(options Options) (RunResult, error) {
	result, err := discovery.Discover(discovery.Request{Path: options.Path, Depth: options.Depth})
	if err != nil {
		exitCode := 1
		var inputErr *discovery.InputError
		if errors.As(err, &inputErr) {
			exitCode = 2
		}
		return RunResult{ExitCode: exitCode}, fmt.Errorf("запуск discovery: %w", err)
	}

	runResult := RunResult{
		Files:      make([]FileResult, 0, len(result.Files)),
		ScanErrors: append([]error(nil), result.ScanErrors...),
		Summary: Summary{
			Found:      len(result.Files),
			ScanErrors: len(result.ScanErrors),
		},
	}
	for _, path := range result.Files {
		runResult.Files = append(runResult.Files, FileResult{
			Path:    path,
			Status:  FileFailed,
			Action:  ActionFailed,
			Message: "обработка DSL ещё не реализована",
		})
	}
	runResult.Summary.Failed = len(result.Files)
	if runResult.Summary.Failed > 0 || runResult.Summary.ScanErrors > 0 {
		runResult.ExitCode = 1
	}
	return runResult, nil
}
