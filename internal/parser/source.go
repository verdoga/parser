package parser

import "dslparser/internal/model"

// physicalSource содержит проверенные физические строки и признак начального BOM.
type physicalSource struct {
	// hasBOM сообщает, начинались ли точные байты с UTF-8 BOM.
	hasBOM bool
	// lines содержит физические строки в исходном порядке.
	lines []physicalLine
}

// physicalLine содержит точное декодированное представление одной физической строки.
type physicalLine struct {
	// number содержит номер строки от единицы.
	number int
	// raw содержит строку без начального BOM и EOL.
	raw string
	// eol содержит точное поддерживаемое окончание строки.
	eol model.LineEnding
}

// sourceFailure описывает первичную фатальную проблему входных байтов.
type sourceFailure struct {
	// code содержит P001 или P002.
	code string
	// fatal сообщает невозможность штатного структурного разбора.
	fatal bool
	// location содержит доступный диапазон либо nil без надёжной позиции.
	location *model.Range
	// fragment содержит исходный проблемный фрагмент для нормативного сообщения.
	fragment string
}

// decodeSource проверяет UTF-8, BOM и переводы строк без нормализации байтов.
func decodeSource(data []byte) (physicalSource, *sourceFailure) { panic("TODO") }

// splitLines выделяет физические строки без создания фиктивной строки после завершающего EOL.
func splitLines(text string) ([]physicalLine, *sourceFailure) { panic("TODO") }

// hasUTF8BOM сообщает true только для BOM в начале и false в остальных случаях.
func hasUTF8BOM(data []byte) bool { panic("TODO") }
