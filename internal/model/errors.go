package model

// InvariantError описывает одно нарушение внутреннего контракта модели.
type InvariantError struct {
	// Path содержит стабильный путь к ошибочному полю.
	Path string
	// Rule содержит понятное описание нарушенного правила.
	Rule string
}

// Error возвращает описание нарушенного инварианта.
func (e *InvariantError) Error() string { panic("TODO") }
