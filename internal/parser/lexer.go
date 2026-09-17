package parser

import "dslparser/internal/grammar"

// tokenKind задаёт внутренний вид лексического фрагмента.
type tokenKind int

const (
	// tokenTag обозначает имя тега.
	tokenTag tokenKind = iota
	// tokenText обозначает текстовый фрагмент.
	tokenText
	// tokenBlockOpen обозначает структурную открывающую скобку.
	tokenBlockOpen
	// tokenBlockClose обозначает структурную закрывающую скобку.
	tokenBlockClose
	// tokenSeparator обозначает структурный разделитель.
	tokenSeparator
)

// token содержит точный фрагмент через внутренние байтовые границы.
type token struct {
	// kind содержит вид фрагмента.
	kind tokenKind
	// raw содержит точный фрагмент исходной строки.
	raw string
	// byteStart содержит включительное байтовое смещение.
	byteStart int
	// byteEnd содержит исключающее байтовое смещение.
	byteEnd int
}

// lexicalLine содержит физическую строку и выделенные токены.
type lexicalLine struct {
	// physical содержит исходное представление строки.
	physical physicalLine
	// tokens содержит токены в порядке диапазонов.
	tokens []token
}

// lexLine выбирает допустимый лексический режим для текущего контекста.
func lexLine(line physicalLine, mode grammar.ContentMode) lexicalLine { panic("TODO") }

// lexDeclaration выделяет токены потенциального объявления DSL.
func lexDeclaration(line physicalLine) lexicalLine { panic("TODO") }

// lexOpaqueContent сохраняет содержимое и распознаёт только разрешённые границы режима.
func lexOpaqueContent(line physicalLine, mode grammar.ContentMode) lexicalLine { panic("TODO") }

// unescapeValue снимает только экранирование, значимое для заданного режима.
func unescapeValue(raw string, mode grammar.EscapeMode) string { panic("TODO") }
