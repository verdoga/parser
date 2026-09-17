package parser

// sourceCase описывает будущий табличный тест декодирования точных байтов.
type sourceCase struct {
	name string
	data []byte
}

// grammarCase описывает будущий табличный тест формы или контекста DSL.
type grammarCase struct {
	name string
	line string
}

// diagnosticCase описывает будущую проверку кода, области и фатальности диагностики.
type diagnosticCase struct {
	name  string
	code  string
	fatal bool
}

// recoveryCase описывает будущую проверку восстановления и лексического хвоста.
type recoveryCase struct {
	name string
	text string
}
