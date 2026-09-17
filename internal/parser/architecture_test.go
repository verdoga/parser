package parser

import "dslparser/internal/model"

// Компиляционные проверки фиксируют внешние контракты парсера и границы грамматики.
var (
	_ Grammar  = architectureGrammar{}
	_ Registry = architectureRegistry{}
	_          = Input{Bytes: []byte{}, FileName: "lesson.txt", FilePath: "/lesson.txt", ProcessingID: "run-1"}
	_          = Result{Document: model.Document{}, Lines: []model.Line{}, Diagnostics: []model.Diagnostic{}}
)

// architectureGrammar представляет тестовую реализацию версии грамматики.
type architectureGrammar struct{}

// Version соответствует контракту Grammar.
func (architectureGrammar) Version() string { panic("TODO") }

// Classify соответствует контракту Grammar.
func (architectureGrammar) Classify(request GrammarRequest) GrammarDecision { panic("TODO") }

// architectureRegistry представляет тестовый реестр версий.
type architectureRegistry struct{}

// Lookup соответствует контракту Registry.
func (architectureRegistry) Lookup(version string) (Grammar, bool) { panic("TODO") }

// Versions соответствует контракту Registry.
func (architectureRegistry) Versions() []string { panic("TODO") }
