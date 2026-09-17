package v1_2

import "dslparser/internal/grammar"

// Компиляционные проверки фиксируют контексты и решения состояния DSL v1.2.
var (
	_ grammar.Context                                                        = contextRoot
	_ grammar.Context                                                        = contextOpaque
	_ grammar.Context                                                        = contextExample
	_ grammar.Context                                                        = contextWordlist
	_ grammar.Context                                                        = contextTable
	_ grammar.Context                                                        = contextText
	_ grammar.Context                                                        = contextEditor
	_ grammar.Context                                                        = contextInstruction
	_ grammar.Context                                                        = contextFragment
	_ grammar.Context                                                        = contextChoice
	_ grammar.Context                                                        = contextMatching
	_ grammar.Context                                                        = contextMultifill
	_ grammar.Context                                                        = contextVariants
	_ func(tagSpec, grammar.TagForm) grammar.Transition                      = transitionFor
	_ func(tagSpec, grammar.TagForm, grammar.Context) grammar.ParentDecision = parentFor
	_ func(tagSpec) grammar.BlockDecision                                    = blockFor
)
