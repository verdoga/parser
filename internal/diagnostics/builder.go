package diagnostics

import "dslparser/internal/model"

// Builder накапливает error-диагностики одного processing-запуска в порядке обнаружения.
// Builder принадлежит одной операции и не предназначен для одновременного использования.
type Builder struct {
	// source содержит непустой processing ID для всех создаваемых диагностик.
	source string
	// items содержит созданные диагностики в порядке обнаружения.
	items []model.Diagnostic
}

// New создаёт пустой накопитель для непустого processing ID.
func New(source string) (*Builder, error) { panic("TODO") }

// Add проверяет запрос и возвращает диагностику с уникальным ID, не сохраняя агрегаты запроса.
// Второй результат равен true для новой диагностики и false для ранее созданной
// диагностики с тем же кодом и основным диапазоном.
func (b *Builder) Add(request Request) (model.Diagnostic, bool, error) { panic("TODO") }

// Items возвращает глубокую копию диагностик в порядке их обнаружения.
func (b *Builder) Items() []model.Diagnostic { panic("TODO") }

// validateRequest проверяет код, область, фатальность, диапазоны и обязательные детали.
func validateRequest(request Request) error { panic("TODO") }

// validateDetails проверяет переменные части сообщения для конкретного кода.
func validateDetails(code Code, details Details) error { panic("TODO") }

// validateLocation проверяет наличие и форму основного диапазона по области диагностики.
func validateLocation(scope model.DiagnosticScope, location *model.Range) error { panic("TODO") }

// validateRelatedLocations проверяет дополнительные диапазоны и отсутствие повторов.
func validateRelatedLocations(location *model.Range, related []model.Range) error { panic("TODO") }

// findDuplicate возвращает копию ранее созданной диагностики и true для того же кода и диапазона.
// Для новой проблемы findDuplicate возвращает нулевую диагностику и false.
func (b *Builder) findDuplicate(code Code, location *model.Range) (model.Diagnostic, bool) {
	panic("TODO")
}

// build создаёт диагностику из уже проверенного запроса без изменения входных срезов.
func (b *Builder) build(request Request, definition definition) model.Diagnostic { panic("TODO") }

// diagnosticID создаёт ID, уникальный между запусками с уникальными processing ID.
func diagnosticID(source string, ordinal int) string { panic("TODO") }

// sameLocation сообщает true для двух одинаковых или одновременно отсутствующих диапазонов.
// sameLocation возвращает false, если присутствует только один диапазон или координаты различаются.
func sameLocation(left, right *model.Range) bool { panic("TODO") }

// cloneRange возвращает отдельную копию диапазона либо nil для отсутствующего диапазона.
func cloneRange(value *model.Range) *model.Range { panic("TODO") }

// cloneRanges возвращает непустую независимую копию либо новый пустой срез вместо nil.
func cloneRanges(values []model.Range) []model.Range { panic("TODO") }

// cloneDiagnostic возвращает диагностику без общих изменяемых срезов и указателей.
func cloneDiagnostic(value model.Diagnostic) model.Diagnostic { panic("TODO") }
