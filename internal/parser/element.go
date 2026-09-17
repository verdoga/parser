package parser

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
)

// elementFromDecision строит элемент и переводит байтовые смещения grammar в Unicode-колонки.
func elementFromDecision(line physicalLine, decision grammar.ElementDecision) model.Element {
	panic("TODO")
}

// emptyElement строит пустой элемент отсутствующего обязательного значения в позиции вставки.
func emptyElement(line physicalLine, column int, elementType model.ElementType) model.Element {
	panic("TODO")
}

// attachElementError добавляет уникальную ссылку на element-диагностику.
func attachElementError(element *model.Element, diagnosticID string) { panic("TODO") }

// sortElements упорядочивает элементы по начальной, затем по исключающей конечной колонке.
func sortElements(elements []model.Element) { panic("TODO") }
