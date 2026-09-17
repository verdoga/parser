package parser

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
)

// Компиляционные проверки фиксируют внешние контракты парсера и границы грамматики.
var (
	_ grammar.Grammar = architectureGrammar{}
	_                 = grammar.Registry{}
	_                 = Input{Bytes: []byte{}, FileName: "lesson.txt", FilePath: "/lesson.txt", ProcessingID: "run-1"}
	_                 = Result{Document: model.Document{}, Lines: []model.Line{}, Diagnostics: []model.Diagnostic{}}
)

// architectureGrammar представляет тестовую реализацию версии грамматики.
type architectureGrammar struct{}

// Version соответствует контракту grammar.Grammar.
func (architectureGrammar) Version() string { panic("TODO") }

// Classify соответствует контракту grammar.Grammar.
func (architectureGrammar) Classify(request grammar.GrammarRequest) grammar.GrammarDecision {
	panic("TODO")
}
