package parser

import "dslparser/internal/model"

// elementFromToken строит элемент и переводит внутренние байтовые смещения в Unicode-колонки.
func elementFromToken(line physicalLine, token token, elementType model.ElementType, value *string) model.Element {
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
