package model

// Компиляционные проверки фиксируют типизированный публичный контракт модели.
var (
	_ error = (*InvariantError)(nil)
	_       = Result{FormatVersion: FormatVersion}
	_       = Document{Metadata: Metadata{ResourceDirs: []string{}}}
	_       = Line{Type: LineContent, EOL: EOLLF, Elements: []Element{}}
	_       = Diagnostic{Severity: SeverityError, Scope: ScopeDocument, RelatedLocations: []Range{}}
)
