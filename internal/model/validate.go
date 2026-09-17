package model

// Validate проверяет модель без исправления, нормализации и изменения её срезов.
func (r Result) Validate() error { panic("TODO") }

// validateDocument проверяет версию формата, сведения документа и вычисленные признаки.
func validateDocument(result Result) error { panic("TODO") }

// validateProcessing проверяет порядок, уникальность и значения запусков инструментов.
func validateProcessing(items []Processing) error { panic("TODO") }

// validateLines проверяет физические строки, родителей, вложенность и элементы.
func validateLines(document Document, lines []Line) error { panic("TODO") }

// validateDiagnostics проверяет диагностики, их ссылки, диапазоны и признаки ошибок.
func validateDiagnostics(processing []Processing, lines []Line, diagnostics []Diagnostic) error {
	panic("TODO")
}

// validateLine проверяет одну строку относительно уже проверенных предшественников.
func validateLine(line Line, previous []Line) error { panic("TODO") }

// validateElement проверяет тип, диапазон, raw и ссылки одного элемента.
func validateElement(line Line, element Element) error { panic("TODO") }

// validateRange проверяет непрерывный диапазон относительно доступных физических строк.
func validateRange(path string, value Range, lines []Line) error { panic("TODO") }

// runeSlice возвращает фрагмент по включительной начальной и исключающей конечной колонкам.
// Второй результат равен true для допустимого диапазона и false для выхода за границы.
func runeSlice(raw string, start, end int) (string, bool) { panic("TODO") }
