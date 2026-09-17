package v1_2

import "dslparser/internal/grammar"

// implementation представляет неизменяемую реализацию грамматики DSL v1.2.
type implementation struct{}

// New создаёт независимую реализацию грамматики DSL v1.2.
func New() grammar.Grammar { panic("TODO") }

// Version возвращает точное поддерживаемое значение версии DSL.
func (implementation) Version() string { panic("TODO") }

// Classify классифицирует одну строку без изменения состояния parser.
func (implementation) Classify(request grammar.GrammarRequest) grammar.GrammarDecision {
	panic("TODO")
}
