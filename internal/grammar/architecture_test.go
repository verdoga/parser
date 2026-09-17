package grammar

// Компиляционные проверки фиксируют границу реестра и реализаций версий.
var (
	_ Grammar                                          = architectureGrammar{}
	_ func(...Grammar) (Registry, error)               = NewRegistry
	_ func(Registry, string) (Grammar, bool)           = Registry.Lookup
	_ func(Registry) []string                          = Registry.Versions
	_ func(TagForm, TagForm) bool                      = TagForm.Allows
	_ error                                            = (*RegistrationError)(nil)
	_ func(*RegistrationError) string                  = (*RegistrationError).Version
	_ func(*RegistrationError) RegistrationErrorReason = (*RegistrationError).Reason
)

// architectureGrammar представляет тестовую реализацию контракта грамматики.
type architectureGrammar struct{}

// Version соответствует контракту Grammar.
func (architectureGrammar) Version() string { panic("TODO") }

// Classify соответствует контракту Grammar.
func (architectureGrammar) Classify(request GrammarRequest) GrammarDecision { panic("TODO") }
