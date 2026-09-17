package v1_2

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
)

// Компиляционные проверки фиксируют внешнюю границу грамматики DSL v1.2.
var (
	_ grammar.Grammar                                                = implementation{}
	_ func() grammar.Grammar                                         = New
	_ func(string) (span, tagName, bool)                             = scanTag
	_ func(string) formResult                                        = detectForm
	_ func(string) resourcePathResult                                = scanResourcePaths
	_ func(string) []span                                            = placeholderRanges
	_ func(string, []span) []span                                    = unescapedBraceRanges
	_ func(string) string                                            = unescapeText
	_ func(string) string                                            = unescapeMediaSource
	_ func(string, contentPolicy) string                             = unescape
	_ func([]grammar.Problem, formResult) grammar.RecoveryDecision   = recoveryFor
	_ func(grammar.GrammarRequest) grammar.GrammarDecision           = classifyBlank
	_ func(string) heading                                           = scanHeading
	_ model.LineType                                                 = model.LineBlank
	_ func(model.ElementType, span, *string) grammar.ElementDecision = elementDecision
	_ func(model.ElementType, int) grammar.ElementDecision           = missingElementDecision
	_ func(span) grammar.ElementDecision                             = unparsedElementDecision
	_ func(span, string) grammar.ElementDecision                     = contentRemainderDecision
)
