package grammar

// Registry предоставляет неизменяемый набор реализаций поддерживаемых версий DSL.
type Registry struct {
	// grammars содержит отдельную копию грамматик в порядке регистрации.
	grammars []Grammar
}

// NewRegistry проверяет версии и создаёт независимый реестр в порядке аргументов.
// При пустой или повторяющейся версии функция возвращает RegistrationError.
func NewRegistry(grammars ...Grammar) (Registry, error) { panic("TODO") }

// Lookup возвращает грамматику версии и true либо nil и false для неподдерживаемой версии.
func (r Registry) Lookup(version string) (Grammar, bool) { panic("TODO") }

// Versions возвращает отдельную копию поддерживаемых версий в порядке регистрации.
func (r Registry) Versions() []string { panic("TODO") }
