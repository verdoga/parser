package model

// Position задаёт позицию в исходнике с нумерацией строк и Unicode-колонок от единицы.
type Position struct {
	// Line содержит номер физической строки от единицы.
	Line int `json:"line"`
	// Column содержит номер колонки в Unicode-кодовых точках от единицы.
	Column int `json:"column"`
}

// Range задаёт непрерывный диапазон с исключающей конечной позицией.
type Range struct {
	// Start содержит включительную начальную позицию.
	Start Position `json:"start"`
	// End содержит исключающую конечную позицию.
	End Position `json:"end"`
}
