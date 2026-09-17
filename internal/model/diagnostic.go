package model

// Severity задаёт закрытое перечисление уровней диагностики.
type Severity string

const (
	// SeverityError обозначает ошибку.
	SeverityError Severity = "error"
	// SeverityWarning обозначает предупреждение.
	SeverityWarning Severity = "warning"
	// SeverityRecommendation обозначает рекомендацию.
	SeverityRecommendation Severity = "recommendation"
)

// DiagnosticScope задаёт закрытое перечисление областей диагностики.
type DiagnosticScope string

const (
	// ScopeElement связывает диагностику с одним или несколькими элементами.
	ScopeElement DiagnosticScope = "element"
	// ScopeLine связывает диагностику с одной строкой.
	ScopeLine DiagnosticScope = "line"
	// ScopeBlock связывает диагностику с диапазоном блока.
	ScopeBlock DiagnosticScope = "block"
	// ScopeDocument связывает диагностику с документом в целом.
	ScopeDocument DiagnosticScope = "document"
)

// Diagnostic содержит одну проблему, созданную конкретным запуском инструмента.
type Diagnostic struct {
	// ID содержит непустой уникальный в пределах результата идентификатор.
	ID string `json:"id"`
	// Source ссылается на идентификатор существующей записи processing.
	Source string `json:"source"`
	// Code содержит стабильный машинный код из реестра инструмента.
	Code string `json:"code"`
	// Severity содержит уровень диагностики.
	Severity Severity `json:"severity"`
	// Message содержит человекочитаемое описание проблемы.
	Message string `json:"message"`
	// Scope содержит область привязки диагностики.
	Scope DiagnosticScope `json:"scope"`
	// Fatal равен true при невозможности штатно продолжить запуск и false иначе.
	Fatal bool `json:"fatal"`
	// Location содержит основной диапазон либо nil при отсутствии исходной позиции.
	Location *Range `json:"location"`
	// RelatedLocations содержит уникальные дополнительные диапазоны в нормативном порядке.
	RelatedLocations []Range `json:"relatedLocations"`
}
