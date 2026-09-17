package diagnostics

import "dslparser/internal/model"

// Details содержит типизированные переменные части нормативного сообщения.
type Details struct {
	// Tag содержит каноническое имя связанного тега без символа @.
	Tag string
	// Fragment содержит исходный проблемный фрагмент объявления.
	Fragment string
	// Required содержит имя обязательного параметра, формы или разделителя.
	Required string
	// Received содержит фактически прочитанное значение.
	Received string
	// Supported содержит поддерживаемые версии в детерминированном порядке.
	Supported []string
	// Path содержит путь источника для технической ошибки чтения.
	Path string
	// Cause содержит исходную техническую ошибку чтения.
	Cause error
}

// Request задаёт одну обнаруженную проблему до назначения диагностического ID.
type Request struct {
	// Code содержит код из единого реестра parser.
	Code Code
	// Scope содержит область, разрешённую реестром для данного кода.
	Scope model.DiagnosticScope
	// Fatal равен true при невозможности штатно продолжить и false при однозначном восстановлении.
	Fatal bool
	// Location содержит основной диапазон либо nil, когда позиция недоступна или не существует.
	Location *model.Range
	// RelatedLocations содержит дополнительные уникальные диапазоны в нормативном порядке.
	RelatedLocations []model.Range
	// Details содержит данные для формирования сообщения реестром.
	Details Details
}
