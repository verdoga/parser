package v1_2

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
	"reflect"
	"testing"
)

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

// TestHeadingScanner проверяет три уровня, обязательный разделитель и название.
func TestHeadingScanner(t *testing.T) {
	tests := []struct {
		value string
		want  heading
	}{
		{value: "# Unit", want: heading{level: span{0, 1}, title: span{2, 6}, value: "1", valid: true}},
		{value: "##\tРаздел", want: heading{level: span{0, 2}, title: span{3, 15}, value: "2", valid: true}},
		{value: "### Step 1", want: heading{level: span{0, 3}, title: span{4, 10}, value: "3", valid: true}},
		{value: "#"}, {value: "# "}, {value: "#### Unsupported"}, {value: "#No separator"}, {value: "text"},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			if got := scanHeading(test.value); got != test.want {
				t.Fatalf("scanHeading(%q) = %#v, want %#v", test.value, got, test.want)
			}
		})
	}
}

// TestClassifyEverySupportedTag проверяет хотя бы одну нормативную форму каждого из 34 тегов.
func TestClassifyEverySupportedTag(t *testing.T) {
	tests := []struct {
		raw      string
		lineType model.LineType
	}{
		{"@dsl-version 1.2", model.LineTag}, {"@document-id doc-1", model.LineTag},
		{"@section Studentbook", model.LineTag}, {"@order 20", model.LineTag},
		{"@resource-dir Audio, \"Shared, files\"", model.LineTag}, {"@header Vocabulary", model.LineTag},
		{"@task 1a", model.LineTag}, {"@endtask", model.LineTag}, {"@step a", model.LineTag},
		{"@speaking", model.LineTag}, {"@newpage", model.LineTag}, {"@editor note", model.LineTag},
		{"@media audio track.mp3", model.LineTag}, {"@example Sample", model.LineTag},
		{"@wordlist word; translation", model.LineTag}, {"@table {", model.LineBlockStart},
		{"@script script.js", model.LineTag}, {"@text {", model.LineBlockStart}, {"@key answer", model.LineTag},
		{"@instr Read", model.LineTag}, {"@note Note", model.LineTag}, {"@alt Alternative", model.LineTag},
		{"@hint Hint", model.LineTag}, {"@fragment intro {", model.LineBlockStart},
		{"@include intro", model.LineTag}, {"@answer yes", model.LineTag}, {"@question Why?", model.LineTag},
		{"@multifill Fill _____{it}", model.LineTag}, {"@choice {", model.LineBlockStart},
		{"@multichoice {", model.LineBlockStart}, {"@matching {", model.LineBlockStart},
		{"@ordering {", model.LineBlockStart}, {"@variants {", model.LineBlockStart},
		{"@variant Student-A", model.LineTag},
	}
	if len(tests) != supportedTagCount {
		t.Fatalf("test table has %d entries, want %d", len(tests), supportedTagCount)
	}
	for _, test := range tests {
		t.Run(test.raw, func(t *testing.T) {
			decision := New().Classify(testRequest(test.raw, contextRoot, grammar.ContentDSL))
			if decision.LineType != test.lineType || len(decision.Problems) != 0 {
				t.Fatalf("Classify() = type %q, problems %#v; want %q without problems", decision.LineType, decision.Problems, test.lineType)
			}
			if len(decision.Elements) == 0 || decision.Elements[0].Type != model.ElementTag || decision.Elements[0].Value == nil {
				t.Fatalf("tag element missing: %#v", decision.Elements)
			}
		})
	}
}

// TestDeclarationDiagnostics проверяет восстанавливаемые ошибки формы объявлений P003–P008.
func TestDeclarationDiagnostics(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		kind grammar.ProblemKind
	}{
		{name: "unknown", raw: "@wat value", kind: grammar.ProblemUnknownTag},
		{name: "separator", raw: "@task-id", kind: grammar.ProblemMissingSeparator},
		{name: "unsupported block", raw: "@task id {", kind: grammar.ProblemUnsupportedForm},
		{name: "missing argument", raw: "@task", kind: grammar.ProblemMissingArgument},
		{name: "extra content", raw: "@endtask tail", kind: grammar.ProblemExtraContent},
		{name: "malformed open", raw: "@text { tail", kind: grammar.ProblemMalformedBlockOpen},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decision := New().Classify(testRequest(test.raw, contextRoot, grammar.ContentDSL))
			if len(decision.Problems) == 0 || decision.Problems[0].Kind != test.kind {
				t.Fatalf("problems = %#v, want first %q", decision.Problems, test.kind)
			}
		})
	}
}

// TestParentForRules проверяет приоритеты логических родителей всех закрытых правил.
func TestParentForRules(t *testing.T) {
	tests := []struct {
		rule parentRule
		want []grammar.ParentKind
	}{
		{parentDocument, nil},
		{parentContent, []grammar.ParentKind{grammar.ParentBlock, grammar.ParentInnerStep, grammar.ParentTask, grammar.ParentVariant, grammar.ParentStep}},
		{parentTask, []grammar.ParentKind{grammar.ParentInnerStep, grammar.ParentTask}},
		{parentStep, []grammar.ParentKind{grammar.ParentStep}},
		{parentVariants, []grammar.ParentKind{grammar.ParentVariants}},
		{parentVariantOrStep, []grammar.ParentKind{grammar.ParentVariant, grammar.ParentStep}},
	}
	for _, test := range tests {
		spec := tagSpec{lineParent: test.rule, blockParent: test.rule}
		for _, form := range []grammar.TagForm{grammar.TagFormLine, grammar.TagFormBlock} {
			if got := parentFor(spec, form, contextRoot); !reflect.DeepEqual(got.Kinds, test.want) {
				t.Errorf("parentFor(rule %d, form %d) = %v, want %v", test.rule, form, got.Kinds, test.want)
			}
		}
	}
}

// TestBlockFor проверяет полное решение об открытии блока.
func TestBlockFor(t *testing.T) {
	spec := tagSpec{name: tagText, block: blockSpec{context: contextText, mode: grammar.ContentOpaque}}
	want := grammar.BlockDecision{Action: grammar.BlockOpen, Tag: "text", Context: contextText, ContentMode: grammar.ContentOpaque}
	if got := blockFor(spec); !reflect.DeepEqual(got, want) {
		t.Fatalf("blockFor() = %#v, want %#v", got, want)
	}
}
