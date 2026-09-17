package parser

// trimmedLine содержит распознаваемый фрагмент и его байтовые границы в raw.
type trimmedLine struct {
	// text содержит строку без крайних U+0020 и TAB.
	text string
	// byteStart содержит включительное начало text в raw.
	byteStart int
	// byteEnd содержит исключающий конец text в raw.
	byteEnd int
}

// runeColumn переводит допустимое байтовое смещение в номер Unicode-колонки от единицы.
func runeColumn(raw string, byteOffset int) int { panic("TODO") }

// runeRange переводит включительную и исключающую байтовые границы в Unicode-колонки.
func runeRange(raw string, byteStart, byteEnd int) (start, end int) { panic("TODO") }

// trimRecognitionSpace удаляет для распознавания только крайние U+0020 и TAB.
func trimRecognitionSpace(raw string) trimmedLine { panic("TODO") }
