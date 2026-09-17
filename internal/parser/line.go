package parser

import (
	"dslparser/internal/diagnostics"
	"dslparser/internal/grammar"
	"dslparser/internal/model"
)

// lineResult содержит завершённую строку и результат изменения состояния.
type lineResult struct {
	// line содержит построенную строку модели.
	line model.Line
	// fatal сообщает необходимость остановить структурный разбор после строки.
	fatal bool
}

// parseLine выбирает разбор строки по контексту и применяет только однозначное решение.
func parseLine(line physicalLine, state *parseState, selected grammar.Grammar, builder *diagnostics.Builder) (lineResult, error) {
	panic("TODO")
}

// parseHeading строит корневой заголовок с нулевой логической глубиной.
func parseHeading(line physicalLine, state parseState) lineResult { panic("TODO") }

// parseDeclaration разбирает форму и параметры объявления через выбранную grammar.
func parseDeclaration(line physicalLine, state *parseState, selected grammar.Grammar, builder *diagnostics.Builder) (lineResult, error) {
	panic("TODO")
}

// parseBlockEnd сопоставляет отдельную закрывающую скобку только с вершиной стека.
func parseBlockEnd(line physicalLine, state *parseState, builder *diagnostics.Builder) (lineResult, error) {
	panic("TODO")
}

// parseContent строит содержимое в соответствии с активным грамматическим контекстом.
func parseContent(line physicalLine, state parseState, selected grammar.Grammar, builder *diagnostics.Builder) (lineResult, error) {
	panic("TODO")
}
