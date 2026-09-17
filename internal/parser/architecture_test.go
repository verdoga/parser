package parser

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
)

// Компиляционные проверки фиксируют внешние контракты парсера и границы грамматики.
var (
	_ grammar.Grammar                               = architectureGrammar{}
	_                                               = grammar.Registry{}
	_                                               = Input{Bytes: []byte{}, FileName: "lesson.txt", FilePath: "/lesson.txt", ProcessingID: "run-1"}
	_                                               = Result{Document: model.Document{}, Lines: []model.Line{}, Diagnostics: []model.Diagnostic{}}
	_ func([]byte, grammar.Registry) (Probe, error) = ProbeSource
	_ func(Input, grammar.Registry) (Result, error) = Parse
)

// architectureGrammar представляет тестовую реализацию версии грамматики.
type architectureGrammar struct{}

// Version соответствует контракту grammar.Grammar.
func (architectureGrammar) Version() string { return "architecture" }

// Classify соответствует контракту grammar.Grammar.
func (architectureGrammar) Classify(grammar.GrammarRequest) grammar.GrammarDecision {
	return grammar.GrammarDecision{}
}
