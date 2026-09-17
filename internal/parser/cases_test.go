package parser

// sourceCase описывает будущий табличный тест декодирования точных байтов.
type sourceCase struct {
	// name содержит имя подтеста.
	name string
	// data содержит точные проверяемые байты.
	data []byte
}

// grammarCase описывает будущий табличный тест формы или контекста DSL.
type grammarCase struct {
	// name содержит имя подтеста.
	name string
	// line содержит проверяемую физическую строку без EOL.
	line string
}

// diagnosticCase описывает будущую проверку кода, области и фатальности диагностики.
type diagnosticCase struct {
	// name содержит имя подтеста.
	name string
	// code содержит ожидаемый стабильный код.
	code string
	// fatal содержит ожидаемый признак остановки разбора.
	fatal bool
}

// recoveryCase описывает будущую проверку восстановления и лексического хвоста.
type recoveryCase struct {
	// name содержит имя подтеста.
	name string
	// text содержит полный проверяемый источник.
	text string
}
