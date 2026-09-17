package grammar

import (
	"errors"
	"reflect"
	"testing"
)

// TestTagFormAllows проверяет одиночные и объединённые маски форм.
func TestTagFormAllows(t *testing.T) {
	tests := []struct {
		name      string
		forms     TagForm
		candidate TagForm
		want      bool
	}{
		{name: "line", forms: TagFormLine, candidate: TagFormLine, want: true},
		{name: "block", forms: TagFormBlock, candidate: TagFormBlock, want: true},
		{name: "combined line", forms: TagFormLine | TagFormBlock, candidate: TagFormLine, want: true},
		{name: "combined block", forms: TagFormLine | TagFormBlock, candidate: TagFormBlock, want: true},
		{name: "different", forms: TagFormLine, candidate: TagFormBlock, want: false},
		{name: "zero candidate", forms: TagFormLine | TagFormBlock, candidate: 0, want: false},
		{name: "combined candidate", forms: TagFormLine | TagFormBlock, candidate: TagFormLine | TagFormBlock, want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.forms.Allows(test.candidate); got != test.want {
				t.Fatalf("Allows(%d) = %t, want %t", test.candidate, got, test.want)
			}
		})
	}
}

// TestParentsOwnsKinds проверяет приоритет и независимость среза родителей.
func TestParentsOwnsKinds(t *testing.T) {
	fallback := []ParentKind{ParentTask, ParentRoot}
	decision := Parents(ParentBlock, fallback...)
	fallback[0] = ParentVariant
	want := []ParentKind{ParentBlock, ParentTask, ParentRoot}
	if !reflect.DeepEqual(decision.Kinds, want) {
		t.Fatalf("Kinds = %v, want %v", decision.Kinds, want)
	}
}

// TestRegistry проверяет регистрацию, порядок, поиск и владение результатами.
func TestRegistry(t *testing.T) {
	first := stubGrammar{version: "first"}
	second := stubGrammar{version: "second"}
	registry, err := NewRegistry(first, second)
	if err != nil {
		t.Fatalf("NewRegistry() error: %v", err)
	}
	versions := registry.Versions()
	if !reflect.DeepEqual(versions, []string{"first", "second"}) {
		t.Fatalf("Versions() = %q", versions)
	}
	versions[0] = "changed"
	if got := registry.Versions()[0]; got != "first" {
		t.Fatalf("Versions() exposed state: %q", got)
	}
	got, ok := registry.Lookup("second")
	if !ok || got.Version() != "second" {
		t.Fatalf("Lookup() = (%v, %t)", got, ok)
	}
	if got, ok := registry.Lookup("missing"); ok || got != nil {
		t.Fatalf("missing Lookup() = (%v, %t)", got, ok)
	}
}

// TestRegistryRejectsInvalidEntries проверяет типизированные причины регистрации.
func TestRegistryRejectsInvalidEntries(t *testing.T) {
	tests := []struct {
		name     string
		grammars []Grammar
		reason   RegistrationErrorReason
		version  string
	}{
		{name: "nil", grammars: []Grammar{nil}, reason: RegistrationNilGrammar},
		{name: "empty", grammars: []Grammar{stubGrammar{}}, reason: RegistrationEmptyVersion},
		{name: "duplicate", grammars: []Grammar{stubGrammar{"duplicate"}, stubGrammar{"duplicate"}}, reason: RegistrationDuplicateVersion, version: "duplicate"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewRegistry(test.grammars...)
			var registrationError *RegistrationError
			if !errors.As(err, &registrationError) {
				t.Fatalf("error = %T, want *RegistrationError", err)
			}
			if registrationError.Reason() != test.reason || registrationError.Version() != test.version {
				t.Fatalf("error = reason %v, version %q", registrationError.Reason(), registrationError.Version())
			}
			if registrationError.Error() == "" {
				t.Fatal("Error() is empty")
			}
		})
	}
}

// stubGrammar реализует управляемую грамматику для реестра.
type stubGrammar struct{ version string }

// Version возвращает настроенную версию.
func (g stubGrammar) Version() string { return g.version }

// Classify возвращает пустое решение, поскольку реестр не классифицирует строки.
func (stubGrammar) Classify(GrammarRequest) GrammarDecision { return GrammarDecision{} }
