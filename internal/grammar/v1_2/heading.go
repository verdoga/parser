package v1_2

import "dslparser/internal/grammar"

// heading содержит синтаксически выделенные части заголовка DSL v1.2.
type heading struct {
	// level содержит диапазон маркера уровня.
	level span
	// title содержит диапазон обязательного свободного названия.
	title span
	// value содержит канонический уровень 1, 2 или 3.
	value string
	// valid равен true для поддерживаемой формы с названием и false для иной строки с начальным #.
	valid bool
}

// scanHeading выделяет уровень и название потенциального заголовка.
func scanHeading(value string) heading { panic("TODO") }

// classifyHeading создаёт решение для строки, начинающей распознавание заголовка.
func classifyHeading(request grammar.GrammarRequest) grammar.GrammarDecision { panic("TODO") }
