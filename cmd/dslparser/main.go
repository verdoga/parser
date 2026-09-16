package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"dslparser/internal/app"
)

// main запускает приложение с аргументами и потоками текущего процесса.
func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run разбирает аргументы и возвращает код завершения приложения.
func run(arguments []string, stdout, stderr io.Writer) int {
	options, err := parseArguments(arguments)
	if err != nil {
		fmt.Fprintf(stderr, "ОШИБКА: %v\n", err)
		return 2
	}
	return app.Run(options, stdout, stderr)
}

// parseArguments проверяет синтаксис командной строки и возвращает параметры запуска.
func parseArguments(arguments []string) (app.Options, error) {
	flags := flag.NewFlagSet("dslparser", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	var replace bool
	var depth optionalDepth
	flags.BoolVar(&replace, "replace", false, "заменить существующий результат")
	flags.Var(&depth, "depth", "глубина обхода каталогов")
	if err := flags.Parse(arguments); err != nil {
		return app.Options{}, fmt.Errorf("недопустимые аргументы: %w", err)
	}
	if flags.NArg() != 1 {
		if flags.NArg() == 0 {
			return app.Options{}, fmt.Errorf("не указан обязательный путь")
		}
		return app.Options{}, fmt.Errorf("ожидался ровно один путь, получено: %d", flags.NArg())
	}

	return app.Options{
		Path:    flags.Arg(0),
		Replace: replace,
		Depth:   depth.pointer(),
	}, nil
}

// optionalDepth хранит значение глубины и признак его явного присутствия.
type optionalDepth struct {
	// value содержит неотрицательную глубину обхода.
	value int
	// set равен true, если флаг передан, и false при отсутствии ограничения.
	set bool
}

// Set разбирает неотрицательную глубину из значения флага.
func (depth *optionalDepth) Set(raw string) error {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fmt.Errorf("глубина должна быть целым числом: %q", raw)
	}
	if value < 0 {
		return fmt.Errorf("глубина не может быть отрицательной: %d", value)
	}
	depth.value = value
	depth.set = true
	return nil
}

// String возвращает заданную глубину или пустую строку при отсутствии значения.
func (depth *optionalDepth) String() string {
	if depth == nil || !depth.set {
		return ""
	}
	return strconv.Itoa(depth.value)
}

// pointer возвращает копию заданной глубины или nil при отсутствии ограничения.
func (depth *optionalDepth) pointer() *int {
	if !depth.set {
		return nil
	}
	value := depth.value
	return &value
}
