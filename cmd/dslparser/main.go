package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"dslparser/internal/app"
)

// version задаётся при сборке через -ldflags "-X main.version=<version>".
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(arguments []string, stdout, stderr io.Writer) int {
	options, err := parseArguments(arguments)
	if err != nil {
		fmt.Fprintf(stderr, "ОШИБКА: %v\n", err)
		return 2
	}
	options.ToolVersion = version
	return app.Run(options, stdout, stderr)
}

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

type optionalDepth struct {
	value int
	set   bool
}

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

func (depth *optionalDepth) String() string {
	if depth == nil || !depth.set {
		return ""
	}
	return strconv.Itoa(depth.value)
}

func (depth *optionalDepth) pointer() *int {
	if !depth.set {
		return nil
	}
	value := depth.value
	return &value
}
