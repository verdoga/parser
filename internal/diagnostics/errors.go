package diagnostics

// RequestError описывает нарушение контракта запроса или накопителя диагностик.
type RequestError struct {
	// Field содержит стабильное имя ошибочного поля.
	Field string
	// Rule содержит описание нарушенного правила.
	Rule string
}

// Error возвращает описание нарушенного контракта diagnostics.
func (e *RequestError) Error() string { panic("TODO") }
