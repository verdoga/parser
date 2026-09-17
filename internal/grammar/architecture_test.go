package grammar

import "testing"

// Компиляционные проверки фиксируют границу реестра и реализаций версий.
var (
	_ Grammar                                          = architectureGrammar{}
	_ func(...Grammar) (Registry, error)               = NewRegistry
	_ func(Registry, string) (Grammar, bool)           = Registry.Lookup
	_ func(Registry) []string                          = Registry.Versions
	_ func(TagForm, TagForm) bool                      = TagForm.Allows
	_ func(ParentKind, ...ParentKind) ParentDecision   = Parents
	_ error                                            = (*RegistrationError)(nil)
	_ func(*RegistrationError) string                  = (*RegistrationError).Version
	_ func(*RegistrationError) RegistrationErrorReason = (*RegistrationError).Reason
)

// TestProblemKindValues проверяет полное типизированное соответствие построчным кодам parser.
func TestProblemKindValues(t *testing.T) {
	tests := []struct {
		name string
		kind ProblemKind
		want string
	}{
		{name: "unknown tag", kind: ProblemUnknownTag, want: "P003"},
		{name: "missing separator", kind: ProblemMissingSeparator, want: "P004"},
		{name: "unsupported form", kind: ProblemUnsupportedForm, want: "P005"},
		{name: "missing argument", kind: ProblemMissingArgument, want: "P006"},
		{name: "extra content", kind: ProblemExtraContent, want: "P007"},
		{name: "malformed block open", kind: ProblemMalformedBlockOpen, want: "P008"},
		{name: "unexpected block close", kind: ProblemUnexpectedBlockClose, want: "P009"},
		{name: "malformed block close", kind: ProblemMalformedBlockClose, want: "P010"},
		{name: "unescaped brace", kind: ProblemUnescapedBrace, want: "P012"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if string(test.kind) != test.want {
				t.Fatalf("ProblemKind = %q, want %q", test.kind, test.want)
			}
		})
	}
}

// architectureGrammar представляет тестовую реализацию контракта грамматики.
type architectureGrammar struct{}

// Version соответствует контракту Grammar.
func (architectureGrammar) Version() string { return "architecture" }

// Classify соответствует контракту Grammar.
func (architectureGrammar) Classify(GrammarRequest) GrammarDecision { return GrammarDecision{} }
